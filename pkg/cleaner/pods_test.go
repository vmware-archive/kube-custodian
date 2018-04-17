package cleaner

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_updatePodsCond(t *testing.T) {
	obj := &corev1.PodList{
		Items: []corev1.Pod{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod3",
					Namespace: "ns3",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubernetes-dashboard-deadbeefed-quack",
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

	t.Logf("Should update all pods")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err := c.updatePodsCond("",
		func(pod *corev1.Pod) bool {
			return true
		})
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 5)
	assertEqual(t, delCnt, 0)

	t.Logf("Should update a single filterIn() Pod")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updatePodsCond("",
		func(pod *corev1.Pod) bool {
			return pod.Labels["created_by"] == "foo"
		})
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 1)
	assertEqual(t, delCnt, 0)

	t.Logf("Should not update any pods")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updatePodsCond("",
		func(pod *corev1.Pod) bool {
			return false
		})
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 0)
}
