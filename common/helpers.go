package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

func getExistingDefectTypeLocatorId(gjson_obj []gjson.Result, defect_type string) (string, bool) {
	for _, v := range gjson_obj {
		defect_type_info := v.Map()
		if defect_type_info["longName"].String() == defect_type {
			return defect_type_info["locator"].String(), true
		}
	}
	return "", false
}

func InitDefectTypes() {
	var client *http.Client = &http.Client{}
	url := "https://reportportal-ccit.apps.ocp4.prod.psi.redhat.com/api/v1/JUNQI_RP/settings"
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("bearer %s", "510256fa-8a43-4b6b-a0d2-c3388d9164a9"))
	if err != nil {
		panic(err)
	}
	req.Header.Add("Content-Type", "application/json")
	resp, _ := client.Do(req)
	d, _ := ioutil.ReadAll(resp.Body)
	ti_sub := gjson.Get(string(d), "subTypes.TO_INVESTIGATE").Array()

	for _, v := range PREDICTED_SUB_TYPES {
		locator, ok := getExistingDefectTypeLocatorId(ti_sub, v["longName"])
		if !ok {
			d, _ := json.Marshal(v)
			req, err := http.NewRequest("POST", "https://reportportal-ccit.apps.ocp4.prod.psi.redhat.com/api/v1/JUNQI_RP/settings/sub-type", bytes.NewBuffer(d))
			req.Header.Add("Authorization", fmt.Sprintf("bearer %s", "510256fa-8a43-4b6b-a0d2-c3388d9164a9"))
			if err != nil {
				panic(err)
			}
			req.Header.Add("Content-Type", "application/json")
			resp, err := client.Do(req)
			if err != nil {
				panic(fmt.Errorf("request to get test ids failed: %s", err))
			}
			data, err := ioutil.ReadAll(resp.Body)
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
