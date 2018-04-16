package cleaner

import (
	"testing"

	"k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_DeleteStatefulSets(t *testing.T) {
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

	t.Logf("Should delete all sts except those in kube-system and monitoring NS")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(obj))
	count, err := c.DeleteStatefulSets()
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	t.Logf("Should delete only sts in ns1")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Namespace = "ns1"
	c.Init(fake.NewSimpleClientset(obj))
	count, err = c.DeleteStatefulSets()
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	t.Logf("Should delete only one sts, as the other two candidates have the 'created_by' label")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	count, err = c.DeleteStatefulSets()
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	t.Logf("Should delete all sts, as namespaceRE and skipLabels don't match any")
	c = *CommonDefaults
	c.SkipNamespaceRE = ".*sYsTEM"
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(obj))
	count, err = c.DeleteStatefulSets()
	assertEqual(t, err, nil)
	assertEqual(t, count, 5)
}
