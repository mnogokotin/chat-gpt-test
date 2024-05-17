package model

import (
	"fmt"
	"reflect"
)

type Event struct {
	ID   int
	Date string
}

type typeA struct {
	fieldA1 int
	fieldA2 string
}

type typeB struct {
	fieldB1 float32
	fieldB2 bool
}

func main() {

	a := []typeA{
		{10, "foo"},
		{20, "boo"},
	}
	b := []typeB{
		{2.5, true},
		{3.5, false},
	}

	printArrayAny(a)
	printArrayAny(b)
}

func printArrayAny(data interface{}) {
	v := reflect.ValueOf(data)
	for i := 0; i < v.Len(); i++ {
		fmt.Printf("printArrayAny row %d: %v\n", i, v.Index(i).Interface())
	}
}
