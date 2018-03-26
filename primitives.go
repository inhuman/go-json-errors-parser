package go_json_errors_parser

import (
	"encoding/json"
	"fmt"
)

type ParsedErrorInterface interface {
	unmarshalJson() error
	transferTo(ps *ParsedErrors, parent string)
}

// String error struct and unmarshal
type stringError struct {
	Error      string
	RawMessage json.RawMessage
}

func (e *stringError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

func (e *stringError) transferTo(ps *ParsedErrors, parent string) {
	r := ParsedError{}
	r.Messages = append(r.Messages, trimQuotes(e.Error))
	r.Parent = parent
	ps.ParsedErrors = append(ps.ParsedErrors, r)
}

// String slice error struct and unmarshal
type sliceStringError struct {
	Error      []string
	RawMessage json.RawMessage
}

func (e *sliceStringError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

func (e *sliceStringError) transferTo(ps *ParsedErrors, parent string) {
	r := ParsedError{}
	r.Messages = append(r.Messages, e.Error...)
	r.Parent = parent
	ps.ParsedErrors = append(ps.ParsedErrors, r)
}

// Map of string slice interface error struct and unmarshal
type mapStringSliceInterfaceError struct {
	Error      map[string][]interface{}
	RawMessage json.RawMessage
}

func (e *mapStringSliceInterfaceError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

func (e *mapStringSliceInterfaceError) transferTo(ps *ParsedErrors, parent string) {
	var tmpMap map[string][]interface{}

	err := json.Unmarshal(e.RawMessage, &tmpMap)
	if err != nil {
		panic(err)
	}

	r := ParsedError{}

	formattedStrs := make(map[string][]string)
	for name, str := range tmpMap {
		var tmpMap []string
		for _, item := range str {
			tmpMap = append(tmpMap, fmt.Sprintf("%v", item))
		}
		formattedStrs[name] = tmpMap
	}
	r.Children = formattedStrs
	r.Parent = parent
	ps.ParsedErrors = append(ps.ParsedErrors, r)

}

//Slice of map string interface error struct and unmarshal
type sliceMapStringInterfaceError struct {
	Error      []map[string]interface{}
	RawMessage json.RawMessage
}

func (e *sliceMapStringInterfaceError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

func (e *sliceMapStringInterfaceError) transferTo(ps *ParsedErrors, parent string) {

	var strs []map[string]interface{}
	err := json.Unmarshal(e.RawMessage, &strs)
	if err != nil {
		panic(err)
	}
	r := ParsedError{}

	for _, value := range strs {
		tmp := make(map[string][]string)
		for key, val := range value {
			tmp[key] = []string{fmt.Sprintf("%v", val)}
		}
		r.Children = tmp
		r.Parent = parent
	}
	ps.ParsedErrors = append(ps.ParsedErrors, r)
}

// Bool struct and unmarshal
type boolValue struct {
	Value      bool
	RawMessage json.RawMessage
}

func (e *boolValue) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Value)
}

func (e *boolValue) transferTo(ps *ParsedErrors, parent string) {}


// Num struct and unmarshal
type numValue struct {
	Value int
	RawMessage json.RawMessage
}

func (e *numValue) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Value)
}

func (e *numValue) transferTo(ps *ParsedErrors, parent string) {}

//TODO: implement batch check
//var typeRegistry = make(map[string]reflect.Type)
//
//func makeInstance(name string) ParsedErrorInterface {
//	v := reflect.New(typeRegistry[name]).Elem()
//
//	ins := v.Interface()
//
//	return ins.(ParsedErrorInterface)
//}
//
//func batchCheck(s *json.RawMessage, ps *ParsedErrors, parent string) error {
//	typeRegistry["stringError"] = reflect.TypeOf(stringError{})
//
//	for name := range typeRegistry {
//		parErr := makeInstance(name)
//
//
//		err := parErr.unmarshalJson()
//		if err == nil {
//			parErr.transferTo(ps, parent)
//		} else {
//			return nil
//		}
//	}
//
//	return nil
//}
