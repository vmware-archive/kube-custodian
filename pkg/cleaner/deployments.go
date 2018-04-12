package cleaner

import (
	log "github.com/sirupsen/logrus"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/jjo/kube-custodian/pkg/utils"
)

// DeleteDeployments ...
func DeleteDeployments(clientset kubernetes.Interface, dryRun bool, namespace string, requiredLabels []string) (int, error) {

	count := 0
	deploys, err := clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List deploys: %v", err)
		return count, err
	}

	deploysArray := make([]appsv1.Deployment, 0)

	for _, deploy := range deploys.Items {
		log.Debugf("Deploy %s.%s ...", deploy.Namespace, deploy.Name)
		if isSystemNS(deploy.Namespace) {
			log.Debugf("Deploy %q in system NS, skipping", deploy.Name)
			continue
		}

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

		log.Infof("%sDeleting Deploy %s.%s ... ", dryRunStr, deploy.Namespace, deploy.Name)
		if !dryRun {
			if err := clientset.AppsV1().Deployments(deploy.Namespace).Delete(deploy.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("failed to delete Deploy: %v", err)
				continue
			}
			count++
		}

	}
	return count, nil
}
