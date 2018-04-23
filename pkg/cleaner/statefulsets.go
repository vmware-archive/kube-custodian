package cleaner

import (
	log "github.com/sirupsen/logrus"
	"k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type stsUpdater struct {
	sts *v1beta1.StatefulSet
}

func (u *stsUpdater) Update(c *Common) error {
	_, err := c.clientset.AppsV1beta1().StatefulSets(u.sts.Namespace).Update(u.sts)
	return err
}

func (u *stsUpdater) Delete(c *Common) error {
	return c.clientset.AppsV1beta1().StatefulSets(u.sts.Namespace).Delete(u.sts.Name, &metav1.DeleteOptions{})
}

func (u *stsUpdater) Meta() *metav1.ObjectMeta {
	return &u.sts.ObjectMeta
}

// updateStatefulSets ...
func (c *Common) updateStatefulSets(namespace string) (int, int, error) {
	updatedCount := 0
	deletedCount := 0
	stss, err := c.clientset.AppsV1beta1().StatefulSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List StatefulSets: %v", err)
		return updatedCount, deletedCount, err
	}

	for _, sts := range stss.Items {
		log.Debugf("StatefulSet %s.%s ...", sts.Namespace, sts.Name)
		if c.skipFromMeta(&sts.ObjectMeta) {
			continue
		}

		log.Debugf("StatefulSet %s.%s about to be touched ...", sts.Namespace, sts.Name)

		updCnt, delCnt := c.updateState(&stsUpdater{sts: &sts})
		updatedCount += updCnt
		deletedCount += delCnt
	}
	return updatedCount, deletedCount, nil
}
