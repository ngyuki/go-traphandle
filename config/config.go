package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Load(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

type Config struct {
	Defaults map[string]string
	Matches  []MatchConfig
}

type MatchConfig struct {
	Trap        string
	Ipaddr      string
	Community   string
	Bindings    map[string]string
	Conditions  map[string]ConditionConfig
	Formats     map[string]string
	Actions     ActionConfig
	Fallthrough bool
}

type ConditionConfig map[string][]string

type ActionConfig struct {
	Emails  []EmailConfig
	Scripts []string
}

type EmailConfig struct {
	Host string
	Port int
	From string
	To   string
}
