package cleaner

import (
	"regexp"
)

// SystemNS has default "system" namespaces regexp
const (
	SystemNS = "kube-.*|.*(-system|monitoring|logging|ingress)"
)

var systemRE *regexp.Regexp

func init() {
	SetSystemNS("")
}

func isSystemNS(namespace string) bool {
	return systemRE.MatchString(namespace)
}

// SetSystemNS is used from cmd/delete.go flags
func SetSystemNS(namespaceRe string) {
	if namespaceRe != "" {
		systemRE = regexp.MustCompile(namespaceRe)
	} else {
		systemRE = regexp.MustCompile(SystemNS)
	}
}
