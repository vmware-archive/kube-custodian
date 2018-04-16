package cleaner

import (
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteJobs ...
func (c *Common) DeleteJobs() (int, error) {

	count := 0
	jobs, err := c.clientset.BatchV1().Jobs(c.Namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List jobs: %v", err)
		return count, err
	}

	for _, job := range jobs.Items {
		log.Debugf("Job %s.%s ...", job.Namespace, job.Name)
		if c.skipFromMeta(&job.ObjectMeta) {
			continue
		}
		if job.Status.Succeeded == 0 {
			log.Debugf("Job %q not finished, skipping", job.Name)
			continue
		}
		log.Debugf("Job %s.%s about to be touched ...", job.Namespace, job.Name)

		count += c.updateState(
			func() error {
				_, err := c.clientset.BatchV1().Jobs(job.Namespace).Update(&job)
				return err
			},
			func() error {
				return c.clientset.BatchV1().Jobs(job.Namespace).Delete(job.Name, &metav1.DeleteOptions{})
			},
			&job.ObjectMeta)
	}
	return count, nil
}
