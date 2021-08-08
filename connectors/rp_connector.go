package connectors

import (
	"bytes"
	"encoding/json"

	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/JunqiZhang0/tfacon/common"
	"github.com/pkg/errors"

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

func (c *RPConnector) Validate() error {
	return errors.New("")
}

func (c RPConnector) String() string {
	v := reflect.ValueOf(c)
	typeOfS := v.Type()
	str := ""
	for i := 0; i < v.NumField(); i++ {
		str = str + fmt.Sprintf("%s: \t %v\n", typeOfS.Field(i).Name, v.Field(i).Interface())
	}
	return str
}

func (c *RPConnector) UpdateAll(updated_list_of_issues common.GeneralUpdatedList) {
	if len(updated_list_of_issues.GetSelf().(UpdatedList).IssuesList) == 0 {
		return
	}
	json_updated_list_of_issues, _ := json.Marshal(updated_list_of_issues)
	// fmt.Println(string(json_updated_list_of_issues))
	// req, err := http.NewRequest("PUT", fmt.Sprintf("%s/api/v1/%s/item", c.RPURL, c.ProjectName), bytes.NewBuffer((json_updated_list_of_issues)))
	// if err != nil {
	// 	panic(fmt.Errorf("%s", err))
	// }
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Content-Type", "application/json")
	// req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	// resp, err := c.Client.Do(req)
	// if err != nil {
	// 	panic(fmt.Errorf("update all failed in sending request: %v", err))
	// }
	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	log.Println("Updating All Test Items With Predictions...")
	url := fmt.Sprintf("%s/api/v1/%s/item", c.RPURL, c.ProjectName)
	method := "PUT"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(json_updated_list_of_issues)
	data, err, success := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(fmt.Sprintf("Updated All failed: %s", err))
	}
	if success {
		fmt.Println("Updated All Test Items Successfully!")
	} else {
		fmt.Println("Updated Failed!")
	}
	fmt.Printf("This is the return info from update: %v\n", string(data))

}

func (c *RPConnector) BuildUpdatedList(ids []string, concurrent bool) common.GeneralUpdatedList {

	return UpdatedList{IssuesList: c.BuildIssues(ids, concurrent)}

}

func (c *RPConnector) BuildIssues(ids []string, concurrent bool) Issues {
	var issues Issues = Issues{}
	if concurrent {
		return c.BuildIssuesConcurrent(ids)
	} else {
		for _, id := range ids {
			issues = append(issues, c.BuildIssueItemHelper(id))
			log.Printf("Getting prediction of test item(id): %s\n", id)
		}
		return issues
	}
}

func (c *RPConnector) BuildIssuesConcurrent(ids []string) Issues {
	var issues Issues = Issues{}
	var issuesChan chan IssueItem = make(chan IssueItem, len(ids))
	var idsChan chan string = make(chan string, len(ids))
	var exitChan chan bool = make(chan bool, len(ids))
	go func() {
		for _, id := range ids {
			idsChan <- id
		}
		close(idsChan)
	}()
	for i := 0; i < len(ids); i++ {
		go c.BuildIssueItemConcurrent(issuesChan, idsChan, exitChan)
	}
	for i := 0; i < len(ids); i++ {
		<-exitChan
	}
	close(issuesChan)
	for issue := range issuesChan {
		issues = append(issues, issue)
	}
	return issues
}

func (c *RPConnector) BuildIssueItemHelper(id string) IssueItem {
	logs := c.GetTestLog(id)
	// Make logs to string(in []byte format)
	log_after_marshal, _ := json.Marshal(logs)
	// This can be the input of GetPrediction
	var tfa_input common.TFAInput = c.BuildTFAInput(id, string(log_after_marshal))
	prediction_json := c.GetPrediction(id, tfa_input)
	prediction := gjson.Get(prediction_json, "result.prediction").String()
	// prediction_code := common.DEFECT_TYPE[prediction]
	prediction_code := common.TFA_DEFECT_TYPE_TO_SUB_TYPE[prediction]["locator"]
	// fmt.Println(prediction_code)
	var issue_info IssueInfo = c.GetIssueInfoForSingleTestId(id)
	issue_info.IssueType = prediction_code
	var issue_item IssueItem = IssueItem{Issue: issue_info, TestItemId: id}
	return issue_item
}

func (c *RPConnector) BuildIssueItemConcurrent(issuesChan chan<- IssueItem, idsChan <-chan string, exitChan chan<- bool) {
	for {
		id, ok := <-idsChan
		if !ok {
			break
		}

		issuesChan <- c.BuildIssueItemHelper(id)
		log.Printf("Getting prediction of test item(id): %s\n", id)
	}
	exitChan <- true
}

