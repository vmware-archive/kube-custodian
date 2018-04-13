package cleaner

import (
	log "github.com/sirupsen/logrus"

	"k8s.io/api/apps/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// DeleteStatefulSets ...
func DeleteStatefulSets(clientset kubernetes.Interface, dryRun bool, namespace string) (int, error) {

	count := 0
	stss, err := clientset.AppsV1beta1().StatefulSets(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List StatefulSets: %v", err)
		return count, err
	}

	stsArray := make([]v1beta1.StatefulSet, 0)

	for _, sts := range stss.Items {
		log.Debugf("StatefulSet %s.%s ...", sts.Namespace, sts.Name)
		if skipFromMeta(&sts.ObjectMeta) {
			continue
		}

		log.Debugf("StatefulSet %q marked for deletion", sts.Name)
		stsArray = append(stsArray, sts)
	}

	dryRunStr := map[bool]string{true: "[dry-run]", false: ""}[dryRun]
	for _, sts := range stsArray {
		log.Debugf("StatefulSet %q about to be deleted", sts.Name)

		log.Infof("%sDeleting StatefulSet %s.%s ... ", dryRunStr, sts.Namespace, sts.Name)
		if !dryRun {
			if err := clientset.AppsV1beta1().StatefulSets(sts.Namespace).Delete(sts.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("failed to delete StatefulSet: %v", err)
				continue
			}
			count++
		}
	}
	return count, nil
}
