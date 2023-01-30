package ref

import (
	"encoding/json"
	"fmt"
	. "reflect"
	"strings"
)

type ref struct {
	ida  int `json:"test"`
	IDb  int
	Name string
	Con  complex128
	Fun  func(int) string
	ch   chan func()
}

func Ref() {
	r := ref{
		ida:  5,
		IDb:  6,
		Name: "asdf",
		Con:  complex(4, 1),
		Fun:  func(i int) string { return "hellp" },
		ch:   make(chan func()),
	}
	_ = r
	m := map[string]interface{}{
		"ida":  5,
		"IDb":  6,
		"Name": "asdf",
	}
	_ = m
	s := []interface{}{1, 2, 5, 6, 9, m}
	_ = s
	a := [4]int{6, 5, 4, 3}
	_ = a
	c := complex(1, 2)
	_ = c
	fmt.Printf("%v\n", structToMap(r))
	j, e := json.Marshal(r)
	fmt.Printf("%s\ne:%v\n", j, e)
}

func structToMap(data interface{}) interface{} {
	if data == nil {
		return data
	}
	val := ValueOf(data)
	ret := make(map[string]interface{}, 30)

	switch val.Kind() {
	case Struct:
		for i := 0; i < val.NumField(); i++ {
			t := val.Type()
			v := val.Field(i)
			k := val.Type().Field(i).Name
			tag := val.Type().Field(i).Tag.Get("json")

			if k == "Fun" || k == "ch" {
				fmt.Printf("type: %v  ", v.Type())
			}
			fmt.Printf("v:%v, k:%v, t:%v, tag:%v\n", v, k, t, decideTag(tag))
		}
	case Map:
		iter := val.MapRange()
		for iter.Next() {
			v := iter.Value()
			k := iter.Key()
			fmt.Printf("v:%v, k:%v\n", v, k)
		}
	case Slice, Array:
		for i := 0; i < val.Len(); i++ {
			v := val.Index(i)
			t := v.Type()
			fmt.Printf("v:%v, t:%v\n", v, t)
		}
	}

	return ret
}

func getVal(rv *Value) (interface{}, bool) {
	switch rv.Kind() {
	case Invalid:
		return "unknown", true
	case Bool:
		return rv.Bool(), true
	case Int8, Int16, Int32, Int64, Uint, Uint8, Uint16, Uint32, Uint64:
		return rv.Int(), true
	case Float32, Float64:
		return rv.Float(), true
	case Complex64, Complex128:
		return fmt.Sprintf("%v", rv.Complex()), true
	case Chan, Func:
		return rv.Type(), true

	}
	return "unknown", false
}

func decideTag(tag string) string {
	tags := strings.Split(tag, ",")
	for _, t := range tags {
		if t != "omitempty" && t != "-" {
			return t
		}
	}
	return ""
}

func iftest()
