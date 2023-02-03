package vulcain

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/dunglas/httpsfv"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.uber.org/zap"
)

// unescape unescapes an extended JSON pointer
func unescape(s string) string {
	s = strings.ReplaceAll(s, "~2", "*")
	s = strings.ReplaceAll(s, "~1", "/")
	return strings.ReplaceAll(s, "~0", "~")
}

// espaceSJSONPath escapes a sjson path
func espaceSJSONPath(s string) string {
	// https://github.com/tidwall/sjson/blob/master/sjson.go#L47
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "#", "\\#")
	s = strings.ReplaceAll(s, "@", "\\@")
	s = strings.ReplaceAll(s, "*", "\\*")

	return strings.ReplaceAll(s, "?", "\\?")
}

// urlRewriter rewrites an URL to propagate the "preload" and "fields" selectors to relations
func urlRewriter(u *url.URL, n *node) {
	p := n.httpList(preload, "")
	f := n.httpList(fields, "")

	q := u.Query()

	if len(p) > 0 {
		if v, err := httpsfv.Marshal(p); err == nil {
			q.Add("preload", v)
		}
	}

	if len(f) > 0 {
		if v, err := httpsfv.Marshal(f); err == nil {
			q.Add("fields", v)
		}
	}

	u.RawQuery = q.Encode()
}

// getBytes retrieves a slice of bytes
func getBytes(r gjson.Result, body []byte) []byte {
	if r.Index > 0 {
		return body[r.Index : r.Index+len(r.Raw)]
	}

	return []byte(r.Raw)
}

// traverseJSON traverses and modify if needed the JSON document
// it pushes the relations specified by a "preload" directive
func (v *Vulcain) traverseJSON(currentBody []byte, tree *node, filter bool, relationHandler func(n *node, v string) string) []byte {
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

	filter = filter && tree.hasChildren(fields)
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
				rawBytes := v.traverseJSON(getBytes(value, currentBody), n, filter, relationHandler)
				newBody, err = sjson.SetRawBytes(newBody, strconv.Itoa(i), rawBytes)
				if err != nil {
					v.logger.Debug("cannot update array", zap.Stringer("node", n), zap.Int("index", i), zap.Error(err))
				}

				i++
				return true
			})
			continue
		}

		path := espaceSJSONPath(unescape(n.path))

		result := gjson.GetBytes(currentBody, path)
		if result.Exists() {
			rawBytes := v.traverseJSON(getBytes(result, currentBody), n, filter, relationHandler)

			newBody, err = sjson.SetRawBytes(newBody, path, rawBytes)
			if err != nil {
				v.logger.Debug("cannot update new document", zap.Stringer("node", n), zap.Error(err))
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
