// common package contains all shared structs(data structures) required for all connectors
package common

type GeneralUpdatedList interface {
	GetSelf() GeneralUpdatedList
}

// type Issues []IssueItem

// type IssueItem struct {
// 	Issue      IssueInfo `json:"issue"`
// 	TestItemId string    `json:"testItemId"`
// }

// type IssueInfo struct {
// 	IssueType            string        `json:"issueType"`
// 	Comment              string        `json:"comment"`
// 	AutoAnalyzed         bool          `json:"autoAnalyzed"`
// 	IgnoreAnalyzer       bool          `json:"ignoreAnalyzer"`
// 	ExternalSystemIssues []interface{} `json:"externalSystemIssues"`
// }

// type UpdatedList struct {
// 	IssuesList Issues `json:"issues"`
// }

type TFAModel map[string]TFAInput
type TFAInput struct {
	Id       string `json:"id"`
	Project  string `json:"project"`
	Messages string `json:"messages"`
}

var DEFECT_TYPE map[string]string = map[string]string{
	"AutomationBug": "ab001",
	"ProductBug":    "pb001",
	"SystemIssue":   "si001",
}
