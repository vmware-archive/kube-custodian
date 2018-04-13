package cleaner

import (
	"regexp"
)

// SystemNS has default "system" namespaces regexp
const (
	SystemNS = ".*(-system|monitoring|logging|ingress)"
)

var systemRE *regexp.Regexp

func init() {
	systemRE = regexp.MustCompile(SystemNS)
}

func isSystemNS(namespace string) bool {
	return systemRE.MatchString(namespace)
}

// SetSystemNS is used from cmd/delete.go flags
func SetSystemNS(namespaceRe string) {
	systemRE = regexp.MustCompile(namespaceRe)
}
