package main

import (
	"encoding/json"
	"fmt"
)

type A interface {
	print()
}

type B struct {
	Test string `json:"test"`
}

func (b *B) print() {
	fmt.Println("hello")
}

// func New(test ...string) *B {
// 	return &B{Test: test}
// }

var asd string = "{\"test\":\"hello\"}"

func main() {
	var a A = &B{}
	// var a A
	json.Unmarshal([]byte(asd), a)
	// a.print()
	fmt.Printf("%+v:\n", a)
}
