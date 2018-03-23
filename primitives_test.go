package go_json_errors_parser

import (
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalString(t *testing.T) {

	var unmarshaledError stringError
	unmarshaledError.RawMessage = json.RawMessage(`"Unauthorized"`)

	err := unmarshaledError.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, "Unauthorized", unmarshaledError.Error)

	parsedErrors := ParsedErrors{}
	unmarshaledError.transferTo(&parsedErrors, "")
	assert.Equal(t, "Unauthorized", parsedErrors.ParsedErrors[0].Messages[0])

	var unmarshaledError2 stringError
	unmarshaledError2.RawMessage = json.RawMessage(`{"error": "Unauthorized"}`)

	err2 := unmarshaledError2.unmarshalJson()
	assert.Error(t, err2, "json: cannot unmarshal object into Go value of type string")
}

func TestUnmarshalSliceString(t *testing.T) {

	var unmarshaledError sliceStringError
	unmarshaledError.RawMessage = json.RawMessage(`["Unauthorized", "Auth required"]`)

	err := unmarshaledError.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, "Unauthorized", unmarshaledError.Error[0])
	assert.Equal(t, "Auth required", unmarshaledError.Error[1])

	parsedErrors := ParsedErrors{}
	unmarshaledError.transferTo(&parsedErrors, "")

	var unmarshaledError2 sliceStringError
	unmarshaledError2.RawMessage = json.RawMessage(`"Errors": ["Unauthorized","Auth required"]`)

	err2 := unmarshaledError2.unmarshalJson()
	assert.Error(t, err2)
}

func TestMapStringSliceInterfaceError(t *testing.T) {

	var unmarshaledError mapStringSliceInterfaceError
	unmarshaledError.RawMessage = json.RawMessage(`{"Errors": ["Unauthorized", "Auth required"]}`)

	err := unmarshaledError.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, "Unauthorized", unmarshaledError.Error["Errors"][0])
}

func TestSliceMapStringInterfaceError(t *testing.T) {

	var unmarshaledError sliceMapStringInterfaceError
	unmarshaledError.RawMessage = json.RawMessage(`[{"secure": false, "name": "ADF", "value": "123"}]`)

	err := unmarshaledError.unmarshalJson()
	assert.NoError(t, err)
	assert.Equal(t, false, unmarshaledError.Error[0]["secure"])
	assert.Equal(t, "ADF", unmarshaledError.Error[0]["name"])

	var unmarshaledError2 sliceMapStringInterfaceError
	unmarshaledError2.RawMessage = json.RawMessage(`{"huembuem" : [{"secure": false, "name": "ADF", "value": "123"}]}`)

	err2 := unmarshaledError2.unmarshalJson()
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
