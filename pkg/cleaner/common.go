package cleaner

import (
	"regexp"
)

const (
	SystemNS = ".*(-system|monitoring|logging|ingress)"
)

// SystemRE set from flags at cmd/delete.go
var SystemRE *regexp.Regexp

func init() {
	SystemRE = regexp.MustCompile(SystemNS)
}

func isSystemNS(namespace string) bool {
	return SystemRE.MatchString(namespace)
}
func SetSystemNS(namespace_re string) {
	SystemRE = regexp.MustCompile(namespace_re)
}
