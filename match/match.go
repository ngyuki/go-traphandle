package match

import (
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/ngyuki/go-traphandle/config"
	"github.com/ngyuki/go-traphandle/types"
)

type Match struct {
	Trap       string
	Ipaddr     *net.IPNet
	Community  string
	Bindings   map[string]*regexp.Regexp
	Conditions map[string][]Condition
}

func NewMatch(cfg *config.MatchConfig) (m *Match, err error) {
	defer func() {
		if e := recover(); e != nil {
			m = nil
			err = e.(error)
		}
	}()
	return newMatch(cfg), nil
}

func newMatch(cfg *config.MatchConfig) *Match {

	m := &Match{}

	if len(cfg.Trap) > 0 {
		raw, err := translateOidToRaw(cfg.Trap)
		if err != nil {
			panic(err)
		}
		m.Trap = raw
	}

	if len(cfg.Ipaddr) > 0 {
		_, ipnet, err := net.ParseCIDR(cfg.Ipaddr)
		if err != nil {
			panic(fmt.Errorf("%v ... %v", err, cfg.Ipaddr))
		}
		m.Ipaddr = ipnet
	}

	if len(cfg.Community) > 0 {
		m.Community = cfg.Community
	}

	m.Bindings = make(map[string]*regexp.Regexp)

	for name, oid := range cfg.Bindings {
		raw, err := translateOidToRaw(oid)
		if err != nil {
			panic(err)
		}
		m.Bindings[name] = compileOidRegexp(raw)
	}

	m.Conditions = make(map[string][]Condition)

	for name, conditions := range cfg.Conditions {
		conds := make([]Condition, 0)
		for key, val := range conditions {
			conds = append(conds, newCondition(key, val))
		}
		m.Conditions[name] = conds
	}

	return m
}

func compileOidRegexp(glob string) *regexp.Regexp {

	suffix := ""

	switch {
	case strings.HasSuffix(glob, ".**"):
		glob = strings.TrimSuffix(glob, ".**")
		suffix = `(\.\d+)+`
	case strings.HasSuffix(glob, ".*"):
		glob = strings.TrimSuffix(glob, ".*")
		suffix = `\.\d+`
	}

	p := `\A` + regexp.QuoteMeta(glob) + suffix + `\z`
	r := regexp.MustCompile(p)
	return r
}

func (m *Match) Match(trap *types.Trap) (map[string]string, bool) {

	variables := make(map[string]string)

	if m.matchBasic(trap) == false {
		return variables, false
	}

	variables = m.lookupVariables(trap.Variables)

	if m.matchConditions(variables) == false {
		return variables, false
	}

	return variables, true
}

func (m *Match) matchBasic(trap *types.Trap) bool {

	if m.Ipaddr != nil && m.Ipaddr.Contains(net.ParseIP(trap.Ipaddr)) == false {
		return false
	}

	if len(m.Community) > 0 && trap.Community != m.Community {
		return false
	}

	if len(m.Trap) > 0 && trap.Trap != m.Trap {
		return false
	}

	return true
}

func (m *Match) lookupVariables(data map[string]string) map[string]string {

	variables := make(map[string]string)

	for name, r := range m.Bindings {

		variables[name] = ""

		for oid, val := range data {
			if r.MatchString(oid) {
				variables[name] = val
				break
			}
		}
	}

	return variables
}

/*
func lookupGlobOid(glob string, dataBinding map[string]string) (string, bool) {

	suffix := ""

	switch {
	case strings.HasSuffix(glob, ".**"):
		glob = strings.TrimSuffix(glob, ".**")
		suffix = `(\.\d+)+`
	case strings.HasSuffix(glob, ".*"):
		glob = strings.TrimSuffix(glob, ".*")
		suffix = `\.\d+`
	}

	p := `\A` + regexp.QuoteMeta(glob) + suffix + `\z`
	r, err := regexp.Compile(p)
	if err != nil {
		return "", false
	}

	for oid, val := range dataBinding {
		if r.MatchString(oid) {
			return val, true
		}
	}
	return "", false
}*/

func (m *Match) matchConditions(variables map[string]string) bool {

	for name, conds := range m.Conditions {
		for _, cond := range conds {
			val, ok := variables[name]
			if ok == false {
				return false
			}
			if cond.IsMatch(val) == false {
				return false
			}
		}
	}

	return true
}
