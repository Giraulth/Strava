package tools

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func ReadConfig(filename string, config interface{}) error {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("Error reading YAML file: %v", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return fmt.Errorf("Error parsing YAML file: %v", err)
	}

	return nil
}
