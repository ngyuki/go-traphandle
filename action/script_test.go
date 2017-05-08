package action

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestExecScript(t *testing.T) {
	defer log.SetOutput(os.Stderr)
	log.SetOutput(ioutil.Discard)

	scriptAct := &scriptAction{"ls / | sort"}
	err := scriptAct.Act(map[string]string{})
	if err != nil {
		t.Error(err)
	}
}
