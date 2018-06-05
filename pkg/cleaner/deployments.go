package cleaner

import (
	log "github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type deployUpdater struct {
	deploy *appsv1.Deployment
}

func (u *deployUpdater) Update(c *Common) error {
	_, err := c.clientset.AppsV1().Deployments(u.deploy.Namespace).Update(u.deploy)
	return err
}

func (u *deployUpdater) Delete(c *Common) error {
	return c.clientset.AppsV1().Deployments(u.deploy.Namespace).Delete(u.deploy.Name, &metav1.DeleteOptions{})
}

func (u *deployUpdater) Meta() *metav1.ObjectMeta {
	return &u.deploy.ObjectMeta
}

// updateDeployments ...
func (c *Common) updateDeployments(namespace string) (int, int, error) {
	updatedCount := 0
	deletedCount := 0
	deploys, err := c.clientset.AppsV1().Deployments(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List deploys: %v", err)
		return updatedCount, deletedCount, err
	}

	for _, deploy := range deploys.Items {
		log.Debugf("Deploy %s.%s ...", deploy.Namespace, deploy.Name)
		if c.skipFromMeta(&deploy.ObjectMeta) {
			continue
		}

		log.Debugf("Deploy %s.%s about to be touched ...", deploy.Namespace, deploy.Name)
		updCnt, delCnt := c.updateState(&deployUpdater{deploy: &deploy})
		updatedCount += updCnt
		deletedCount += delCnt
	}
	return updatedCount, deletedCount, nil
}
