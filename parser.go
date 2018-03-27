package go_json_errors_parser

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"sort"
)

type ParsedError struct {
	Parent   string
	Children map[string][]string
	Messages []string
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

func (pe *ParsedErrors) GetErrors() []error {

	var errs []error

	for _, parsedError := range pe.ParsedErrors {

		// Collect errors from Messages
		if len(parsedError.Messages) > 0 {
			for _, msg := range parsedError.Messages {
				errs = append(errs, errors.New("[] "+msg))
			}
		}

		// Collect errors from children
		for name, children := range parsedError.Children {
			for _, child := range children {
				errStr := "[" + parsedError.Parent + "][" + name + "] " + child
				errs = append(errs, errors.New(errStr))
			}
		}
	}

	sort.Slice(errs[:], func(i, j int) bool {
		return errs[i].Error() < errs[j].Error()
	})

	return errs
}

// Main method
func ParseErrors(jsn string) *ParsedErrors {

	errs := ParsedErrors{}

	// Unmarshal given json to temporary map
	var tmpMap map[string]*json.RawMessage

	if err := json.Unmarshal([]byte(jsn), &tmpMap); err != nil {
		panic(err)
	}

	walk(tmpMap, &errs, "")
	debugMessage("Final result struct:")
	debugStruct(errs)

	sort.Slice(errs.ParsedErrors[:], func(i, j int) bool {
		return len(errs.ParsedErrors[i].Messages) < len(errs.ParsedErrors[j].Messages)
	})

	return &errs
}

// Recursively walks throw entire json, unmarshal and
// search substring 'error' by regexp (case insensitive match) in keys and values
// and puts found errors into struct
func walk(item map[string]*json.RawMessage, ps *ParsedErrors, parent string) {

	debugMessage("intermediate result")
	debugStruct(ps)

	re := regexp.MustCompile(`(?i)(.+|.?)(error)(.+|.?)`)

	for key, s := range item {

		debugMessagef("Key: %s\n", key)
		debugMessagef("Value: %s\n", s)

		// check if errors in value
		str := fmt.Sprintf("%s", s)
		if re.MatchString(str) {
			debugMessagef("ERROR FOUND IN VALUE: %s\n", str)

			var unmarshaledError stringError
			unmarshaledError.RawMessage = *s
			err := unmarshaledError.unmarshalJson()
			if err == nil {
				unmarshaledError.transferTo(ps, parent)
				continue
			} else {
				debugMessage("Error-in-value not works: can't unmarshal")
				debugMessage(err.Error())
			}
		}

		if re.MatchString(string(key)) {
			debugMessagef("ERROR FOUND IN KEY: %s\n", key)

			err := batchExtract(*s, ps, parent)
			if err == nil {
				continue
			} else {
				debugMessage(err.Error())
			}

		} else {

			debugMessage("ELSE fires:")
			debugStruct(s)

			if s == nil {
				debugMessage("detect s is nil, skipping..")
				continue
			}

			debugMessage("Checking all with continue on mapStringSliceInterfaceError")
			err, cont := batchCheck(*s, []string{"mapStringSliceInterfaceError"})
			if (err == nil) && cont {

				var tmpMap map[string]*json.RawMessage
				err = json.Unmarshal(*s, &tmpMap)
				checkErr(err)

				debugMessage("detect mapStringSliceInterfaceError, going deeper..")
				walk(tmpMap, ps, key)

				continue
			} else {
				if err != nil {
					debugMessage(err.Error())
				}
			}

			debugMessage("Checking all with continue on sliceMapStringInterfaceError")
			err, cont = batchCheck(*s, []string{"sliceMapStringInterfaceError"})
			if err == nil {
				if cont {
					var tmpMap []map[string]*json.RawMessage
					err = json.Unmarshal(*s, &tmpMap)
					checkErr(err)

					debugMessage("detect sliceMapStringInterfaceError, going deeper..")

					for _, value := range tmpMap {

						debugMessage("parsing sub struct")
						debugMessagef("%s", value)

						walk(value, ps, key)
					}
				}

				debugMessage("Err is nil but cont is false")

				continue
			} else {
				debugMessage(err.Error())
			}

			var tmpMap map[string]*json.RawMessage
			err = json.Unmarshal(*s, &tmpMap)
			checkErr(err)

			debugMessage("PARENT set to: " + key)
			walk(tmpMap, ps, key)
		}
	}
}
