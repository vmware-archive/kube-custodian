package cleaner

import (
	log "github.com/sirupsen/logrus"

	appsv1 "k8s.io/api/apps/v1"
	// corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/jjo/kube-custodian/pkg/utils"
)

const (
	kubeJobNameLabel = "job-name"
)

// DeleteJobs ...
func DeleteDeployments(clientset *kubernetes.Clientset, dryRun bool, namespace string, requiredLabels []string) error {
	deploys, err := clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List deploys: %v", err)
		return err
	}

	if len(requiredLabels) < 1 {
		log.Fatal("At least one required-label is needed")
	}
	log.Infof("Required labels: %v ...", requiredLabels)
	if err != nil {
		log.Error(err)
		return err
	}

	deploysArray := make([]appsv1.Deployment, 0)

	for _, deploy := range deploys.Items {
		log.Debugf("Deploy %q ...", deploy.Name)

		if utils.LabelsSubSet(deploy.Labels, requiredLabels) {
			log.Debugf("Deploy %q has required labels (%v), skipping", deploy.Name, deploy.Labels)
			continue
		}

		log.Debugf("Deploy %q missing required labels, will be marked for deletion", deploy.Name)
		deploysArray = append(deploysArray, deploy)
	}

	dryRunStr := map[bool]string{true: "[dry-run]", false: ""}[dryRun]
	for _, deploy := range deploysArray {
		log.Debugf("Deploy %q about to be deleted", deploy.Name)

		log.Infof("Deleting Deploy %s.%s ... %s", deploy.Namespace, deploy.Name, dryRunStr)
		if !dryRun {
			if err := clientset.AppsV1().Deployments(deploy.Namespace).Delete(deploy.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("failed to delete Deploy: %v", err)
			}
		}

	}
	return nil
}
