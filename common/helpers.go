package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func SendHTTPRequest(method, url, auth_token string, body *bytes.Buffer, client *http.Client) ([]byte, error, bool) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", auth_token))
	if err != nil {
		return nil, err, false
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err, false
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err, false
	}
	if resp.StatusCode == 200 {
		return d, err, true
	} else {
		return d, err, false
	}
}
