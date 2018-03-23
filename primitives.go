package go_json_errors_parser

import "encoding/json"


// String error struct and unmarshal
type stringError struct {
	Error      string
	RawMessage json.RawMessage
}

func (e *stringError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

// String slice error struct and unmarshal
type sliceStringError struct {
	Error []string
	RawMessage json.RawMessage
}
func (e *sliceStringError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

// Map of string slice interface error struct and unmarshal
type mapStringSliceInterfaceError struct {
	Error map[string][]interface{}
	RawMessage json.RawMessage
}
func (e *mapStringSliceInterfaceError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

//Slice of map string interface error struct and unmarshal
type sliceMapStringInterfaceError struct {
	Error []map[string]interface{}
	RawMessage json.RawMessage
}
func (e *sliceMapStringInterfaceError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

// Bool struct and unmarshal
type boolValue struct {
	Value bool
	RawMessage json.RawMessage
}

func (e *boolValue) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Value)
}