func (c *RPConnector) GetIssueInfoForSingleTestId(id string) IssueInfo {
	// req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.id=%s&filter.eq.launchId=%s&isLatest=false&launchesLimit=0", c.RPURL, c.ProjectName, id, c.LaunchId), nil)
	// if err != nil {
	// 	panic(fmt.Errorf("%s", err))
	// }
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	// if err != nil {
	// 	panic(fmt.Errorf("request to get test ids failed: %s", err))
	// }
	// resp, err := c.Client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// data, _ := ioutil.ReadAll(resp.Body)
	url := fmt.Sprintf("%s/api/v1/%s/item?filter.eq.id=%s&filter.eq.launchId=%s&isLatest=false&launchesLimit=0", c.RPURL, c.ProjectName, id, c.LaunchId)
	method := "GET"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, err, _ := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(err)
	}
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
	// req, err := http.NewRequest("POST", c.TFAURL, bytes.NewBuffer(model))
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Add("Content-Type", "application/json")
	// resp, err := c.Client.Do(req)
	// if err != nil {
	// 	panic(fmt.Errorf("request to get test ids failed: %s", err))
	// }
	// data, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(fmt.Errorf("read response body failed: %v", err))
	// }
	url := c.TFAURL
	method := "POST"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(model)
	data, err, _ := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(err)
	}
	return string(data)
}

func (c *RPConnector) BuildTFAInput(test_id, messages string) common.TFAInput {
	return common.TFAInput{Id: test_id, Project: c.ProjectName, Messages: messages}
}

func (c *RPConnector) GetAllTestIds() []string {
	// req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/item?filter.eq.issueType=ti001&filter.eq.launchId=%s&filter.eq.status=FAILED&isLatest=false&launchesLimit=0", c.RPURL, c.ProjectName, c.LaunchId), nil)
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	// resp, err := c.Client.Do(req)
	// if err != nil {
	// 	panic(fmt.Errorf("request to get test ids failed: %s", err))
	// }
	// data, _ := ioutil.ReadAll(resp.Body)
	url := fmt.Sprintf("%s/api/v1/%s/item?filter.eq.issueType=ti001&filter.eq.launchId=%s&filter.eq.status=FAILED&isLatest=false&launchesLimit=0", c.RPURL, c.ProjectName, c.LaunchId)
	method := "GET"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, err, _ := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(err)
	}
	a := gjson.Get(string(data), "content")
	var ret []string
	a.ForEach(func(_, m gjson.Result) bool {
		ret = append(ret, m.Get("id").String())
		return true
	})
	return ret
}

func (c *RPConnector) GetTestLog(test_id string) []string {
	// req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/log?filter.eq.item=%s&filter.eq.launchId=%s", c.RPURL, c.ProjectName, test_id, c.LaunchId), nil)
	// req.Header.Add("Accept", "application/json")
	// req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	// if err != nil {
	// 	panic(err)
	// }
	// resp, _ := c.Client.Do(req)
	// data, _ := ioutil.ReadAll(resp.Body)
	url := fmt.Sprintf("%s/api/v1/%s/log?filter.eq.item=%s&filter.eq.launchId=%s", c.RPURL, c.ProjectName, test_id, c.LaunchId)
	method := "GET"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, err, _ := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(err)
	}
	a := gjson.Get(string(data), "content")
	var ret []string
	a.ForEach(func(_, m gjson.Result) bool {
		ret = append(ret, m.Get("message").String())
		return true
	})
	return ret
}

func getExistingDefectTypeLocatorId(gjson_obj []gjson.Result, defect_type string) (string, bool) {
	for _, v := range gjson_obj {
		defect_type_info := v.Map()
		if defect_type_info["longName"].String() == defect_type {
			return defect_type_info["locator"].String(), true
		}
	}
	return "", false
}

func (c *RPConnector) InitConnector() {
	fmt.Println("Initializing defect types...")
	// req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v1/%s/settings", c.RPURL, c.ProjectName), nil)
	// req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
	// if err != nil {
	// 	panic(err)
	// }
	// req.Header.Add("Content-Type", "application/json")
	// resp, _ := c.Client.Do(req)
	// d, _ := ioutil.ReadAll(resp.Body)
	url := fmt.Sprintf("%s/api/v1/%s/settings", c.RPURL, c.ProjectName)
	method := "GET"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, err, _ := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(err)
	}
	ti_sub := gjson.Get(string(data), "subTypes.TO_INVESTIGATE").Array()

	for _, v := range common.PREDICTED_SUB_TYPES {
		locator, ok := getExistingDefectTypeLocatorId(ti_sub, v["longName"])
		if !ok {
			d, _ := json.Marshal(v)
			// req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v1/%s/settings/sub-type", c.RPURL, c.ProjectName), bytes.NewBuffer(d))
			// req.Header.Add("Authorization", fmt.Sprintf("bearer %s", c.AuthToken))
			// if err != nil {
			// 	panic(err)
			// }
			// req.Header.Add("Content-Type", "application/json")
			// resp, err := c.Client.Do(req)
			// if err != nil {
			// 	panic(fmt.Errorf("request to get test ids failed: %s", err))
			// }
			// data, err := ioutil.ReadAll(resp.Body)
			url := fmt.Sprintf("%s/api/v1/%s/settings/sub-type", c.RPURL, c.ProjectName)
			method := "POST"
			auth_token := c.AuthToken
			body := bytes.NewBuffer(d)
			data, err, _ := common.SendHTTPRequest(method, url, auth_token, body, c.Client)

			if err != nil {
				panic(fmt.Errorf("read response body failed: %v", err))
			}
			v["locator"] = gjson.Get(string(data), "locator").String()
		} else {
			v["locator"] = locator
		}
	}
	// fmt.Println(PREDICTED_SUB_TYPES)
}
