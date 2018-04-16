package cleaner

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteStatefulSets ...
func (c *Common) DeleteStatefulSets() (int, error) {

	count := 0
	stss, err := c.clientset.AppsV1beta1().StatefulSets(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List StatefulSets: %v", err)
		return count, err
	}

	for _, sts := range stss.Items {
		log.Debugf("StatefulSet %s.%s ...", sts.Namespace, sts.Name)
		if c.skipFromMeta(&sts.ObjectMeta) {
			continue
		}

		log.Debugf("StatefulSet %s.%s about to be touched ...", sts.Namespace, sts.Name)

		count += c.updateState(
			func() error {
				_, err := c.clientset.AppsV1beta1().StatefulSets(sts.Namespace).Update(&sts)
				return err
			},
			func() error {
				return c.clientset.AppsV1beta1().StatefulSets(sts.Namespace).Delete(sts.Name, &metav1.DeleteOptions{})
			},
			&sts.ObjectMeta)
	}
	return count, nil
}
