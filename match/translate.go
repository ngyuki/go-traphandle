package match

import (
	"errors"
	"os/exec"
	"strings"
)

func translateOidToRaw(oid string) (string, error) {

	suffix := ""
	for _, s := range []string{".**", ".*"} {
		if strings.HasSuffix(oid, s) {
			oid = strings.TrimSuffix(oid, s)
			suffix = s
			break
		}
	}

	raw, err := runSnmpTranslate(oid)
	if err != nil {
		return "", err
	}

	return raw + suffix, nil
}

func runSnmpTranslate(oid string) (string, error) {
	out, err := exec.Command("snmptranslate", "-On", oid).Output()
	if err != nil {
		return "", errors.New("snmptranslate -On " + oid + " ... " + err.Error())
	}
	return strings.TrimSpace(string(out)), nil
}
