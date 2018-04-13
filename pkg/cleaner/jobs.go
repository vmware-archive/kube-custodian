package cleaner

import (
	log "github.com/sirupsen/logrus"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/jjo/kube-custodian/pkg/utils"
)

const (
	kubeJobNameLabel = "job-name"
)

// DeleteJobs ...
func DeleteJobs(clientset kubernetes.Interface, dryRun bool, namespace string, requiredLabels []string) (int, error) {
	jobs, err := clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})

	count := 0
	if err != nil {
		log.Errorf("List jobs: %v", err)
		return count, err
	}

	jobArray := make([]batchv1.Job, 0)

	for _, job := range jobs.Items {
		if isSystemNS(job.Namespace) {
			log.Debugf("Job %q in system NS, skipping", job.Name)
			continue
		}
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
		return count, err
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

		log.Infof("%sDeleting Job %s.%s ...", dryRunStr, job.Namespace, job.Name)
		if !dryRun {
			if err := clientset.BatchV1().Jobs(job.Namespace).Delete(job.Name, &metav1.DeleteOptions{}); err != nil {
				log.Errorf("failed to delete Job: %v", err)
				continue
			}
			count++
		}

		for _, pod := range jobPods {
			log.Infof("%s  Deleting Pod %s.%s ...", dryRunStr, pod.Namespace, pod.Name)
			if !dryRun {
				if err := clientset.CoreV1().Pods(pod.Namespace).Delete(pod.Name, &metav1.DeleteOptions{}); err != nil {
					log.Errorf("failed to delete Pod: %v", err)
					continue
				}
				count++
			}
		}
	}
	return count, nil
}
