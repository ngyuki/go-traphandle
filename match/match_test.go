package match

import (
	"testing"

	"github.com/ngyuki/go-traphandle/config"
	"github.com/ngyuki/go-traphandle/types"
)

func TestMatchTrapByIpAddr(t *testing.T) {

	test := func(data string, match string, expect bool) {
		trap := &types.Trap{}
		if len(data) > 0 {
			trap.Ipaddr = data
		}
		cfg := &config.MatchConfig{}
		if len(match) > 0 {
			cfg.Ipaddr = match
		}
		_, actual := newMatch(cfg).Match(trap)

		if actual != expect {
			t.Errorf("MatchTrap(%v, %v) must be return %v ... %v", cfg, trap, expect, actual)
		}
	}

	test("192.168.2.123", "192.168.2.0/24", true)
	test("192.168.3.123", "192.168.2.0/24", false)
	test("192.168.2.123", "192.168.2.123/32", true)
	test("192.168.2.124", "192.168.2.123/32", false)

	test("", "", true)
	test("192.168.2.123", "", true)
	test("", "192.168.3.0/24", false)
	test("", "192.168.2.123/32", false)
}

func TestMatchTrapByCommunity(t *testing.T) {

	test := func(data string, match string, expect bool) {
		trap := &types.Trap{}
		if len(data) > 0 {
			trap.Community = data
		}
		cfg := &config.MatchConfig{}
		if len(match) > 0 {
			cfg.Community = match
		}
		_, actual := newMatch(cfg).Match(trap)
		if actual != expect {
			t.Errorf("MatchTrap(%v, %v) must be return %v ... %v", cfg, trap, expect, actual)
		}
	}

	test("aaa", "aaa", true)
	test("aaa", "bbb", false)
	test("aaa", "", true)
	test("", "aaa", false)
}

func TestMatchTrapByTrap(t *testing.T) {

	test := func(data string, match string, expect bool) {
		trap := &types.Trap{}
		if len(data) > 0 {
			trap.Trap = data
		}
		cfg := &config.MatchConfig{}
		if len(match) > 0 {
			cfg.Trap = match
		}
		_, actual := newMatch(cfg).Match(trap)
		if actual != expect {
			t.Errorf("MatchTrap(%v, %v) must be return %v ... %v", cfg, trap, expect, actual)
		}
	}

	test(".1.3.6.1.2.1.1.4", "SNMPv2-MIB::sysContact", true)
	test(".1.3.6.1.2.1.1.4", "SNMP-COMMUNITY-MIB::snmpTrapCommunity.0", false)
	test(".1.3.6.1.2.1.1.4", "", true)
	test("", "SNMPv2-MIB::sysContact", false)
}

func TestMatchTrapByVariables(t *testing.T) {

	cfg := &config.MatchConfig{}
	cfg.Bindings = map[string]string{
		"val1": ".1.2.3.0",
		"val2": ".1.2.4.*",
		"val3": ".1.2.5.**",
	}
	cfg.Conditions = map[string]config.ConditionConfig{
		"val1": {
			"eq": {"1", "2", "3"},
		},
		"val2": {
			"eq": {"4"},
		},
		"val3": {
			"eq": {"5"},
		},
	}

	test := func(data map[string]string, expect bool) {
		trap := &types.Trap{}
		trap.Variables = data
		_, actual := newMatch(cfg).Match(trap)
		if actual != expect {
			t.Errorf("MatchTrap(cfg, %v) must be return %v ... %v", trap, expect, actual)
		}
	}

	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9": "4", ".1.2.5.9": "5"}, true)
	test(map[string]string{".1.2.3.9": "1", ".1.2.4.9": "4", ".1.2.5.9": "5"}, false)
	test(map[string]string{".1.2.3.0": "2", ".1.2.4.9": "4", ".1.2.5.9": "5"}, true)
	test(map[string]string{".1.2.3.0": "3", ".1.2.4.9": "4", ".1.2.5.9": "5"}, true)

	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9": "4", ".1.2.5.9": "5"}, true)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9": "0", ".1.2.5.9": "5"}, false)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.99": "4", ".1.2.5.9": "5"}, true)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.999": "4", ".1.2.5.9": "5"}, true)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9999": "4", ".1.2.5.9": "5"}, true)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9.9": "4", ".1.2.5.9": "5"}, false)

	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9": "4", ".1.2.5.9": "5"}, true)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9": "4", ".1.2.5.9": "0"}, false)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9": "4", ".1.2.5.9.9": "5"}, true)
	test(map[string]string{".1.2.3.0": "1", ".1.2.4.9": "4", ".1.2.5.9.9.9": "5"}, true)
}
