# go-json-errors-parser 

[![Build Status](https://travis-ci.org/inhuman/go-json-errors-parser.svg?branch=master)](https://travis-ci.org/inhuman/go-json-errors-parser)


Tiny go lib for parsing error fields in json

The lib may be useful, when you don't know where in received json will be an error field and in which format it will be, 
for example:
```json
{
  "Errors": [
    "Unauthorized",
    "Auth required"
  ]
}
``` 
OR
```json
{
  "message": "Validations failed for pipeline 'FromTemplate2'. Error(s): [Validation failed.]. Please correct and resubmit.",
  "data": {
    "Errors": {
      "materials": [
        "A pipeline must have at least one material"
      ],
      "label_template": [
        "Invalid label '123'. Label should be composed of alphanumeric text, it can contain the build number as ${COUNT}, can contain a material revision as ${<material-name>} of ${<material-name>[:<number>]}, or use params as #{<param-name>}."
      ]
    }
  }
}
```
OR
```
{
  "error": "Unauthorized",
  "code": 401
}
```

The lib parse json and find errors, and put them to usable struct like this:
```json
{
  "ParsedErrors": [
    {
      "Children": null,
      "Message": [
        "Validations failed for pipeline 'FromTemplate3'. Error(s): [Validation failed.]. Please correct and resubmit."
      ],
      "Parent": ""
    },
    {
      "Children": {
        "destination": [
          "Invalid Destination Directory. Every material needs a different destination directory and the directories should not be nested.",
          "The destination directory must be unique across materials."
        ]
      },
      "Message": null,
      "Parent": "materials"
    }
  ]
}
```


### Usage

```go
import jerrparser "github.com/inhuman/go-json-errors-parser"


...

json := `{
           "message": "Validations failed for pipeline 'FromTemplate2'. Error(s): [Validation failed.]. Please correct and resubmit.",
           "data": {
             "errors": {
               "materials": [
                 "A pipeline must have at least one material"
               ],
               "label_template": [
                 "Invalid label '123'. Label should be composed of alphanumeric text, it can contain the build number as ${COUNT}, can contain a material revision as ${<material-name>} of ${<material-name>[:<number>]}, or use params as #{<param-name>}."
               ]
             },
             "_links": {
               "doc": {
                 "href": "https://api.gocd.org/#pipeline-config"
               },
               "find": {
                 "href": "http://localhost:8153/go/api/admin/pipelines/:pipeline_name"
               }
             },
             "label_template": "123",
             "origin": null
           }
         }`

errs := jerrparser.ParseErrors(json)

if errs.IsErrors() {
	
    fmt.Println("errors count", errs.GetCount())

    for _, err := range errs.ParsedErrors {
        // handle error messages   
    }
}
```