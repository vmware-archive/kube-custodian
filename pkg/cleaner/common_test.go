package cleaner

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_SkipMeta(t *testing.T) {
	SetSkipMeta("", nil)
	assertEqual(t, skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-system"}), true)
	assertEqual(t, skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-foo"}), true)
	assertEqual(t, skipFromMeta(&metav1.ObjectMeta{Namespace: "foo-system"}), true)
	assertEqual(t, skipFromMeta(&metav1.ObjectMeta{Namespace: "monitoring"}), true)
	assertEqual(t, skipFromMeta(&metav1.ObjectMeta{Namespace: "bob"}), false)
	SetSkipMeta("xyz", nil)
	assertEqual(t, skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-system"}), false)
}
