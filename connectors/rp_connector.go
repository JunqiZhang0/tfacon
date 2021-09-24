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

func (c *RPConnector) Validate(verbose bool) (bool, error) {
	fmt.Print("Validating....\n")
	_, validateRPURLAndAuthToken, err := c.validateRPURLAndAuthToken()

	if err != nil {
		err = errors.Errorf("%s", err)
		return false, err
	}
	if verbose {
		fmt.Printf("RPURLValidate: %t\n", validateRPURLAndAuthToken)
	}
	validateTFA, err := c.validateTFAURL()
	if verbose {
		fmt.Printf("TFAURLValidate: %t\n", validateTFA)
	}
	if err != nil {
		err = errors.Errorf("%s", err)
		return false, err
	}
	projectnameNotEmpty := c.ProjectName != ""
	if verbose {
		fmt.Printf("projectnameValidate: %t\n", projectnameNotEmpty)
	}
	if !projectnameNotEmpty {
		err = errors.Errorf("%s", "You need to input project name")
		return false, err
	}
	launchidNotEmpty := c.LaunchId != ""
	if verbose {
		fmt.Printf("lauchidValidate: %t\n", launchidNotEmpty)
	}
	if !launchidNotEmpty {
		err = errors.Errorf("%s", "You need to input launch id")
		return false, err
	}
	ret := validateRPURLAndAuthToken && validateTFA && projectnameNotEmpty && launchidNotEmpty
	return ret, nil
}

func (c *RPConnector) validateTFAURL() (bool, error) {
	body := `{"data": {"id": "123", "project": "rhv", "messages": ""}}`
	_, err, success := common.SendHTTPRequest("POST", c.TFAURL, "", bytes.NewBuffer([]byte(body)), c.Client)
	return success, err
}

func (c *RPConnector) validateRPURLAndAuthToken() ([]byte, bool, error) {
	data, err, success := common.SendHTTPRequest("GET", c.RPURL+"/api/v1/project/list", c.AuthToken, bytes.NewBuffer(nil), c.Client)
	return data, success, err
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

func (c *RPConnector) UpdateAll(updated_list_of_issues common.GeneralUpdatedList, verbose bool) {
	if len(updated_list_of_issues.GetSelf().(UpdatedList).IssuesList) == 0 {
		return
	}
	json_updated_list_of_issues, _ := json.Marshal(updated_list_of_issues)
	log.Println("Updating All Test Items With Predictions...")
	url := fmt.Sprintf("%s/api/v1/%s/item", c.RPURL, c.ProjectName)
	method := "PUT"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(json_updated_list_of_issues)
	data, err, success := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(fmt.Sprintf("Updated All failed: %s", err))
	}
	if verbose {
		fmt.Printf("This is the return info from update: %v\n", string(data))
	}

	if success {
		fmt.Println()
		common.PrintGreen("Updated All Test Items Successfully!")
	} else {
		common.PrintRed("Updated Failed!")
	}

}

func (c *RPConnector) BuildUpdatedList(ids []string, concurrent bool, add_attributes bool) common.GeneralUpdatedList {
	return UpdatedList{IssuesList: c.BuildIssues(ids, concurrent, add_attributes)}

}

func (c *RPConnector) BuildIssues(ids []string, concurrent bool, add_attributes bool) Issues {
	var issues Issues = Issues{}
	if concurrent {
		return c.BuildIssuesConcurrent(ids, add_attributes)
	} else {
		for _, id := range ids {
			issues = append(issues, c.BuildIssueItemHelper(id, add_attributes))
			log.Printf("Getting prediction of test item(id): %s\n", id)
		}
		return issues
	}
}

