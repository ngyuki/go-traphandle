package action

import (
	"testing"

	"github.com/ngyuki/go-traphandle/config"
)

func TestNewActions(t *testing.T) {
	cfg := &config.ActionConfig{
		Emails: []config.EmailConfig{
			{
				From: "from@example.com",
				To:   "to@example.com",
			},
			{
				From: "from@example.net",
				To:   "to@example.net",
				Host: "mail.example.net",
				Port: 2525,
			},
		},
		Scripts: []string{
			"/path/to/scripts/foo",
			"/path/to/scripts/bar",
		},
	}

	actions, err := NewActions(cfg)
	if err != nil {
		t.Errorf("NewActions returned error ... %v", err)
	}

	{
		act := actions[0].(*emailAction)
		actions = actions[1:]
		if exp := "from@example.com"; act.From != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.From)
		}
		if exp := "localhost"; act.Host != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.Host)
		}
		if exp := 25; act.Port != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.Port)
		}
	}
	{
		act := actions[0].(*emailAction)
		actions = actions[1:]
		if exp := "to@example.net"; act.To != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.From)
		}
		if exp := "mail.example.net"; act.Host != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.Host)
		}
		if exp := 2525; act.Port != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.Port)
		}
	}
	{
		act := actions[0].(*scriptAction)
		actions = actions[1:]
		if exp := "/path/to/scripts/foo"; act.script != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.script)
		}
	}
	{
		act := actions[0].(*scriptAction)
		actions = actions[1:]
		if exp := "/path/to/scripts/bar"; act.script != exp {
			t.Errorf("shouled be %v ... given %v", exp, act.script)
		}
	}
}

func TestEmailMissing(t *testing.T) {
	{
		cfg := &config.ActionConfig{
			Emails: []config.EmailConfig{
				{
					//From: "from@example.net",
					To:   "to@example.net",
					Host: "mail.example.net",
					Port: 2525,
				},
			},
		}

		_, err := NewActions(cfg)
		if err == nil {
			t.Errorf("shouled be return err ... returned nil")
		}
	}
	{
		cfg := &config.ActionConfig{
			Emails: []config.EmailConfig{
				{
					From: "from@example.net",
					//To:   "to@example.net",
					Host: "mail.example.net",
					Port: 2525,
				},
			},
		}

		_, err := NewActions(cfg)
		if err == nil {
			t.Errorf("shouled be return err ... returned nil")
		}
	}
}
