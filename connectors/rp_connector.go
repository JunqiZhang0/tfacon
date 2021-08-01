package connectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"tfactl/common"

	"github.com/tidwall/gjson"
)

type RPConnector struct {
	launch_id    string
	project_name string
	auth_token   string
	rp_url       string
	client       *http.Client
	tfa_url      string
}

func (c *RPConnector) UpdateAll(updated_list_of_issues common.UpdatedList) {
	if len(updated_list_of_issues.IssuesList) == 0 {
		return
	}
	json_updated_list_of_issues, _ := json.Marshal(updated_list_of_issues)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/%s/item", c.rp_url, c.project_name), bytes.NewBuffer((json_updated_list_of_issues)))
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.auth_token))
	resp, err := c.client.Do(req)
	if err != nil {
		panic(fmt.Errorf("update all failed in sending request: %v", err))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("This is the return info from update: %v\n", string(data))

}

func (c *RPConnector) BuildUpdatedList(ids []string) common.UpdatedList {
	return common.UpdatedList{IssuesList: c.BuildIssues(ids)}
}

func (c *RPConnector) BuildIssues(ids []string) common.Issues {
	var issues common.Issues = common.Issues{}
	for _, id := range ids {
		fmt.Println(id)
		logs := c.GetTestLog(id)
		// Make logs to string(in []byte format)
		log_after_marshal, _ := json.Marshal(logs)
		// This can be the input of GetPrediction
		var tfa_input common.TFAInput = c.BuildTFAInput(id, string(log_after_marshal))
		prediction_json := c.GetPrediction(id, tfa_input)
		prediction := gjson.Get(prediction_json, "result.prediction").String()
		prediction_code := common.DEFECT_TYPE[prediction]
		var issue_info common.IssueInfo = c.GetIssueInfoForSingleTestId(id)
		issue_info.IssueType = prediction_code
		var issue_item common.IssueItem = common.IssueItem{Issue: issue_info, TestItemId: id}
		issues = append(issues, issue_item)
	}
	return issues
}

func (c *RPConnector) GetIssueInfoForSingleTestId(id string) common.IssueInfo {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.id=%s&filter.eq.launchId=%s&isLatest=false&launchesLimit=0", c.rp_url, c.project_name, id, c.launch_id), nil)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.auth_token))
	if err != nil {
		panic(fmt.Errorf("request to get test ids failed: %s", err))
	}
	resp, err := c.client.Do(req)
	if err != nil {
		panic(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	issue_info_str := gjson.Get(string(data), "content.0.issue").String()
	var issue_info common.IssueInfo
	json.Unmarshal([]byte(issue_info_str), &issue_info)
	return issue_info

}

func (c *RPConnector) GetPrediction(id string, tfa_input common.TFAInput) string {
	tfa_model := common.TFAModel{"data": tfa_input}
	model, err := json.Marshal(tfa_model)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	req, err := http.NewRequest("POST", c.tfa_url, bytes.NewBuffer(model))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		panic(fmt.Errorf("request to get test ids failed: %s", err))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("read response body failed: %v", err))
	}
	return string(data)
}

func (c *RPConnector) BuildTFAInput(test_id, messages string) common.TFAInput {
	return common.TFAInput{Id: test_id, Project: c.project_name, Messages: messages}
}

func NewConnector(auth_token, project_name, rp_url, launch_id, tfa_input string) *RPConnector {
	return &RPConnector{launch_id: launch_id, rp_url: rp_url, client: &http.Client{}, auth_token: auth_token, project_name: project_name, tfa_url: tfa_input}
}

func (c *RPConnector) GetAllTestIds() []string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.issueType=ti001&filter.eq.launchId=%s&filter.eq.status=FAILED&isLatest=false&launchesLimit=0", c.rp_url, c.project_name, c.launch_id), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.auth_token))
	resp, err := c.client.Do(req)
	if err != nil {
		panic(fmt.Errorf("request to get test ids failed: %s", err))
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

func (c *RPConnector) GetTestLog(test_id string) []string {
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
