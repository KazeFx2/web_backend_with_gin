package Net

import (
	"bytes"
	"net/http"
)

func Request(method string, url string, headers map[string]string, buf *bytes.Buffer) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, buf)
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return http.DefaultClient.Do(req)
}

func Get(url string, headers map[string]string, params map[string]string) (*http.Response, error) {
	if len(params) > 0 {
		url += "?"
		for k, v := range params {
			url += k + "=" + v + "&"
		}
	}
	return Request("GET", url, headers, &bytes.Buffer{})
}
