package utils

import (
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_LabelsSubSet(t *testing.T) {
	assertEqual(t, LabelsSubSet(map[string]string{"a": "x", "b": "y"}, []string{"a"}), true)
	assertEqual(t, LabelsSubSet(map[string]string{"a": "x", "b": "y"}, []string{"q"}), false)
}
