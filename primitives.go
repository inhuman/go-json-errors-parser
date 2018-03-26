package go_json_errors_parser

import (
	"encoding/json"
	"fmt"
)

type ParsedErrorInterface interface {
	setRawMessage(m json.RawMessage)
	unmarshalJson() error
	transferTo(ps *ParsedErrors, parent string)
}

// String error struct and unmarshal
type stringError struct {
	Error      string
	RawMessage json.RawMessage
}

func (e *stringError) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
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

func (e *sliceStringError) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
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

func (e *mapStringSliceInterfaceError) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
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

func (e *sliceMapStringInterfaceError) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
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

func (e *boolValue) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
}

func (e *boolValue) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Value)
}

func (e *boolValue) transferTo(ps *ParsedErrors, parent string) {}

// Num struct and unmarshal
type numValue struct {
	Value      int
	RawMessage json.RawMessage
}

func (e *numValue) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
}

func (e *numValue) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Value)
}

func (e *numValue) transferTo(ps *ParsedErrors, parent string) {}

//TODO: implement batch check

var typeRegistry = []string{
	"stringError",
	"sliceStringError",
	"mapStringSliceInterfaceError",
	"sliceMapStringInterfaceError",
	"boolValue",
	"numValue",
}

func makeInstanceStr(name string) ParsedErrorInterface {

	switch name {
	case "stringError":
		return &stringError{}
	case "sliceStringError":
		return &sliceStringError{}
	case "mapStringSliceInterfaceError":
		return &mapStringSliceInterfaceError{}
	case "sliceMapStringInterfaceError":
		return &sliceMapStringInterfaceError{}
	case "boolValue":
		return &boolValue{}
	case "numValue":
		return &numValue{}
	default:
		return nil
	}
}

func batchCheck(s json.RawMessage, ps *ParsedErrors, parent string) error {

	var mainError []error

	for _, name := range typeRegistry {

		parErr := makeInstanceStr(name)
		parErr.setRawMessage(s)

		err := parErr.unmarshalJson()
		if err == nil {
			parErr.transferTo(ps, parent)
		} else {
			mainError = append(mainError, err)
		}
	}

	if len(mainError) > 0 {
		for _, err := range mainError {
			debugMessage(err.Error())
		}
	}

	return nil
}
