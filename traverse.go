package vulcain

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/dunglas/httpsfv"
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
	preload := n.httpList(Preload, "")
	fields := n.httpList(Fields, "")

	q := u.Query()

	if len(preload) > 0 {
		if v, err := httpsfv.Marshal(preload); err == nil {
			q.Add("preload", v)
		}
	}

	if len(fields) > 0 {
		if v, err := httpsfv.Marshal(fields); err == nil {
			q.Add("fields", v)
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

func traverseJSON(currentBody []byte, tree *node, filter bool, relationHandler func(n *node, v string) string) []byte {
	var (
		newBody []byte
		err     error
	)

	result := gjson.ParseBytes(currentBody)
	switch result.Type {
	// Maybe a relation
	case gjson.String:
		return handleRelation(currentBody, result.String(), tree, relationHandler)
	case gjson.Number:
		return handleRelation(currentBody, strconv.FormatInt(result.Int(), 10), tree, relationHandler)
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
					log.WithFields(log.Fields{"node": n.String(), "reason": err, "index": i}).Debug("Cannot update array")
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
				log.WithFields(log.Fields{"node": n.String(), "reason": err}).Debug("Cannot update new document")
			}
		}
	}

	return newBody
}

func handleRelation(currentBody []byte, rel string, tree *node, relationHandler func(n *node, v string) string) []byte {
	if newValue := relationHandler(tree, rel); newValue != "" {
		newBody, _ := json.Marshal(newValue)
		return newBody
	}

	return currentBody
}
