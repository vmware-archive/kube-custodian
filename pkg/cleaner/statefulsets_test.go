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
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
			},
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
			},
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "sts3",
					Namespace: "ns3",
				},
			},
			// sysNS will be skipped:
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "fatDb",
					Namespace: "kube-system",
				},
			},
			v1beta1.StatefulSet{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "prometheus",
					Namespace: "monitoring",
				},
			},
		},
	}
	t.Logf("Should delete all sts except those in kube-system and monitoring NS")
	SetSkipMeta("", []string{"xxx"})
	clientset := fake.NewSimpleClientset(obj)
	count, err := DeleteStatefulSets(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	t.Logf("Should delete only sts in ns1")
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteStatefulSets(clientset, false, "ns1")
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	t.Logf("Should delete only one sts, as the other two candidates have the 'created_by' label")
	SetSkipMeta("", nil)
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteStatefulSets(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	t.Logf("Should delete all sts, as namespaceRE and skipLabels don't match any")
	SetSkipMeta(".*sYsTEM", []string{"xxx"})
	clientset = fake.NewSimpleClientset(obj)
	count, err = DeleteStatefulSets(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 5)
}
