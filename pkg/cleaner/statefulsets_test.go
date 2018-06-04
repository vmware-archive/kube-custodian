package cleaner

import (
	"testing"

	"k8s.io/api/apps/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_updateStatefulSets(t *testing.T) {
	nss := &corev1.NamespaceList{
		Items: []corev1.Namespace{
			{ObjectMeta: metav1.ObjectMeta{Name: "ns1"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "ns2"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "ns3"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "kube-system"}},
			{ObjectMeta: metav1.ObjectMeta{Name: "monitoring"}},
		},
	}
	obj := &v1beta1.StatefulSetList{
		Items: []v1beta1.StatefulSet{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts3",
					Namespace: "ns3",
				},
			},
			// sysNS will be skipped:
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fatDb",
					Namespace: "kube-system",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prometheus",
					Namespace: "monitoring",
				},
			},
		},
	}
	var c Common

	t.Logf("Should update all sts except those in kube-system and monitoring NS")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(nss, obj))
	updCnt, delCnt, errCnt := c.Run()
	assertEqual(t, errCnt, 0)
	assertEqual(t, updCnt, 3)
	assertEqual(t, delCnt, 0)

	t.Logf("Should update only sts in ns1")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Namespace = "ns1"
	c.Init(fake.NewSimpleClientset(nss, obj))
	updCnt, delCnt, errCnt = c.Run()
	assertEqual(t, errCnt, 0)
	assertEqual(t, updCnt, 1)
	assertEqual(t, delCnt, 0)

	t.Logf("Should update only one sts, as the other two candidates have the 'created_by' label")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(nss, obj))
	updCnt, delCnt, errCnt = c.Run()
	assertEqual(t, errCnt, 0)
	assertEqual(t, updCnt, 1)
	assertEqual(t, delCnt, 0)

	t.Logf("Should update all sts, as namespaceRE and skipLabels don't match any")
	c = *CommonDefaults
	c.SkipNamespaceRE = ".*sYsTEM"
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(nss, obj))
	updCnt, delCnt, errCnt = c.Run()
	assertEqual(t, errCnt, 0)
	assertEqual(t, updCnt, 5)
	assertEqual(t, delCnt, 0)
	t.Logf("... second call should not touch anything")
	updCnt, delCnt, _ = c.Run()
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 0)
	t.Logf("... another call with a zero TTL should delete all marked ones")
	c.tagTTL = 0
	updCnt, delCnt, _ = c.Run()
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 5)
}
