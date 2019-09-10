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

func traverseJSON(currentBody []byte, tree *node, filter bool, relationHandler func(u *url.URL, n *node)) []byte {
	var (
		newBody []byte
		err     error
	)

	result := gjson.ParseBytes(currentBody)
	if result.Type == gjson.String {
		// Maybe a relation
		return handleRelation(currentBody, tree, relationHandler)
	}

	filter = filter && tree.hasChildren(Fields)
	if filter {
		if result.IsArray() {
			newBody = []byte("[]")
		} else {
			newBody = []byte("{}")
		}
	} else {
		newBody = currentBody
	}

	for _, n := range tree.children {
		if filter {
			if !n.fields {
				// Don't push for nothing
				continue
			}
		}

		if n.path == "*" {
			var i int
			result.ForEach(func(_, value gjson.Result) bool {
				// TODO: support iterating over objects
				rawBytes := traverseJSON(getBytes(value, currentBody), n, filter, relationHandler)
				newBody, err = sjson.SetRawBytes(newBody, strconv.Itoa(i), rawBytes)
				if err != nil {
					log.WithFields(log.Fields{"path": n.path, "reason": err, "index": i}).Debug("Cannot update array")
				}

				i++
				return true
			})
			continue
		}

		path := unescape(n.path)
		result := gjson.GetBytes(currentBody, path)
		if result.Exists() {
			rawBytes := traverseJSON(getBytes(result, currentBody), n, filter, relationHandler)

			newBody, err = sjson.SetRawBytes(newBody, path, rawBytes)
			if err != nil {
				log.WithFields(log.Fields{"path": n.path, "reason": err}).Debug("Cannot update new document")
			}
		}
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
