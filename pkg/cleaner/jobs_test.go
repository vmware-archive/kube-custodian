package cleaner

import (
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_updateJobs(t *testing.T) {
	obj := &batchv1.JobList{
		Items: []batchv1.Job{
			{
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
			{
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
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "job3",
					Namespace: "ns3",
				},
				Status: batchv1.JobStatus{
					Succeeded: 0,
				},
			},
			{
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

	var c Common

	t.Logf("Should update all jobs except those in kube-system and monitoring NS")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err := c.updateJobs()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 2)
	assertEqual(t, delCnt, 0)

	t.Logf("Should update only the jobs in ns1")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Namespace = "ns1"
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updateJobs()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 1)
	assertEqual(t, delCnt, 0)

	t.Logf("Should not update any, as the first two have the required label")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updateJobs()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 0)

	t.Logf("Should update all jobs, as namespaceRE and skipLabels don't match any")
	c = *CommonDefaults
	c.SkipNamespaceRE = ".*sYsTEM"
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updateJobs()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 3)
	assertEqual(t, delCnt, 0)
	t.Logf("... second call should not touch anything")
	updCnt, delCnt, err = c.updateJobs()
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 0)
	t.Logf("... another call with a zero TTL should delete all marked ones")
	c.tagTTL = 0
	updCnt, delCnt, err = c.updateJobs()
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 3)
}
