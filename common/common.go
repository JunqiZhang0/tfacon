// common package contains all shared structs(data structures) required for all connectors
package common

// GeneralUpdatedList is an updated list of object, with the prediction from TFA classifier
// each connector should have it's own UpdatedList structure and implement the
// GeneralUpdatedList interface
type GeneralUpdatedList interface {
	GetSelf() GeneralUpdatedList
}

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
