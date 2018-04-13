package cleaner

import (
	log "github.com/sirupsen/logrus"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeletePods is main entry point from cmd/delete.go
func DeletePods(clientset kubernetes.Interface, dryRun bool, namespace string) (int, error) {
	return DeletePodsCond(clientset, dryRun, namespace,
		func(pod *corev1.Pod) bool {
			if skipFromMeta(&pod.ObjectMeta) {
				return false
			}
			return true
		})
}

// DeletePodsCond is passed a generic closure to select Pods to delete
func DeletePodsCond(clientset kubernetes.Interface, dryRun bool, namespace string, filterIn func(*corev1.Pod) bool) (int, error) {

	count := 0
	pods, err := clientset.Core().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List pods: %v", err)
		return count, err
	}

	podsArray := []corev1.Pod{}

	for _, pod := range pods.Items {
		log.Debugf("Pod %s.%s ...", pod.Namespace, pod.Name)
		if !filterIn(&pod) {
			continue
		}

		log.Debugf("Pod %q marked for deletion", pod.Name)
		podsArray = append(podsArray, pod)
	}

	dryRunStr := map[bool]string{true: "[dry-run]", false: ""}[dryRun]
	for _, pod := range podsArray {
		log.Debugf("Pod %q about to be deleted", pod.Name)

		log.Infof("%s  Deleting Pod %s.%s ...", dryRunStr, pod.Namespace, pod.Name)
		if !dryRun {
			if err := clientset.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("failed to delete Pod: %v", err)
				continue
			}
			count++
		}
	}
	return count, nil
}
