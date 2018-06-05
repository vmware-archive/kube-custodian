package cleaner

import (
	log "github.com/sirupsen/logrus"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type jobUpdater struct {
	job *batchv1.Job
}

func (u *jobUpdater) Update(c *Common) error {
	_, err := c.clientset.BatchV1().Jobs(u.job.Namespace).Update(u.job)
	return err
}

func (u *jobUpdater) Meta() *metav1.ObjectMeta {
	return &u.job.ObjectMeta
}

func (u *jobUpdater) Delete(c *Common) error {
	return c.clientset.BatchV1().Jobs(u.job.Namespace).Delete(u.job.Name, &metav1.DeleteOptions{})
}

// updateJobs ...
func (c *Common) updateJobs(namespace string) (int, int, error) {
	updatedCount := 0
	deletedCount := 0
	jobs, err := c.clientset.BatchV1().Jobs(namespace).List(metav1.ListOptions{})
	if err != nil {
		log.Errorf("List jobs: %v", err)
		return updatedCount, deletedCount, err
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

		updCnt, delCnt := c.updateState(&jobUpdater{job: &job})
		updatedCount += updCnt
		deletedCount += delCnt
	}
	return updatedCount, deletedCount, nil
}
