package parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/ngyuki/go-traphandle/types"
)

const (
	oidUpTime     = ".1.3.6.1.2.1.1.3.0"
	oidCommunity  = ".1.3.6.1.6.3.18.1.4.0"
	oidTrap       = ".1.3.6.1.6.3.1.1.4.1.0"
	oidAgent      = ".1.3.6.1.6.3.18.1.3.0"
	oidEnterprise = ".1.3.6.1.6.3.1.1.4.3.0"
)

func Parse(input []byte) *types.Trap {

	inputStr := normalizeDquote(string(input))

	lines := strings.Split(inputStr, "\n")
	addr := lines[1]
	lines = lines[2:]

	trap := &types.Trap{}

	r, _ := regexp.Compile(`^\w+:\s+\[(\d+\.\d+\.\d+\.\d+)\]`)
	m := r.FindStringSubmatch(addr)
	if len(m) >= 2 {
		trap.Ipaddr = m[1]
	}

	trap.Variables = make(map[string]string)

	for _, str := range lines {
		kv := strings.Split(str, " ")
		if len(kv) >= 2 {
			oid := kv[0]
			val := strings.Join(kv[1:], " ")
			if len(oid) > 1 && oid[0] == '.' {
				trap.Variables[oid] = val
			}
		}
	}

	if val, ok := trap.Variables[oidCommunity]; ok {
		delete(trap.Variables, oidCommunity)
		trap.Community = val
	}

	if val, ok := trap.Variables[oidTrap]; ok {
		delete(trap.Variables, oidTrap)
		trap.Trap = val
	}

	if _, ok := trap.Variables[oidUpTime]; ok {
		delete(trap.Variables, oidUpTime)
	}

	if _, ok := trap.Variables[oidAgent]; ok {
		delete(trap.Variables, oidAgent)
	}

	if _, ok := trap.Variables[oidEnterprise]; ok {
		delete(trap.Variables, oidEnterprise)
	}

	return trap
}

func normalizeDquote(input string) string {
	r, _ := regexp.Compile(`"([^"]*)"`)
	return r.ReplaceAllStringFunc(input, func(s string) string {
		s = strings.Trim(s, `"`)
		s = strings.TrimSpace(s)
		s = strings.Replace(s, "\n", " ", -1)
		s = strings.Replace(s, "  ", " ", -1)
		return hexStringToRawString(s)
	})
}

func hexStringToRawString(str string) string {
	b := []byte{}
	for _, hex := range strings.Split(str, " ") {
		n, err := strconv.ParseUint(hex, 16, 8)
		if err != nil {
			return str
		}
		b = append(b, byte(n))
	}
	return string(b)
}
