package action

import (
	"github.com/ngyuki/go-traphandle/config"
)

type Acter interface {
	Act(values map[string]string) error
}

func NewActions(cfg *config.ActionConfig) ([]Acter, error) {

	actions := make([]Acter, 0)

	for _, email := range cfg.Emails {
		act, err := newEmailAction(&email)
		if err != nil {
			return nil, err
		}
		actions = append(actions, act)
	}

	for _, script := range cfg.Scripts {
		act := &scriptAction{script}
		actions = append(actions, act)
	}

	return actions, nil
}
