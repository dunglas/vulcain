package gateway

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

func unescape(s string) string {
	s = strings.ReplaceAll(s, "~2", "*")
	s = strings.ReplaceAll(s, "~1", "/")
	return strings.ReplaceAll(s, "~0", "~")
}

func urlRewriter(u *url.URL, n *node) {
	q := u.Query()
	for _, pp := range n.strings(Preload, "") {
		if pp != "/" {
			q.Add("preload", pp)
		}
	}
	for _, fp := range n.strings(Fields, "") {
		if fp != "/" {
			q.Add("fields", fp)
		}
	}
	u.RawQuery = q.Encode()
}

func getBytes(r gjson.Result, body []byte) []byte {
	if r.Index > 0 {
		return body[r.Index : r.Index+len(r.Raw)]
	}

	return []byte(r.Raw)
}

func (g *Gateway) traverseJSON(currentBody []byte, tree *node, filter bool, relationHandler func(u *url.URL, n *node)) []byte {
	var (
		newBody []byte
		err     error
	)

	if len(tree.children) == 0 {
		// Leaf, maybe a relation
		newBody = handleRelation(currentBody, tree, relationHandler)

		return currentBody
	}

	if filter {
		if len(tree.children) == 1 && tree.children[0].path == "*" {
			newBody = []byte("[]")
		} else {
			newBody = []byte("{}")
		}
	} else {
		newBody = currentBody
	}

	for _, node := range tree.children {
		if filter {
			if !node.fields {
				// Don't push for nothing
				continue
			}
		}

		if node.path == "*" {
			result := gjson.ParseBytes(currentBody)
			var i int
			result.ForEach(func(_, value gjson.Result) bool {
				// TODO: support iterating over objects
				rawBytes := getBytes(value, currentBody)
				rawBytes = g.traverseJSON(rawBytes, node, filter, relationHandler)
				newBody, err = sjson.SetRawBytes(newBody, strconv.Itoa(i), rawBytes)
				if err != nil {
					log.WithFields(log.Fields{"path": node.path, "reason": err, "index": i}).Debug("Cannot update array")
				}

				i++
				return true
			})
			continue
		}

		path := unescape(node.path)
		result := gjson.GetBytes(currentBody, path)
		if result.Exists() {
			rawBytes := g.traverseJSON(getBytes(result, currentBody), node, filter, relationHandler)

			newBody, err = sjson.SetRawBytes(newBody, path, rawBytes)
			if err != nil {
				log.WithFields(log.Fields{"path": node.path, "reason": err}).Debug("Cannot update new document")
			}

			continue
		}

		// Probably a relation, push it with the current context
		newBody = handleRelation(currentBody, tree, relationHandler)
		break
	}

	return newBody
}

func handleRelation(currentBody []byte, tree *node, relationHandler func(u *url.URL, n *node)) []byte {
	result := gjson.ParseBytes(currentBody)

	if result.Type != gjson.String {
		return currentBody
	}

	uStr := result.String()
	u, err := url.Parse(uStr)
	if err != nil {
		log.WithFields(log.Fields{"path": tree.path, "relation": uStr}).Debug("Invalid relation")
		return currentBody
	}

	relationHandler(u, tree)

	newBody, _ := json.Marshal(u.String())
	return newBody
}
