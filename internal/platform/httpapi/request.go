package httpapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func pathID(r *http.Request) string {
	p := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(p) > 2 {
		return p[2]
	}
	return ""
}
func decodeBody(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("request body required")
	}
	dec := json.NewDecoder(http.MaxBytesReader(nil, r.Body, 1<<20))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}
