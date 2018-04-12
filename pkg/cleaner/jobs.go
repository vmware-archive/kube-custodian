package cleaner

import (
	log "github.com/sirupsen/logrus"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/jjo/kube-custodian/pkg/utils"
)

// DeleteJobs ...
func DeleteJobs(clientset *kubernetes.Clientset, dryRun bool, namespace string, requiredLabels []string) error {
	jobs, err := clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List jobs: %v", err)
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

	jobArray := make([]batchv1.Job, 0)

	for _, job := range jobs.Items {
		log.Debugf("Job %q ...", job.Name)
		if job.Status.Succeeded == 0 {
			log.Debugf("Job %q not finished, skipping", job.Name)
			continue
		}

		if utils.LabelsSubSet(job.Labels, requiredLabels) {
			log.Debugf("Job %q has required labels (%v), skipping", job.Name, job.Labels)
			continue
		}

		log.Debugf("Job %q missing required labels, will be marked for deletion", job.Name)
		jobArray = append(jobArray, job)
	}

	pods, err := clientset.Core().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List pods: %v", err)
		return err
	}

	jobPods := []corev1.Pod{}

	for _, pod := range pods.Items {
		if pod.Labels[kubeJobNameLabel] == "" {
			continue
		}

		if !(pod.Status.Phase == corev1.PodSucceeded ||
			pod.Status.Phase == corev1.PodFailed) {
			log.Debugf("Pod %q still running, skipping", pod.Name)
			continue
		}

		jobPods = append(jobPods, pod)
	}

	dryRunStr := map[bool]string{true: "[dry-run]", false: ""}[dryRun]
	for _, job := range jobArray {
		log.Debugf("Job %q about to be deleted", job.Name)

		log.Infof("Deleting Job %s.%s ... %s", job.Namespace, job.Name, dryRunStr)
		if !dryRun {
			if err := clientset.BatchV1().Jobs(job.Namespace).Delete(job.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("failed to delete Job: %v", err)
			}
		}

		for _, pod := range jobPods {
			log.Infof("  Deleting Pod %s.%s ... %s", pod.Namespace, pod.Name, dryRunStr)
			if !dryRun {
				if err := clientset.Core().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
					log.Errorf("failed to delete Pod: %v", err)
				}
			}
		}
	}
	return nil
}
