package main

import (
	"fmt"
	"reflect"
)

func detectType(v interface{}) string {
	switch v.(type) {
	case int:
		return "int"
	case string:
		return "string"
	case bool:
		return "bool"
	case chan interface{}:
		return "chan"
	default:
		return reflect.TypeOf(v).String()
	}
}

func main() {
	var (
		num  int           = 42
		str  string        = "hello"
		flag bool          = true
		ch   chan struct{} = make(chan struct{})
	)

	fmt.Println("num is:", detectType(num))   // int
	fmt.Println("str is:", detectType(str))   // string
	fmt.Println("flag is:", detectType(flag)) // bool
	fmt.Println("ch is:", detectType(ch))     // chan

}
