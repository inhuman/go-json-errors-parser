package go_json_errors_parser

import (
	"testing"
	"io/ioutil"
	"github.com/stretchr/testify/assert"
)

func TestParseErrorsExample1(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example1.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, 2, errs.GetCount())
	assert.Equal(t, true, errs.IsErrors())
	assert.Equal(t, "data", errs.ParsedErrors[0].Parent)
	assert.Equal(t, "Validations failed for package 'c4b10faf-62f9-4b75-ae7f-9dc042e3d310'. Error(s): [Validation failed.]. Please correct and resubmit.", errs.ParsedErrors[1].Messages[0])
	assert.Equal(t, "Package spec not specified", errs.ParsedErrors[0].Children["PACKAGE_SPEC"][0])
}

func TestParseErrorsExample2(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example2.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, "", errs.ParsedErrors[0].Parent)
	assert.Equal(t, "Unauthorized", errs.ParsedErrors[0].Messages[0])
}

func TestParseErrorsExample3(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example3.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, "data", errs.ParsedErrors[0].Parent)
	assert.Equal(t, "A pipeline must have at least one material", errs.ParsedErrors[0].Children["materials"][0])
	assert.Equal(t, "Invalid label '123'. Label should be composed of alphanumeric text, it can contain the build number as ${COUNT}, can contain a material revision as ${<material-name>} of ${<material-name>[:<number>]}, or use params as #{<param-name>}.", errs.ParsedErrors[0].Children["label_template"][0])
}

func TestParseErrorsExample4(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example4.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, "data", errs.ParsedErrors[0].Parent)
	assert.Equal(t, "A pipeline must have at least one material", errs.ParsedErrors[0].Children["materials"][0])
}

func TestParseErrorsExample5(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example5.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, "materials", errs.ParsedErrors[0].Parent)
	assert.Equal(t, 0, len(errs.ParsedErrors[0].Messages))
	assert.Equal(t, "Invalid Destination Directory. Every material needs a different destination directory and the directories should not be nested.", errs.ParsedErrors[1].Children["destination"][0])
}

func TestParseErrorsExample6(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example6.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, "data", errs.ParsedErrors[0].Parent)
	assert.Equal(t, "Unauthorized", errs.ParsedErrors[0].Messages[0])
}

func TestParseErrorsExample7(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example7.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, "data", errs.ParsedErrors[0].Parent)
	assert.Equal(t, "some error", errs.ParsedErrors[0].Children["FieldName"][0])
	assert.Equal(t, "some error", errs.ParsedErrors[0].Children["FieldName2"][0])
}

func TestParseErrorsExample8(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example8.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))

	assert.Equal(t, false, errs.IsErrors())
}

func TestParseErrorsExample9(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example9.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))
	assert.Equal(t, true, errs.IsErrors())
	assert.Equal(t, "Unauthorized", errs.ParsedErrors[0].Messages[0])
	assert.Equal(t, "Auth required", errs.ParsedErrors[0].Messages[1])

}

func TestParseErrorsExample10(t *testing.T) {

	file, e := ioutil.ReadFile("tests/example10.json")
	assert.NoError(t, e)

	errs := ParseErrors(string(file))
	assert.Equal(t, true, errs.IsErrors())

	prettyPrintStruct(errs)

}

func TestParsedErrors_GetErrors(t *testing.T) {
	file, e := ioutil.ReadFile("tests/example3.json")
	assert.NoError(t, e)
	errs := ParseErrors(string(file))

	errors := errs.GetErrors()
	assert.Equal(t, "[data][materials] A pipeline must have at least one material", errors[2].Error())
}

//func TestVoid1(t *testing.T) {
//	runWithFile(1, t)
//}
//
//func TestVoid2(t *testing.T) {
//	runWithFile(2, t)
//}
//
//func TestVoid3(t *testing.T) {
//	runWithFile(3, t)
//}
//
//func TestVoid5(t *testing.T) {
//	runWithFile(5, t)
//}
//
//func TestVoid6(t *testing.T) {
//	runWithFile(6, t)
//}
//
//func TestVoid9(t *testing.T) {
//	runWithFile(9, t)
//}
//
//func TestVoid7(t *testing.T) {
//	runWithFile(7, t)
//}
//
//func TestVoid10(t *testing.T) {
//	runWithFile(10, t)
//}
//
//func runWithFile(i int, t *testing.T) {
//
//	file, e := ioutil.ReadFile("tests/example" + strconv.Itoa(i) + ".json")
//	assert.NoError(t, e)
//
//	ParseErrors(string(file))
//}
