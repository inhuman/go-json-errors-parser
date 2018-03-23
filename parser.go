package go_json_errors_parser

import (
	"encoding/json"
	"fmt"
	"regexp"
	"github.com/hokaccha/go-prettyjson"
)

type ParsedError struct {
	Parent   string
	Children map[string][]string
	Message  []string
}

type ParsedErrors struct {
	ParsedErrors []ParsedError
}

func (pe *ParsedErrors) IsErrors() bool {
	return len(pe.ParsedErrors) > 0
}

func (pe *ParsedErrors) GetCount() int {
	return len(pe.ParsedErrors)
}

func ParseErrors(jsn string) *ParsedErrors {

	errs := ParsedErrors{}

	// Unmarshal given json to temporary map
	var tmpMap map[string]*json.RawMessage

	if err := json.Unmarshal([]byte(jsn), &tmpMap); err != nil {
		panic(err)
	}

	walk(tmpMap, &errs)

	//fmt.Printf("%+v\n", errs)
	fmt.Printf("result struct:\n")
	PrettyPrintStruct(errs)

	return &errs
}

func walk(item map[string]*json.RawMessage, ps *ParsedErrors) {

	fmt.Printf("result struct:\n")
	PrettyPrintStruct(ps)

	re := regexp.MustCompile(`(?i)error`)

	for key, s := range item {

		fmt.Printf("%s\n", key)
		fmt.Printf("%s\n", s)

		if re.MatchString(string(key)) {

			//TODO: try to unmarshal to string slice
			//TODO: try to unmarshal to map[string]*json.RawMessage

			fmt.Println("ERROR PARSED IN KEY: " + string(key))

			// try to unmarshal to string
			str, err := tryUnmarshalToString(s)
			if err == nil {
				addStringError(*str, ps)
				continue
			}

			strs, err := tryUnmarshalToStringSlice(s)
			if err == nil {
				addStringSliceError(*strs, ps)
				continue
			}

			maps, err := tryUnmarshalToStringSliceMap(s)
			if err == nil {
				addStringSliceMapError(*maps, ps)
				continue
			}

			objMaps, err := tryUnmarshalToObjectsSliceMap(s)
			fmt.Println(err)

			if err == nil {
				fmt.Println("UNMARSHAL to obj slice map")

				addObjectSliceMapError(*objMaps, ps)
				continue
			}


		} else {

			if s == nil {
				continue
			}

			_, err := tryUnmarshalToString(s)
			if err == nil {
				continue
			}

			_, err = tryUnmarshalToNum(s)
			if err == nil {
				continue
			}

			_, err = tryUnmarshalToStringSlice(s)
			if err == nil {
				continue
			}

			_, err = tryUnmarshalToBool(s)
			if err == nil {
				continue
			}

			_, err = tryUnmarshalToObjectsSliceMap(s)
			fmt.Println("tryUnmarshalToObjectsSliceMap", err)
			if err == nil {
				continue
			}

			_, err = tryUnmarshalToStringSliceMap(s)
			if err == nil {
				continue
			}


			//TODO: finish the parsor

			fmt.Println("ELSE fires:")
			PrettyPrintStruct(s)

			var tmpMap map[string]*json.RawMessage
			err = json.Unmarshal(*s, &tmpMap)
			checkErr(err)

			walk(tmpMap, ps)
		}

		//b, err := s.MarshalJSON()
		//checkErr(err)
		//
		//var tmpMap map[string]*json.RawMessage
		//
		//if err := json.Unmarshal(b, &tmpMap); err != nil {
		//
		//	if re.MatchString(string(b)) {
		//		fmt.Println("ERROR PARSED IN VALUE: " + string(b))
		//
		//		addStringError(string(b), ps)
		//	}
		//
		//	fmt.Println(string(b))
		//} else {
		//	var tmpMap map[string]*json.RawMessage
		//	err := json.Unmarshal(b, &tmpMap)
		//	if err != nil {
		//		panic(err)
		//	}
		//
		//	walk(tmpMap, ps)
		//}
	}
}

func tryUnmarshalToString(i *json.RawMessage) (*string, error) {
	var tmpMap string

	err := json.Unmarshal(*i, &tmpMap)

	if err != nil {
		return nil, err
	}

	return &tmpMap, nil
}

func tryUnmarshalToNum(i *json.RawMessage) (*int, error) {
	var tmpMap int

	err := json.Unmarshal(*i, &tmpMap)

	if err != nil {
		return nil, err
	}

	return &tmpMap, nil
}



func tryUnmarshalToStringSlice(i *json.RawMessage) (*[]string, error) {
	var tmpMap []string

	err := json.Unmarshal(*i, &tmpMap)

	if err != nil {
		return nil, err
	}

	return &tmpMap, nil
}

func tryUnmarshalToStringSliceMap(i *json.RawMessage) (*map[string][]interface{}, error) {
	var tmpMap map[string][]interface{}

	err := json.Unmarshal(*i, &tmpMap)

	if err != nil {
		return nil, err
	}

	return &tmpMap, nil
}

func tryUnmarshalToBool(i *json.RawMessage) (*bool, error) {
	var tmpMap bool

	err := json.Unmarshal(*i, &tmpMap)

	if err != nil {
		return nil, err
	}

	return &tmpMap, nil
}

func tryUnmarshalToObjectsSliceMap(i *json.RawMessage) (*[]map[string]string, error) {

	var tmpMap []map[string]string

	err := json.Unmarshal(*i, &tmpMap)

	if err != nil {
		return nil, err
	}

	return &tmpMap, nil
}

func addStringError(str string, ps *ParsedErrors) {

	e := ParsedError{}
	e.Message = append(e.Message, trimQuotes(str))
	ps.ParsedErrors = append(ps.ParsedErrors, e)
}

func addStringSliceError(strs []string, ps *ParsedErrors) {

	e := ParsedError{}
	e.Message = append(e.Message, strs...)
	ps.ParsedErrors = append(ps.ParsedErrors, e)
}

func addStringSliceMapError(strs map[string][]interface{}, ps *ParsedErrors) {
	e := ParsedError{}

	formattedStrs := make(map[string][]string)

	for name, str := range strs {

		var tmpMap []string

		for _, item := range str {
			tmpMap = append(tmpMap, fmt.Sprintf("%v", item))
		}

		formattedStrs[name] = tmpMap
	}

	e.Children = formattedStrs
	ps.ParsedErrors = append(ps.ParsedErrors, e)
}

func addObjectSliceMapError(strs []map[string]string, ps *ParsedErrors) {
	e := ParsedError{}

	for _, value := range strs {

		tmp := make(map[string][]string)

		for key, val := range value {
			tmp[key] = []string{val}
		}

		e.Children = tmp
	}

	ps.ParsedErrors = append(ps.ParsedErrors, e)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func trimQuotes(s string) string {
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s
}

func PrettyPrintStruct(strct interface{}) {

	s, _ := prettyjson.Marshal(strct)
	fmt.Println(string(s))
}
