package httpclient

import (
	"net/http"
	"strings"
)

var sensitiveHeaders = map[string]bool{
	"authorization": true,
}

// NonSensitiveHeaders TODOC
func NonSensitiveHeaders(h http.Header) (h2 http.Header) {
	h2 = http.Header{}

	for k, vs := range h {
		if sensitiveHeaders[strings.ToLower(k)] {
			h2.Set(k, "<redacted>")
			continue
		}

		for _, v := range vs {
			h2.Add(k, v)
		}
	}

	return
}
