package main

import "fmt"

// type A interface {
// 	print()
// }

// type B struct {
// 	Test string `json:"test"`
// }

// func (b *B) print() {
// 	fmt.Println("hello")
// }

// // func New(test ...string) *B {
// // 	return &B{Test: test}
// // }

// var asd string = "{\"test\":\"hello\"}"

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

func (u UpdatedList) GetSelf() GenralUpdatedList {
	return u
}

type GenralUpdatedList interface {
	GetSelf() GenralUpdatedList
}

func main() {
	var a GenralUpdatedList
	a = UpdatedList{}.GetSelf()
	fmt.Printf("%T\n", a)
	// var a A = &B{}
	// // var a A
	// json.Unmarshal([]byte(asd), a)
	// // a.print()
	// fmt.Printf("%+v:\n", a)
}
