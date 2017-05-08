package action

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"strings"
	"time"

	"github.com/ngyuki/go-traphandle/config"
)

type emailAction struct {
	config.EmailConfig
}

func newEmailAction(cfg *config.EmailConfig) (*emailAction, error) {

	act := &emailAction{EmailConfig: *cfg}

	if len(act.From) == 0 {
		return nil, fmt.Errorf("must be set 'from' by %+v", act.EmailConfig)
	}

	if len(act.To) == 0 {
		return nil, fmt.Errorf("must be set 'to' by %v", act.EmailConfig)
	}

	if len(act.Host) == 0 {
		act.Host = "localhost"
	}

	if act.Port == 0 {
		act.Port = 25
	}

	return act, nil
}

func (act *emailAction) Act(values map[string]string) error {

	subject := values["subject"]
	body := values["body"]

	m := gomail.NewMessage()
	m.SetDateHeader("Date", time.Now())
	m.SetHeader("From", act.From)
	m.SetHeader("To", act.To)
	m.SetHeader("Subject", strings.TrimSpace(subject))
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(act.Host, act.Port, "", "")
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
