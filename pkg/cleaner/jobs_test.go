package cleaner

import (
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
)

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%v != %v", a, b)
	}
}

func Test_DeleteJob(t *testing.T) {
	// 2of3 Jobs Succeeded
	obj := &batchv1.JobList{
		Items: []batchv1.Job{
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job1",
					Namespace: "ns1",
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job2",
					Namespace: "ns2",
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job3",
					Namespace: "ns3",
				},
				Status: batchv1.JobStatus{
					Succeeded: 0,
				},
			},
		},
	}
	clientset := fake.NewSimpleClientset(obj)

	count, err := DeleteJobs(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 2)
}
