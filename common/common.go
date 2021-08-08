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
	"SystemBug":     "si001",
}

var DEFECT_TYPE_TEST map[string]string = map[string]string{
	"AutomationBug": "",
	"ProductBug":    "",
	"SystemBug":     "",
}

type PREDICTED_SUB_TYPE map[string]string

var PREDICTED_SUB_TYPES map[string]PREDICTED_SUB_TYPE = map[string]PREDICTED_SUB_TYPE{
	"PREDICTED_AUTOMATION_BUG": PREDICTED_AUTOMATION_BUG,
	"PREDICTED_SYSTEM_BUG":     PREDICTED_SYSTEM_BUG,
	"PREDICTED_PRODUCT_BUG":    PREDICTED_PRODUCT_BUG,
}

var PREDICTED_AUTOMATION_BUG PREDICTED_SUB_TYPE = PREDICTED_SUB_TYPE{
	"typeRef":   "TO_INVESTIGATE",
	"longName":  "Predicted Automation Bug",
	"shortName": "TIA",
	"color":     "#ffeeaa",
}

var PREDICTED_SYSTEM_BUG PREDICTED_SUB_TYPE = PREDICTED_SUB_TYPE{
	"typeRef":   "TO_INVESTIGATE",
	"longName":  "Predicted System Issue",
	"shortName": "TIS",
	"color":     "#aaaaff",
}

var PREDICTED_PRODUCT_BUG PREDICTED_SUB_TYPE = PREDICTED_SUB_TYPE{
	"typeRef":   "TO_INVESTIGATE",
	"longName":  "Predicted Product Bug",
	"shortName": "TIP",
	"color":     "#ffaaaa",
}
