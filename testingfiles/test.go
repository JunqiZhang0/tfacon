package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

var DEFECT_TYPE map[string]string = map[string]string{
	"AutomationBug": "ab001",
	"ProductBug":    "pb001",
	"SystemIssue":   "si001",
}

// type innerList []map[string]interface{}

type Connector struct {
	launch_id    string
	project_name string
	auth_token   string
	rp_url       string
	client       *http.Client
	tfa_url      string
}

type TFAInput struct {
	Id       string `json:"id"`
	Project  string `json:"project"`
	Messages string `json:"messages"`
}

func main() {

	launch_id := "909"
	project_name := "TEFLO_RP"
	auth_token := "510256fa-8a43-4b6b-a0d2-c3388d9164a9"
	rp_url := "https://reportportal-ccit.apps.ocp4.prod.psi.redhat.com"
	tfa_url := "https://dave.corp.redhat.com:443/models/5f248eb11e43c7000602300b/latest/model"
	connector := NewConnector(auth_token, project_name, rp_url, launch_id, tfa_url)
	ids := connector.GetAllTestIds()
	updated_list_of_issues := connector.BuildUpdatedList(ids)
	// json, _ := json.Marshal(updated_list_of_issues)
	// fmt.Println(string(json))
	connector.UpdateAll(updated_list_of_issues)

}

func (c *Connector) UpdateAll(updated_list_of_issues UpdatedList) {
	if len(updated_list_of_issues.IssuesList) == 0 {
		return
	}
	json_updated_list_of_issues, _ := json.Marshal(updated_list_of_issues)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/%s/item", c.rp_url, c.project_name), bytes.NewBuffer((json_updated_list_of_issues)))
	if err != nil {
		fmt.Errorf("%s", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.auth_token))
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Errorf("update all failed in sending request: %v", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("This is the return info from update: %v\n", string(data))

}

type Issues []IssueItem

type IssueItem struct {
	Issue      IssueInfo `json:"issue"`
	TestItemId string    `json:"testItemId"`
}

type IssueInfo struct {
	IssueType            string        `json:"issueType"`
	Comment              string        `json:"comment"`
	AutoAnalyzed         bool          `json:"autoAnalyzed"`
	IgnoreAnalyzer       bool          `json:"ignoreAnalyzer"`
	ExternalSystemIssues []interface{} `json:"externalSystemIssues"`
}

type UpdatedList struct {
	IssuesList Issues `json:"issues"`
}

func (c *Connector) BuildUpdatedList(ids []string) UpdatedList {
	return UpdatedList{IssuesList: c.BuildIssues(ids)}
}

func (c *Connector) BuildIssues(ids []string) Issues {
	var issues Issues = Issues{}
	for _, id := range ids {
		fmt.Println(id)
		logs := c.GetTestLog(id)
		// Make logs to string(in []byte format)
		log_after_marshal, _ := json.Marshal(logs)
		// This can be the input of GetPrediction
		var tfa_input TFAInput = c.BuildTFAInput(id, string(log_after_marshal))
		prediction_json := c.GetPrediction(id, tfa_input)
		prediction := gjson.Get(prediction_json, "result.prediction").String()
		prediction_code := DEFECT_TYPE[prediction]
		var issue_info IssueInfo = c.GetIssueInfoForSingleTestId(id)
		issue_info.IssueType = prediction_code
		var issue_item IssueItem = IssueItem{Issue: issue_info, TestItemId: id}
		issues = append(issues, issue_item)
	}
	return issues
}

func (c *Connector) GetIssueInfoForSingleTestId(id string) IssueInfo {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.id=%s&filter.eq.launchId=%s&isLatest=false&launchesLimit=0", c.rp_url, c.project_name, id, c.launch_id), nil)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.auth_token))
	if err != nil {
		fmt.Errorf("request to get test ids failed: ", err)
	}
	resp, err := c.client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	issue_info_str := gjson.Get(string(data), "content.0.issue").String()
	var issue_info IssueInfo
	json.Unmarshal([]byte(issue_info_str), &issue_info)
	return issue_info

}

type TFAModel map[string]TFAInput

func (c *Connector) GetPrediction(id string, tfa_input TFAInput) string {
	tfa_model := TFAModel{"data": tfa_input}
	model, err := json.Marshal(tfa_model)
	if err != nil {
		fmt.Errorf("%s", err)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s", c.tfa_url), bytes.NewBuffer(model))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Errorf("request to get test ids failed: ", err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("Read response body failed:", err)
	}
	return string(data)
}

func (c *Connector) BuildTFAInput(test_id, messages string) TFAInput {
	return TFAInput{Id: test_id, Project: c.project_name, Messages: messages}
}

func NewConnector(auth_token, project_name, rp_url, launch_id, tfa_input string) *Connector {
	return &Connector{launch_id: launch_id, rp_url: rp_url, client: &http.Client{}, auth_token: auth_token, project_name: project_name, tfa_url: tfa_input}
}

func (c *Connector) GetAllTestIds() []string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.issueType=ti001&filter.eq.launchId=%s&filter.eq.status=FAILED&isLatest=false&launchesLimit=0", c.rp_url, c.project_name, c.launch_id), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.auth_token))
	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Errorf("request to get test ids failed: ", err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	a := gjson.Get(string(data), "content")
	var ret []string
	a.ForEach(func(_, m gjson.Result) bool {
		ret = append(ret, m.Get("id").String())
		return true
	})
	return ret
}

func (c *Connector) GetTestLog(test_id string) []string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/log?filter.eq.item=%s&filter.eq.launchId=%s", c.rp_url, c.project_name, test_id, c.launch_id), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.auth_token))
	if err != nil {
		panic(err)
	}
	resp, _ := c.client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	a := gjson.Get(string(data), "content")
	var ret []string
	a.ForEach(func(_, m gjson.Result) bool {
		ret = append(ret, m.Get("message").String())
		return true
	})
	return ret
}
