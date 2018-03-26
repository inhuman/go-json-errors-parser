package go_json_errors_parser

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalString(t *testing.T) {

	// success unmarshal
	var unmarshaledErrorSuccess stringError
	unmarshaledErrorSuccess.RawMessage = json.RawMessage(`"Unauthorized"`)

	err := unmarshaledErrorSuccess.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, "Unauthorized", unmarshaledErrorSuccess.Error)

	// transfer to parsed errors struct
	parsedErrors := ParsedErrors{}
	unmarshaledErrorSuccess.transferTo(&parsedErrors, "")
	assert.Equal(t, "Unauthorized", parsedErrors.ParsedErrors[0].Messages[0])

	// fault unmarshal
	var unmarshaledErrorFault stringError
	unmarshaledErrorFault.RawMessage = json.RawMessage(`{"error": "Unauthorized"}`)

	err2 := unmarshaledErrorFault.unmarshalJson()
	assert.Error(t, err2, "json: cannot unmarshal object into Go value of type string")
}

func TestUnmarshalSliceString(t *testing.T) {

	// success unmarshal
	var unmarshaledErrorSuccess sliceStringError
	unmarshaledErrorSuccess.RawMessage = json.RawMessage(`["Unauthorized", "Auth required"]`)
	err := unmarshaledErrorSuccess.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, "Unauthorized", unmarshaledErrorSuccess.Error[0])
	assert.Equal(t, "Auth required", unmarshaledErrorSuccess.Error[1])

	// transfer to parsed errors struct
	parsedErrors := ParsedErrors{}
	unmarshaledErrorSuccess.transferTo(&parsedErrors, "parent")
	assert.Equal(t, "Unauthorized", parsedErrors.ParsedErrors[0].Messages[0])
	assert.Equal(t, "Auth required", parsedErrors.ParsedErrors[0].Messages[1])
	assert.Equal(t, "parent", parsedErrors.ParsedErrors[0].Parent)


	// fault unmarshal
	var unmarshaledErrorFault sliceStringError
	unmarshaledErrorFault.RawMessage = json.RawMessage(`"Errors": ["Unauthorized","Auth required"]`)

	err2 := unmarshaledErrorFault.unmarshalJson()
	assert.Error(t, err2)
}

func TestMapStringSliceInterfaceError(t *testing.T) {

	// success unmarshal
	var unmarshaledError mapStringSliceInterfaceError
	unmarshaledError.RawMessage = json.RawMessage(`{"Errors": ["Unauthorized", "Auth required"]}`)

	err := unmarshaledError.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, "Unauthorized", unmarshaledError.Error["Errors"][0])

	// transfer to parsed errors struct
	parsedErrors := ParsedErrors{}
	unmarshaledError.transferTo(&parsedErrors, "TestParent")

	assert.Equal(t, "[TestParent][Errors] Auth required", parsedErrors.GetErrors()[0].Error())
	assert.Equal(t, "[TestParent][Errors] Unauthorized", parsedErrors.GetErrors()[1].Error())

	// fault unmarshal
	var unmarshaledErrorFault mapStringSliceInterfaceError
	unmarshaledErrorFault.RawMessage = json.RawMessage(`["Unauthorized","Auth required"]`)
	err2 := unmarshaledErrorFault.unmarshalJson()
	assert.Error(t, err2)
}

func TestSliceMapStringInterfaceError(t *testing.T) {

	// success unmarshal
	var unmarshaledError sliceMapStringInterfaceError
	unmarshaledError.RawMessage = json.RawMessage(`[{"secure": false, "name": "ADF", "value": "123"}]`)

	err := unmarshaledError.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, false, unmarshaledError.Error[0]["secure"])
	assert.Equal(t, "ADF", unmarshaledError.Error[0]["name"])

	// transfer to parsed errors struct
	parsedErrors := ParsedErrors{}
	unmarshaledError.transferTo(&parsedErrors, "TestParent")

	assert.Equal(t, "[TestParent][name] ADF", parsedErrors.GetErrors()[0].Error())
	assert.Equal(t, "[TestParent][secure] false", parsedErrors.GetErrors()[1].Error())

	// unmarshal error fault
	var unmarshaledErrorFault sliceMapStringInterfaceError
	unmarshaledErrorFault.RawMessage = json.RawMessage(`{"huembuem" : [{"secure": false, "name": "ADF", "value": "123"}]}`)

	err2 := unmarshaledErrorFault.unmarshalJson()
	assert.Error(t, err2)
}

func TestBoolValue(t *testing.T) {

	var unmarshaledError boolValue
	unmarshaledError.RawMessage = json.RawMessage(`false`)
	err := unmarshaledError.unmarshalJson()
	assert.NoError(t, err)

	var unmarshaledError2 boolValue
	unmarshaledError2.RawMessage = json.RawMessage(`123`)
	err2 := unmarshaledError2.unmarshalJson()
	assert.Error(t, err2)
}

//func TestBatchCheck(t *testing.T) {
//
//	rawMessage := json.RawMessage(`"Unauthorized"`)
//	parsedErrors := ParsedErrors{}
//
//	batchCheck(&rawMessage, &parsedErrors, "")
//
//}
