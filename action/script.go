package action

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

type scriptAction struct {
	script string
}

func (act *scriptAction) Act(values map[string]string) error {

	cmd := exec.Command("/bin/sh", "-c", act.script)
	cmd.Env = os.Environ()

	for k, v := range values {
		cmd.Env = append(cmd.Env, "traphandle_"+k+"="+v)
	}

	out, err := cmd.CombinedOutput()

	if len(out) > 0 {
		log.Println("script output ...")
		for _, s := range strings.Split(string(out), "\n") {
			log.Printf("  %s\n", s)
		}
	}

	return err
}
