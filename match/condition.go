package match

import (
	"fmt"
	"regexp"
)

type Condition interface {
	IsMatch(string) bool
}

type callbackCondition struct {
	callback func(v string) bool
	not      bool
}

func (c *callbackCondition) IsMatch(v string) bool {
	return c.callback(v)
}

func newCondition(op string, vals []string) Condition {
	switch op {
	case "eq", "not_eq":
		vmap := make(map[string]bool)
		for _, v := range vals {
			vmap[v] = true
		}
		callback := func(v string) bool {
			return vmap[v]
		}
		return &callbackCondition{callback, op == "not_eq"}

	case "regexp", "not_regexp":
		regs := make([]*regexp.Regexp, 0)
		for _, v := range vals {
			regs = append(regs, regexp.MustCompile(v))
		}
		callback := func(v string) bool {
			for _, r := range regs {
				if r.MatchString(v) == false {
					return false
				}
			}
			return true
		}
		return &callbackCondition{callback, op == "not_regexp"}
	}
	panic(fmt.Errorf("Unknown condition %v", op))
}
