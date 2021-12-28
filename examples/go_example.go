package main

import "github.com/ilyapt/go-ansible-module/ansible_module"

// any number and type of input params
type inputParams struct {
	InputParam        int `json:"input_param"`
	AnotherInputParam string `json:"another_input_param"`
}

func main() {
	var input inputParams
	module := ansible_module.New(&input)

	// simple logging
	module.LogPrint("Hi, go_ansible_module")

	// set returned value
	module.Set("example_response_key", "example_value")

	if input.InputParam < 5 {
		module.LogPrint("Input value is good")
		module.Done(true)
	}
	module.FailWithErrorf("Input value too big")
}