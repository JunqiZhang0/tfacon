package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HTTPHandler struct {
	Client     *http.Client
	Method     string
	URL        string
	AUTH_TOKEN string
	Body       *bytes.Buffer
}

func (h *HTTPHandler) Do() []byte {
	req, err := http.NewRequest(h.Method, h.URL, h.Body)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", h.AUTH_TOKEN))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := h.Client.Do(req)
	if err != nil {
		panic(err)
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return d
}
