package cleaner

import (
	"fmt"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_CommonSkipMeta(t *testing.T) {
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

func Test_CommonUpdater(t *testing.T) {
	var c Common
	c = *CommonDefaults
	pod := corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "pod1",
			Namespace: "ns1",
		},
	}
	t.Logf("Should add %s annotation with timeStamp value", kubeCustodianAnnotationTime)
	c.Init(fake.NewSimpleClientset(&pod))
	u := &podUpdater{pod: &pod}
	uCnt, dCnt := c.updateState(u)
	assertEqual(t, uCnt, 1)
	assertEqual(t, dCnt, 0)
	assertEqual(t, pod.ObjectMeta.Annotations[kubeCustodianAnnotationTime], fmt.Sprintf("%d", c.timeStamp))
}

func Test_CommonNSLabel(t *testing.T) {
	nss := &corev1.NamespaceList{
		Items: []corev1.Namespace{
			{ObjectMeta: metav1.ObjectMeta{Name: "ns1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "ns2"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "ns3", Labels: map[string]string{"created_by": "sre"}}},
			{ObjectMeta: metav1.ObjectMeta{Name: "ns4"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "monitoring"}},
		},
	}
	dps := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "dp1", Namespace: "ns1"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "dp2", Namespace: "ns2"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "dp3", Namespace: "ns3"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "foo", Namespace: "foo-system"},
			},
		},
	}
	pods := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns1"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns2", Labels: map[string]string{"created_by": "bar"}},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "pod1", Namespace: "ns3"},
			},
			{
				ObjectMeta: metav1.ObjectMeta{Name: "bar", Namespace: "monitoring"},
			},
		},
	}
	var c Common

	t.Logf("Should update all deploys except those in with proper label either on its namespace or self")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(nss, dps, pods))
	updCnt, delCnt, errCnt := c.Run()
	assertEqual(t, errCnt, 0)
	assertEqual(t, updCnt, 3)
	assertEqual(t, delCnt, 0)
}
