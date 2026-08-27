package httpapi

import (
	"net/url"
	"strings"
)

func pathParts(path, prefix string) []string {
	p := strings.TrimPrefix(path, prefix)
	p = strings.Trim(p, "/")
	if p == "" {
		return nil
	}
	raw := strings.Split(p, "/")
	out := make([]string, 0, len(raw))
	for _, x := range raw {
		if y, e := url.PathUnescape(x); e == nil && y != "" {
			out = append(out, y)
		}
	}
	return out
}
func validID(id string) bool {
	if len(id) < 4 || len(id) > 128 {
		return false
	}
	return !strings.ContainsAny(id, "/\\")
}
func queryBool(v string) bool {
	return v == "1" || strings.EqualFold(v, "true") || strings.EqualFold(v, "yes")
}
