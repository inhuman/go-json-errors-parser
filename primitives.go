package go_json_errors_parser

import (
	"encoding/json"
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

//Slice of map string interface error struct and unmarshal
type sliceMapStringInterfaceError struct {
	Error      []map[string]interface{}
	RawMessage json.RawMessage
}

func (e *sliceMapStringInterfaceError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

// Bool struct and unmarshal
type boolValue struct {
	Value      bool
	RawMessage json.RawMessage
}

func (e *boolValue) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Value)
}

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
