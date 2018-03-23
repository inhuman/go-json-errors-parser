package go_json_errors_parser

import (
	"github.com/hokaccha/go-prettyjson"
	"fmt"
	"os"
)

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

func prettyPrintStruct(strct interface{}) {
	s, _ := prettyjson.Marshal(strct)
	fmt.Println(string(s))
}

const DebugEnvVarName = "GOJSONPARSER_DEBUG"

func debugMessage(string string) {
	if os.Getenv(DebugEnvVarName) == "1" {
		fmt.Println(string)
	}
}

func debugMessagef(i interface{}, format string) {
	if os.Getenv(DebugEnvVarName) == "1" {
		fmt.Printf(format, i)
	}
}

func debugStruct(strct interface{}) {
	if os.Getenv(DebugEnvVarName) == "1" {
		prettyPrintStruct(strct)
	}
}