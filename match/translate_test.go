package match

import (
	"testing"
)

func TestTranslateOidToRaw_ok(t *testing.T) {
	test := func(oid string, exp string) {
		act, err := translateOidToRaw(oid)
		if err != nil {
			t.Error(err)
		}
		if act != exp {
			t.Fatalf("translateOidToRaw(%v) must be %v ... %v", oid, exp, act)
		}
	}
	test("SNMPv2-MIB::sysContact", ".1.3.6.1.2.1.1.4")
	test("SNMPv2-MIB::sysContact.0", ".1.3.6.1.2.1.1.4.0")
	test("SNMPv2-MIB::sysContact.*", ".1.3.6.1.2.1.1.4.*")
	test("SNMPv2-MIB::sysContact.**", ".1.3.6.1.2.1.1.4.**")
}

func TestTranslateOidToRaw_err(t *testing.T) {
	oid := "invalid-mib::oid"
	if _, err := translateOidToRaw(oid); err == nil {
		t.Fatalf("translateOidToRaw err is nil")
	}
}
