package cleaner

import (
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeletePods is main entry point from cmd/delete.go
func (c *Common) DeletePods() (int, error) {
	return c.DeletePodsCond(c.Namespace,
		func(pod *corev1.Pod) bool {
			if c.skipFromMeta(&pod.ObjectMeta) {
				return false
			}
			return true
		})
}

// DeletePodsCond is passed a generic closure to select Pods to delete
func (c *Common) DeletePodsCond(namespace string, filterIn func(*corev1.Pod) bool) (int, error) {

	count := 0
	pods, err := c.clientset.Core().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List pods: %v", err)
		return count, err
	}

	for _, pod := range pods.Items {
		log.Debugf("Pod %s.%s ...", pod.Namespace, pod.Name)
		if !filterIn(&pod) {
			continue
		}

		log.Debugf("Pod %s.%s about to be touched ...", pod.Namespace, pod.Name)

		count += c.updateState(
			func() error {
				_, err := c.clientset.CoreV1().Pods(pod.Namespace).Update(&pod)
				return err
			},
			func() error {
				return c.clientset.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{})
			},
			&pod.ObjectMeta)
	}
	return count, nil
}
