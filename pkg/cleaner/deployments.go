package cleaner

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteDeployments ...
func (c *Common) DeleteDeployments() (int, error) {
	count := 0
	deploys, err := c.clientset.AppsV1().Deployments(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List deploys: %v", err)
		return count, err
	}

	for _, deploy := range deploys.Items {
		log.Debugf("Deploy %s.%s ...", deploy.Namespace, deploy.Name)
		if c.skipFromMeta(&deploy.ObjectMeta) {
			continue
		}

		log.Debugf("Deploy %s.%s about to be touched ...", deploy.Namespace, deploy.Name)

		count += c.updateState(
			func() error {
				_, err := c.clientset.AppsV1().Deployments(deploy.Namespace).Update(&deploy)
				return err
			},
			func() error {
				return c.clientset.AppsV1().Deployments(deploy.Namespace).Delete(deploy.Name, &metav1.DeleteOptions{})
			},
			&deploy.ObjectMeta)
	}
	return count, nil
}
