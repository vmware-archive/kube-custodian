package cleaner

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type podUpdater struct {
	pod *corev1.Pod
}

func (u *podUpdater) Update(c *Common) error {
	_, err := c.clientset.CoreV1().Pods(u.pod.Namespace).Update(u.pod)
	return err
}

func (u *podUpdater) Delete(c *Common) error {
	return c.clientset.CoreV1().Pods(u.pod.Namespace).Delete(u.pod.Name, &metav1.DeleteOptions{})
}

func (u *podUpdater) Meta() *metav1.ObjectMeta {
	return &u.pod.ObjectMeta
}

// updatePods is main entry point from cmd/delete.go
func (c *Common) updatePods() (int, int, error) {
	return c.updatePodsCond(c.Namespace,
		func(pod *corev1.Pod) bool {
			if c.skipFromMeta(&pod.ObjectMeta) {
				return false
			}
			return true
		})
}

// updatePodsCond is passed a generic closure to select Pods to delete
func (c *Common) updatePodsCond(namespace string, filterIn func(*corev1.Pod) bool) (int, int, error) {

	updatedCount := 0
	deletedCount := 0
	pods, err := c.clientset.Core().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List pods: %v", err)
		return updatedCount, deletedCount, err
	}

	for _, pod := range pods.Items {
		log.Debugf("Pod %s.%s ...", pod.Namespace, pod.Name)
		if !filterIn(&pod) {
			continue
		}

		log.Debugf("Pod %s.%s about to be touched ...", pod.Namespace, pod.Name)

		updCnt, delCnt := c.updateState(&podUpdater{pod: &pod})
		updatedCount += updCnt
		deletedCount += delCnt
	}
	return updatedCount, deletedCount, nil
}
