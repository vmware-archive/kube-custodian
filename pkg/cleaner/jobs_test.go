package cleaner

import (
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_DeleteJobs(t *testing.T) {
	job_obj := &batchv1.JobList{
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
			// will be skipped from its .*-system namespace
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

	pod_obj := &corev1.PodList{
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
	SetSystemNS("")
	// two non system Succeeded Jobs and one Pod
	clientset := fake.NewSimpleClientset(job_obj, pod_obj)
	count, err := DeleteJobs(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 3)

	// no one, as the 1st two now have the required label
	clientset = fake.NewSimpleClientset(job_obj, pod_obj)
	count, err = DeleteJobs(clientset, false, "", []string{"created_by"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 0)

	// only job1 in ns1 and its pod1
	clientset = fake.NewSimpleClientset(job_obj, pod_obj)
	count, err = DeleteJobs(clientset, false, "ns1", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 2)

	// 3of4 Jobs Succeeded (+ pod1), as sysNS has been overridden
	SetSystemNS(".*sYsTEM")
	clientset = fake.NewSimpleClientset(job_obj, pod_obj)
	count, err = DeleteJobs(clientset, false, "", []string{"xxx"})
	assertEqual(t, err, nil)
	assertEqual(t, count, 4)

}
