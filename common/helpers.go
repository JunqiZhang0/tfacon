package common

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/fatih/color"
)

// PrintGreen is a helper function that
// prints str in green to terminal.
func PrintGreen(str string) {
	color.Green(str)
}

// PrintRed is a helper function that
// prints str in red to terminal.
func PrintRed(str string) {
	color.Red(str)
}

// PrintHeader is a helper function
// for the whole program to print
// header information.
func PrintHeader(version string) {
	fmt.Println("--------------------------------------------------")
	fmt.Printf("tfacon  %s\n", version)
	fmt.Println("Copyright (C) 2021, Red Hat, Inc.")
	fmt.Print("-------------------------------------------------\n\n\n")
}

// SendHTTPRequest is a helper function that
// deals with all http operation for tfacon.
func SendHTTPRequest(method, url, auth_token string, body *bytes.Buffer, client *http.Client) ([]byte, bool, error) {
	req, err := http.NewRequest(method, url, body)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", auth_token))
	if err != nil {
		return nil, false, err
	}

	req.Header.Add("Content-Type", "application/json")

	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, false, err
	}
	d, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}
	if resp.StatusCode == 200 {
		return d, true, err
	}
	if method == "POST" && resp.StatusCode == 201 {
		return d, true, err
	}

	fmt.Printf("status code is:%v\n", resp.StatusCode)
	return d, false, err
}

// HandleError is the Error handler
// for the whole tfacon.
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
