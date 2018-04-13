package cleaner

import (
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_SysNS(t *testing.T) {
	assertEqual(t, isSystemNS("kube-system"), true)
}
