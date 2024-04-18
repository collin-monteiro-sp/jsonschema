package jsonschema_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type dclarkRegexp regexp2.Regexp

func (re *dclarkRegexp) MatchString(s string) bool {
	matched, err := (*regexp2.Regexp)(re).MatchString(s)
	return err == nil && matched
}

func (re *dclarkRegexp) String() string {
	return (*regexp2.Regexp)(re).String()
}

func dlclarkCompile(s string) (jsonschema.Regexp, error) {
	re, err := regexp2.Compile(s, regexp2.ECMAScript)
	if err != nil {
		return nil, err
	}
	return (*dclarkRegexp)(re), nil
}

func Example_customRegexEngine() {
	// golang regexp does not support escape sequence: `\c`
	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(`{
		"type": "string",
		"pattern": "^\\cc$"
	}`))
	if err != nil {
		log.Fatal(err)
	}
	inst, err := jsonschema.UnmarshalJSON(strings.NewReader(`"\u0003"`))
	if err != nil {
		log.Fatal(err)
	}

	c := jsonschema.NewCompiler()
	c.UseRegexpEngine(dlclarkCompile)
	if err := c.AddResource("schema.json", schema); err != nil {
		log.Fatal(err)
	}

	sch, err := c.Compile("schema.json")
	if err != nil {
		log.Fatal(err)
	}
	err = sch.Validate(inst)
	fmt.Println("valid:", err == nil)
	// Output:
	// valid: true
}
