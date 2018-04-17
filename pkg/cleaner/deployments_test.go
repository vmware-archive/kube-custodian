package cleaner

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

func Test_updateDeployments(t *testing.T) {
	obj := &appsv1.DeploymentList{
		Items: []appsv1.Deployment{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp1",
					Namespace: "ns1",
					Labels: map[string]string{
						"created_by": "bar",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp2",
					Namespace: "ns2",
					Labels: map[string]string{
						"created_by": "foo",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp3",
					Namespace: "ns3",
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "dp4",
					Namespace: "ns4",
					Annotations: map[string]string{
						kubeCustodianAnnotationTime: "0",
					},
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "kubernetes-dashboard",
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

	t.Logf("Should update all deploys except those in kube-system and monitoring NS")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err := c.updateDeployments()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 3)
	assertEqual(t, delCnt, 1)

	t.Logf("Should update only deploys in ns1")
	c = *CommonDefaults
	c.SkipLabels = []string{"xxx"}
	c.Namespace = "ns1"
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updateDeployments()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 1)
	assertEqual(t, delCnt, 0)

	t.Logf("Should update only one deploy, as the other two candidates have the 'created_by' label")
	c = *CommonDefaults
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updateDeployments()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 1)
	assertEqual(t, delCnt, 1)

	t.Logf("Should update all deploys, as namespaceRE and skipLabels don't match any")
	c = *CommonDefaults
	c.SkipNamespaceRE = ".*sYsTEM"
	c.SkipLabels = []string{"xxx"}
	c.Init(fake.NewSimpleClientset(obj))
	updCnt, delCnt, err = c.updateDeployments()
	assertEqual(t, err, nil)
	assertEqual(t, updCnt, 5)
	assertEqual(t, delCnt, 1)
	t.Logf("... second call should not touch anything")
	updCnt, delCnt, err = c.updateDeployments()
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 0)
	t.Logf("... another call with a zero TTL should delete all marked ones")
	c.tagTTL = 0
	updCnt, delCnt, err = c.updateDeployments()
	assertEqual(t, updCnt, 0)
	assertEqual(t, delCnt, 5)

}
