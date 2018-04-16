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
	var c Common
	c = *CommonDefaults
	c.Init(nil)
	t.Logf("Should skip all resources in 'system' namespaces")
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-system"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-foo"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "foo-system"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "monitoring"}), true)
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "bob"}), false)
	c = *CommonDefaults
	c.SkipNamespaceRE = "xyz"
	c.Init(nil)
	t.Logf("Should match resources from changed --skip-namespace-re")
	assertEqual(t, c.skipFromMeta(&metav1.ObjectMeta{Namespace: "kube-system"}), false)
}
