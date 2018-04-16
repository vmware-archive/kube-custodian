package cleaner

import (
	"testing"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_DeletePodsCond(t *testing.T) {
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

	t.Logf("Should delete all pods")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	count, err := c.DeletePodsCond("",
		func(pod *corev1.Pod) bool {
			return true
		})
	assertEqual(t, err, nil)
	assertEqual(t, count, 5)

	t.Logf("Should delete a single filterIn() Pod")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	count, err = c.DeletePodsCond("",
		func(pod *corev1.Pod) bool {
			return pod.Labels["created_by"] == "foo"
		})
	assertEqual(t, err, nil)
	assertEqual(t, count, 1)

	t.Logf("Should not delete any pods")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	count, err = c.DeletePodsCond("",
		func(pod *corev1.Pod) bool {
			return false
		})
	assertEqual(t, err, nil)
	assertEqual(t, count, 0)
}
