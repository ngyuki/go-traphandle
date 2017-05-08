package format

import (
	"testing"
)

func TestApplyTemplate(t *testing.T) {

	text := "My name of {{.name}}. I like {{.foods}}."

	tmpl, err := NewTemplate("hoge", text)
	if err != nil {
		t.Error(err)
	}

	out, err := tmpl.Apply(map[string]string{"name": "ore", "foods": "susi"})
	if err != nil {
		t.Error(err)
	}
	if out != "My name of ore. I like susi." {
		t.Errorf("unexcepted output ... %v", out)
	}
}
