package config

import (
	"os"
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	config, err := Load("../_files/config.yml")
	if err != nil {
		t.Error(err)
	}

	// Defaults
	if val, ok := config.Defaults["url"]; ok != true || val != "http://example.com/" {
		t.Errorf("config.Defaults ... %v", config.Defaults)
	}

	// Matches
	if len(config.Matches) != 1 {
		t.Errorf("config.Matches ... %v", config.Matches)
	}

	m := config.Matches[0]

	if m.Trap != "IF-MIB::linkDown" {
		t.Errorf("config.Matches ... %v", m)
	}
	if val, ok := m.Bindings["status"]; ok != true || val != "RFC1213-MIB::ifOperStatus.*" {
		t.Errorf("config.Matches ... %v", m)
	}

	// Conditions
	if _, ok := m.Conditions["status"]; ok != true {
		t.Errorf("config.Matches.Conditions ... %v", m.Conditions)
	}
	if _, ok := m.Conditions["status"]["regexp"]; ok != true {
		t.Errorf("config.Matches.Conditions ... %v", m.Conditions)
	}
	if len(m.Conditions["status"]["regexp"]) != 1 {
		t.Errorf("config.Matches.Conditions ... %v", m.Conditions)
	}
	if m.Conditions["status"]["regexp"][0] != "あいうえお" {
		t.Errorf("config.Matches.Conditions ... %v", m.Conditions)
	}

	// Formats
	if _, ok := m.Formats["subject"]; ok != true {
		t.Errorf("config.Matches.Formats ... %v", m.Formats)
	}

	// Actions
	if len(m.Actions.Emails) == 0 {
		t.Errorf("config.Actions ... %v", m.Actions)
	}
	if len(m.Actions.Scripts) == 0 {
		t.Errorf("config.Actions ... %v", m.Actions)
	}
	if val := m.Actions.Emails[0].From; val != "root@example.com" {
		t.Errorf("config.Actions ... %v", m.Actions)
	}
	if val := m.Actions.Scripts[0]; val != "env | sort | logger -t snmptrap" {
		t.Errorf("config.Actions ... %v", m.Actions)
	}
}

func TestLoadConfig_notfound(t *testing.T) {
	_, err := Load("../_files/xxxxx.yml")
	if err == nil {
		t.Error("must be instanceof *os.PathError ... nil")
	}
	if _, ok := err.(*os.PathError); ok == false {
		t.Errorf("must be instanceof *os.PathError ... %v", reflect.TypeOf(err))
	}
}

func TestLoadConfig_invalid(t *testing.T) {
	_, err := Load("../_files/trap.txt")
	if err == nil {
		t.Error("must be err")
	}
}
