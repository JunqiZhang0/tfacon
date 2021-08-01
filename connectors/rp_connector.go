package connectors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"tfacon/common"

	"github.com/tidwall/gjson"
)

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

func (u UpdatedList) GetSelf() common.GeneralUpdatedList {
	return u
}

type RPConnector struct {
	LaunchId    string `mapstructure:"LAUNCH_ID"`
	ProjectName string `mapstructure:"PROJECT_NAME"`
	AuthToken   string `mapstructure:"AUTH_TOKEN"`
	RPURL       string `mapstructure:"PLATFORM_URL"`
	Client      *http.Client
	TFAURL      string `mapstructure:"TFA_URL"`
}

func (c *RPConnector) UpdateAll(updated_list_of_issues common.GeneralUpdatedList) {
	if len(updated_list_of_issues.GetSelf().(UpdatedList).IssuesList) == 0 {
		return
	}
	json_updated_list_of_issues, _ := json.Marshal(updated_list_of_issues)
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/%s/item", c.RPURL, c.ProjectName), bytes.NewBuffer((json_updated_list_of_issues)))
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	resp, err := c.Client.Do(req)
	if err != nil {
		panic(fmt.Errorf("update all failed in sending request: %v", err))
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("This is the return info from update: %v\n", string(data))

}

func (c *RPConnector) BuildUpdatedList(ids []string) common.GeneralUpdatedList {
	return UpdatedList{IssuesList: c.BuildIssues(ids)}
}

func (c *RPConnector) BuildIssues(ids []string) Issues {
	var issues Issues = Issues{}
	for _, id := range ids {
		log.Printf("Getting prediction of test item(id): %s\n", id)
		logs := c.GetTestLog(id)
		// Make logs to string(in []byte format)
		log_after_marshal, _ := json.Marshal(logs)
		// This can be the input of GetPrediction
		var tfa_input common.TFAInput = c.BuildTFAInput(id, string(log_after_marshal))
		prediction_json := c.GetPrediction(id, tfa_input)
		prediction := gjson.Get(prediction_json, "result.prediction").String()
		prediction_code := common.DEFECT_TYPE[prediction]
		var issue_info IssueInfo = c.GetIssueInfoForSingleTestId(id)
		issue_info.IssueType = prediction_code
		var issue_item IssueItem = IssueItem{Issue: issue_info, TestItemId: id}
		issues = append(issues, issue_item)
	}
	return issues
}

func (c *RPConnector) GetIssueInfoForSingleTestId(id string) IssueInfo {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.id=%s&filter.eq.launchId=%s&isLatest=false&launchesLimit=0", c.RPURL, c.ProjectName, id, c.LaunchId), nil)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	if err != nil {
		panic(fmt.Errorf("request to get test ids failed: %s", err))
	}
	resp, err := c.Client.Do(req)
	if err != nil {
		panic(err)
	}
	data, _ := ioutil.ReadAll(resp.Body)
	issue_info_str := gjson.Get(string(data), "content.0.issue").String()
	var issue_info IssueInfo
	json.Unmarshal([]byte(issue_info_str), &issue_info)
	return issue_info

}

func (c *RPConnector) GetPrediction(id string, tfa_input common.TFAInput) string {
	tfa_model := common.TFAModel{"data": tfa_input}
	model, err := json.Marshal(tfa_model)
	if err != nil {
		panic(fmt.Errorf("%s", err))
	}
	req, err := http.NewRequest("POST", c.TFAURL, bytes.NewBuffer(model))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := c.Client.Do(req)
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
	return common.TFAInput{Id: test_id, Project: c.ProjectName, Messages: messages}
}

func (c *RPConnector) GetAllTestIds() []string {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.issueType=ti001&filter.eq.launchId=%s&filter.eq.status=FAILED&isLatest=false&launchesLimit=0", c.RPURL, c.ProjectName, c.LaunchId), nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	resp, err := c.Client.Do(req)
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
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/log?filter.eq.item=%s&filter.eq.launchId=%s", c.RPURL, c.ProjectName, test_id, c.LaunchId), nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	if err != nil {
		panic(err)
	}
	resp, _ := c.Client.Do(req)
	data, _ := ioutil.ReadAll(resp.Body)
	a := gjson.Get(string(data), "content")
	var ret []string
	a.ForEach(func(_, m gjson.Result) bool {
		ret = append(ret, m.Get("message").String())
		return true
	})
	return ret
}