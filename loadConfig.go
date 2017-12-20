package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

func LoadConfig(filename string) (apiv1.AgentOptions, error) {
	var ao apiv1.AgentOptions
	configFromFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiv1.AgentOptions{}, err
	}

	err = json.Unmarshal(configFromFile, &ao)
	if err != nil {
		return apiv1.AgentOptions{}, err
	}
	return ao, nil
}
