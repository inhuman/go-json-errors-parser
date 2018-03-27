package go_json_errors_parser

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
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

// Map of string interface struct and unmarshal
type mapStringInterfaceError struct {
	Error      map[string]interface{}
	RawMessage json.RawMessage
}

func (e *mapStringInterfaceError) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
}

func (e *mapStringInterfaceError) unmarshalJson() error {
	//TODO how check that in interface bool or float or num or string

	err := json.Unmarshal(e.RawMessage, &e.Error)

	if err == nil {

		for _, value := range e.Error {
			s := fmt.Sprintf("%s", value)

			if len(s) >= 2 {
				b := s[0]
				e := s[len(s)-1]

				if (string(b) == "{") && (string(e) == "}") {
					debugMessage("OBJECT FOUND IN INTERFACE")
					return errors.New("Json object found in interface value")
				}

				if (string(b) == "[") && (string(e) == "]") {
					debugMessage("ARRAY FOUND IN INTERFACE")
					return errors.New("Array object found in interface value")
				}
			}
		}
	}

	return err
}

func (e *mapStringInterfaceError) transferTo(ps *ParsedErrors, parent string) {
	r := ParsedError{}
	r.Parent = parent
	tmp := make(map[string][]string)

	for name, err := range e.Error {
		tmp[name] = []string{fmt.Sprintf("%v", err)}
	}

	r.Children = tmp
	ps.ParsedErrors = append(ps.ParsedErrors, r)
}

// Map of string string struct and unmarshal
type mapStringStringError struct {
	Error      map[string]string
	RawMessage json.RawMessage
}

func (e *mapStringStringError) setRawMessage(m json.RawMessage) {
	e.RawMessage = m
}

func (e *mapStringStringError) unmarshalJson() error {
	return json.Unmarshal(e.RawMessage, &e.Error)
}

func (e *mapStringStringError) transferTo(ps *ParsedErrors, parent string) {
	r := ParsedError{}
	r.Parent = parent
	tmp := make(map[string][]string)

	for name, err := range e.Error {
		tmp[name] = []string{fmt.Sprintf("%v", err)}
	}

	r.Children = tmp
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

var typeRegistry = []string{
	"boolValue",
	"numValue",
	"stringError",
	"sliceStringError",
	//"mapStringStringError",

	"sliceMapStringInterfaceError",
	"mapStringSliceInterfaceError",
	"mapStringInterfaceError",
}

func makeInstanceStr(name string) ParsedErrorInterface {

	switch name {
	case "boolValue":
		return &boolValue{}
	case "numValue":
		return &numValue{}
	case "stringError":
		return &stringError{}
	case "sliceStringError":
		return &sliceStringError{}
		//case "mapStringStringError":
		//	return &mapStringStringError{}
	case "mapStringInterfaceError":
		return &mapStringInterfaceError{}
	case "mapStringSliceInterfaceError":
		return &mapStringSliceInterfaceError{}
	case "sliceMapStringInterfaceError":
		return &sliceMapStringInterfaceError{}
	default:
		return nil
	}
}

func batchExtract(s json.RawMessage, ps *ParsedErrors, parent string) error {

	var mainError []error

	debugMessage("Starting batch extractor..")

	for _, name := range typeRegistry {

		debugMessagef("Trying extract to %s\n", name)

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
		return errors.New("There are errors while unmarshaling")
	}

	return nil
}

func batchCheck(s json.RawMessage, continueStructs []string) (error, bool) {

	var mainError []error

	for _, name := range typeRegistry {

		parErr := makeInstanceStr(name)
		parErr.setRawMessage(s)
		err := parErr.unmarshalJson()
		if err == nil {
			if stringInSlice(name, continueStructs) {
				return nil, true
			}
			return nil, false
		} else {
			mainError = append(mainError, err)
		}
	}

	if len(mainError) > 0 {
		for _, err := range mainError {
			debugMessage(err.Error())
		}
		return errors.New("There are errors while unmarshaling"), false
	}

	return nil, true
}

func batchCheckCallback(s json.RawMessage, funcMap map[string]func()) error {

	var mainError []error

	for _, name := range typeRegistry {

		parErr := makeInstanceStr(name)
		parErr.setRawMessage(s)
		err := parErr.unmarshalJson()

		fmt.Println("Check unmarshal to", name)
		fmt.Println(err)
		fmt.Printf("Errored value: %s\n", s)

		if err == nil {
			fmt.Println("Check if calback exists")
			if f, ok := funcMap[name]; ok {
				fmt.Println("Fire callback")
				f()
			}

			return nil
		}
	}

	if len(mainError) > 0 {
		for _, err := range mainError {
			debugMessage(err.Error())
		}
		return errors.New("There are errors while unmarshaling")
	}

	return nil
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
