package parser

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {

	input, err := ioutil.ReadFile("../_files/trap.txt")
	if err != nil {
		panic(err)
	}

	trap := Parse(input)

	if trap.Variables[".1.3.6.1.2.1.2.2.1.1.5"] != "5" {
		v, _ := trap.Variables[".1.3.6.1.2.1.2.2.1.1.5"]
		t.Errorf("trap .1.3.6.1.2.1.2.2.1.1.5 must be equal 5 ... %v", v)
	}

	if trap.Variables[".1.3.6.1.2.1.2.2.1.8.5"] != "あいうえお" {
		v, _ := trap.Variables[".1.3.6.1.2.1.2.2.1.8.5"]
		t.Errorf("trap .1.3.6.1.2.1.2.2.1.8.5 must be equal あいうえお ... %v", v)
	}
}

func TestParse2(t *testing.T) {

	input := []string{
		`<UNKNOWN>`,
		`UDP: [192.0.2.100]:45111->[127.0.0.1]:162`,
		`.1.3.6.1.2.1.1.3.0 7:4:08:13.14`,
		`.1.3.6.1.2.1.2.2.1.8.5 "E3 81 82 E3 81 84 E3 `,
		`81 86 E3 81 88 E3 81 8A "`,
	}

	trap := Parse([]byte(strings.Join(input, "\n")))

	if trap.Variables[".1.3.6.1.2.1.2.2.1.8.5"] != "あいうえお" {
		v, _ := trap.Variables[".1.3.6.1.2.1.2.2.1.8.5"]
		t.Errorf("trap .1.3.6.1.2.1.2.2.1.8.5 must be equal あいうえお ... %v", v)
	}
}
