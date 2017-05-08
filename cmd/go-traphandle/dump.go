package main

import (
	"gopkg.in/yaml.v2"
)

func dump(obj interface{}) string {
	out, _ := yaml.Marshal(obj)
	return "---\n" + string(out) + "---\n"
}
