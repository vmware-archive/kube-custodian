package cleaner

import (
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_DeleteJobs(t *testing.T) {
	jobObj := &batchv1.JobList{
		Items: []batchv1.Job{
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
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
			batchv1.Job{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job4",
					Namespace: "kube-system",
				},
				Status: batchv1.JobStatus{
					Succeeded: 1,
				},
			},
		},
	}

	podObj := &corev1.PodList{
		Items: []corev1.Pod{
			corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "pod1",
					Namespace: "ns1",
					Labels: map[string]string{
						kubeJobNameLabel: "job1",
					},
				},
				Status: corev1.PodStatus{
					Phase: corev1.PodSucceeded,
				},
			},
		},
	}
	t.Logf("Should delete all jobs except those in kube-system and monitoring NS")
	SetSkipMeta("", []string{"xxx"})
	clientset := fake.NewSimpleClientset(jobObj, podObj)
	count, err := DeleteJobs(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	t.Logf("Should delete only the jobs in ns1 and its pod")
	clientset = fake.NewSimpleClientset(jobObj, podObj)
	count, err = DeleteJobs(clientset, false, "ns1")
	assertEqual(t, err, nil)
	assertEqual(t, count, 2)

	SetSkipMeta("", nil)
	t.Logf("Should not delete any, as the first two have the required label")
	clientset = fake.NewSimpleClientset(jobObj, podObj)
	count, err = DeleteJobs(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 0)

	t.Logf("Should delete all jobs, as namespaceRE and skipLabels don't match any")
	SetSkipMeta(".*sYsTEM", []string{"xxx"})
	clientset = fake.NewSimpleClientset(jobObj, podObj)
	count, err = DeleteJobs(clientset, false, "")
	assertEqual(t, err, nil)
	assertEqual(t, count, 4)

}
