package utils

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type ConfigData struct {
	Credentials struct {
		User string `yaml:"user"`
		Pwd  string `yaml:"pass"`
		Host string `yaml:"host"`
		Port string `yaml:"port"`
		Env  string `yaml:"env"`
	}
	Monitor struct {
		Queue []string `yaml:"queue"`
	}
}

func ReadConf(filename string) (*ConfigData, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &ConfigData{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}
