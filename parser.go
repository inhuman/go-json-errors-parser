package go_json_errors_parser

type ParsedError struct {
	Parent string
	Children map[string][]string
	Message []string
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


func ParseErrors(json string) *ParsedErrors {

	errs := ParsedErrors{}


	return &errs
}