func (c *RPConnector) BuildIssuesConcurrent(ids []string, add_attributes bool) Issues {
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
		go c.BuildIssueItemConcurrent(issuesChan, idsChan, exitChan, add_attributes)
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

func (c *RPConnector) BuildIssueItemHelper(id string, add_attributes bool) IssueItem {
	logs := c.GetTestLog(id)
	// Make logs to string(in []byte format)
	log_after_marshal, _ := json.Marshal(logs)
	// This can be the input of GetPrediction
	testlog := string(log_after_marshal)
	var tfa_input common.TFAInput = c.BuildTFAInput(id, testlog)
	prediction_json := c.GetPrediction(id, tfa_input)
	prediction := gjson.Get(prediction_json, "result.prediction").String()
	// prediction_code := common.DEFECT_TYPE[prediction]
	prediction_code := common.TFA_DEFECT_TYPE_TO_SUB_TYPE[prediction]["locator"]
	// fmt.Println(prediction_code)
	var issue_info IssueInfo = c.GetIssueInfoForSingleTestId(id)
	issue_info.IssueType = prediction_code
	var issue_item IssueItem = IssueItem{Issue: issue_info, TestItemId: id}
	if add_attributes {
		prediction_name := common.TFA_DEFECT_TYPE_TO_SUB_TYPE[prediction]["longName"]
		err := c.updateAttributesForPrediction(id, prediction_name)
		if err != nil {
			panic(err)
		}
	}
	return issue_item
}

func (c *RPConnector) BuildIssueItemConcurrent(issuesChan chan<- IssueItem, idsChan <-chan string, exitChan chan<- bool, add_attributes bool) {
	for {
		id, ok := <-idsChan
		if !ok {
			break
		}

		issuesChan <- c.BuildIssueItemHelper(id, add_attributes)
		log.Printf("Getting prediction of test item(id): %s\n", id)
	}
	exitChan <- true
}

func (c *RPConnector) GetIssueInfoForSingleTestId(id string) IssueInfo {
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

type attribute map[string]string

func (c *RPConnector) updateAttributesForPrediction(id, prediction string) error {
	updated_attribute := map[string][]attribute{
		"attributes": {
			attribute{
				"key":   "AI Prediction",
				"value": prediction},
		},
	}
	url := fmt.Sprintf("%s/api/v1/%s/item/%s/update", c.RPURL, c.ProjectName, id)
	method := "PUT"
	auth_token := c.AuthToken
	d, err := json.Marshal(updated_attribute)
	if err != nil {
		panic(err)
	}
	body := bytes.NewBuffer(d)
	_, err, _ = common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	// fmt.Printf("This is the return from updating attributes: %s\n", string(data))
	log.Printf("Updated the test item(id): %s with it's prediction %s\n", id, prediction)
	return err

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
	fmt.Println("Initializing Defect Types...")
	url := fmt.Sprintf("%s/api/v1/%s/settings", c.RPURL, c.ProjectName)
	method := "GET"
	auth_token := c.AuthToken
	body := bytes.NewBuffer(nil)
	data, err, success := common.SendHTTPRequest(method, url, auth_token, body, c.Client)
	if err != nil {
		panic(err)
	}
	if !success {
		panic(fmt.Errorf("created defect types failed, please use superadmin auth_token"))
	}
	ti_sub := gjson.Get(string(data), "subTypes.TO_INVESTIGATE").Array()

	for _, sub_type := range common.PREDICTED_SUB_TYPES {

		locator, ok := getExistingDefectTypeLocatorId(ti_sub, sub_type["longName"])
		if !ok {
			d, _ := json.Marshal(sub_type)
			url := fmt.Sprintf("%s/api/v1/%s/settings/sub-type", c.RPURL, c.ProjectName)
			method := "POST"
			auth_token := c.AuthToken
			body := bytes.NewBuffer(d)
			data, err, success := common.SendHTTPRequest(method, url, auth_token, body, c.Client)

			if err != nil {
				panic(fmt.Errorf("read response body failed: %v", err))
			}
			if !success {
				panic(fmt.Errorf("creation of the defect type failed %v", string(data)))
			}
			sub_type["locator"] = gjson.Get(string(data), "locator").String()
		} else {
			sub_type["locator"] = locator
		}
	}
}
