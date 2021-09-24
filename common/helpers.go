package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

func PrintGreen(str string) {
	color.Green(str)
}

func PrintRed(str string) {
	color.Red(str)
}

func PrintHeader(version string) {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("tfacon  %s\n", version)
	fmt.Println("Copyright (C) 2021, Red Hat, Inc.")
	fmt.Print("-------------------------------------------------\n\n\n")
}

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
		if method == "POST" {
			if resp.StatusCode == 201 {
				return d, err, true
			}
		} else {
			fmt.Printf("status code is:%v\n", resp.StatusCode)
			return d, err, false
		}
	}
	return d, err, false
}
