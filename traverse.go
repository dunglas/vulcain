package vulcain

import (
	"bytes"
	"encoding/json/jsontext"
	"net/url"
	"strconv"
	"strings"

	"github.com/dunglas/httpsfv"
	"go.uber.org/zap"
)

// unescape unescapes an extended JSON pointer
func unescape(s string) string {
	s = strings.ReplaceAll(s, "~2", "*")
	s = strings.ReplaceAll(s, "~1", "/")
	return strings.ReplaceAll(s, "~0", "~")
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

// traverser walks a JSON document in a single pass, splicing replacements over the original bytes
type traverser struct {
	dec             *jsontext.Decoder
	body            []byte
	key             []byte
	relationHandler func(n *node, v string) string
}

// traverseJSON traverses and modify if needed the JSON document
// it pushes the relations specified by a "preload" directive
func (v *Vulcain) traverseJSON(body []byte, tree *node, filter bool, relationHandler func(n *node, v string) string) []byte {
	t := &traverser{
		dec:             jsontext.NewDecoder(bytes.NewBuffer(body), jsontext.AllowDuplicateNames(true), jsontext.AllowInvalidUTF8(true)),
		body:            body,
		relationHandler: relationHandler,
	}

	start, end, value, err := t.walk(tree, filter)
	if err != nil {
		v.logger.Debug("cannot traverse JSON document", zap.Error(err))

		return body
	}

	if value == nil {
		return body
	}

	newBody := make([]byte, 0, len(body)-(end-start)+len(value))
	newBody = append(newBody, body[:start]...)
	newBody = append(newBody, value...)

	return append(newBody, body[end:]...)
}

// walk consumes the next value and returns its offsets in the body and its replacement, nil if unchanged
func (t *traverser) walk(n *node, filter bool) (int, int, []byte, error) {
	switch t.dec.PeekKind() {
	case '"', '0':
		if n.preload || n.hasChildren(preload) || n.hasChildren(fields) {
			return t.relation(n)
		}
	case '{', '[':
		if n.hasChildren(preload) || n.hasChildren(fields) {
			return t.container(n, filter && n.hasChildren(fields))
		}
	}

	raw, err := t.dec.ReadValue()
	if err != nil {
		return 0, 0, nil, err
	}

	end := int(t.dec.InputOffset())

	return end - len(raw), end, nil, nil
}

// relation passes a string or number value to the relation handler
func (t *traverser) relation(n *node) (int, int, []byte, error) {
	raw, err := t.dec.ReadValue()
	if err != nil {
		return 0, 0, nil, err
	}

	end := int(t.dec.InputOffset())
	start := end - len(raw)

	rel := string(raw)
	if raw.Kind() == '"' {
		// Unquoting replaces the invalid Unicode accepted by the decoder.
		b, _ := jsontext.AppendUnquote(nil, raw)
		rel = string(b)
	}

	newValue := t.relationHandler(n, rel)
	if newValue == "" {
		return start, end, nil, nil
	}

	value, _ := jsontext.AppendQuote(nil, newValue)

	return start, end, value, nil
}

// container walks the members of an object or the elements of an array
func (t *traverser) container(n *node, filter bool) (int, int, []byte, error) {
	tok, err := t.dec.ReadToken()
	if err != nil {
		return 0, 0, nil, err
	}

	start := int(t.dec.InputOffset()) - 1
	isObject := tok.Kind() == '{'
	closing := jsontext.Kind(']')
	if isObject {
		closing = '}'
	}

	// Filtered containers are rebuilt from scratch as compact JSON
	var value []byte
	if filter {
		value = []byte{byte(tok.Kind())}
	}

	last := start
	nextIndex := 0
	for i := 0; t.dec.PeekKind() != closing; i++ {
		var name []byte
		if isObject {
			raw, err := t.dec.ReadValue()
			if err != nil {
				return 0, 0, nil, err
			}

			nameEnd := int(t.dec.InputOffset())
			name = t.body[nameEnd-len(raw) : nameEnd]

			t.key, _ = jsontext.AppendUnquote(t.key[:0], name)
		} else {
			t.key = strconv.AppendInt(t.key[:0], int64(i), 10)
		}

		child := n.child(t.key)
		if child == nil || (filter && !child.fields) {
			if err := t.dec.SkipValue(); err != nil {
				return 0, 0, nil, err
			}

			continue
		}

		childStart, childEnd, childValue, err := t.walk(child, filter)
		if err != nil {
			return 0, 0, nil, err
		}

		if filter {
			if !isObject {
				for nextIndex < i {
					if len(value) > 1 {
						value = append(value, ',')
					}
					value = append(value, "null"...)
					nextIndex++
				}
				nextIndex++
			}
			if len(value) > 1 {
				value = append(value, ',')
			}
			if isObject {
				value = append(append(value, name...), ':')
			}
			if childValue == nil {
				childValue = t.body[childStart:childEnd]
			}

			value = append(value, childValue...)
		} else if childValue != nil {
			value = append(value, t.body[last:childStart]...)
			value = append(value, childValue...)
			last = childEnd
		}
	}

	if _, err := t.dec.ReadToken(); err != nil {
		return 0, 0, nil, err
	}

	end := int(t.dec.InputOffset())

	switch {
	case filter:
		value = append(value, byte(closing))
	case value != nil:
		value = append(value, t.body[last:end]...)
	}

	return start, end, value, nil
}
