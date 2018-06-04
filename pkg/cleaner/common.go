package cleaner

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/bitnami-labs/kube-custodian/pkg/utils"
)

// Common represents flags set from cmd/, conversion of these
// into more convenient type(s), and runtime clientset
type Common struct {
	DryRun          bool
	Namespace       string
	SkipNamespaceRE string
	SkipLabels      []string
	TagTTL          string

	skipNamespaceRegexp *regexp.Regexp
	timeStamp           int64
	tagTTL              int64
	dryRunStr           string

	clientset kubernetes.Interface
}

// CommonDefaults has opinionated default values for Common
var CommonDefaults = &Common{
	SkipNamespaceRE: "kube-.*|.*(-system|monitoring|logging|ingress)",
	SkipLabels:      []string{"created_by"},
	TagTTL:          "24h",
}

const (
	kubeCustodianAnnotationTime = "kube-custodian.bitnami.com/expiration-time"
)

type updater interface {
	Update(*Common) error
	Delete(*Common) error
	Meta() *metav1.ObjectMeta
}

// Init initializes Common obj with ready-to-use values,
// must be called by callers before Run()
func (c *Common) Init(clientset kubernetes.Interface) {
	var err error
	c.skipNamespaceRegexp = regexp.MustCompile(c.SkipNamespaceRE)
	c.timeStamp = time.Now().Unix()
	c.dryRunStr = map[bool]string{true: "[dry-run]", false: ""}[c.DryRun]
	tagTTL, err := time.ParseDuration(c.TagTTL)
	if err != nil {
		log.Fatalf("Failed for parse %q as time.Duration", c.TagTTL)
	}
	c.tagTTL = int64(tagTTL / time.Second)
	c.clientset = clientset
}

// Run is main entry point for this package, will loop over all namespaces
// skip those matching SkipLabels or SkipNamespaceRE, then call update<Type>()
// on each
func (c Common) Run() (int, int, int) {
	var nss *corev1.NamespaceList
	var err error
	var updatedCount, deletedCount, errorCount int

	if c.Namespace != "" {
		ns, err := c.clientset.Core().Namespaces().Get(c.Namespace, metav1.GetOptions{})
		if err != nil {
			log.Errorf("Error getting namespace: %s", c.Namespace)
			return 0, 0, 1
		}
		nss = &corev1.NamespaceList{Items: []corev1.Namespace{*ns}}

	} else {
		nss, err = c.clientset.Core().Namespaces().List(metav1.ListOptions{})
		if err != nil {
			log.Errorf("List namespaces: %v", err)
			return 0, 0, 1
		}
	}

	for _, ns := range nss.Items {
		if c.skipNamespaceRegexp.MatchString(ns.Name) {
			log.Debugf("Namespace %s skipped", ns.Name)
			continue
		}
		if c.skipFromMeta(&ns.ObjectMeta) {
			continue
		}
		log.Debugf("Scanning namespace: %s", ns.Name)
		updateFunctions := []func(string) (int, int, error){
			c.updateDeployments,
			c.updateStatefulSets,
			c.updateJobs,
			c.updatePods,
		}
		for _, updateFunc := range updateFunctions {
			updCnt, delCnt, err := updateFunc(ns.Name)
			if err != nil {
				errorCount++
			}
			updatedCount += updCnt
			deletedCount += delCnt
		}
	}
	return updatedCount, deletedCount, errorCount
}

func (c *Common) skipFromMeta(meta *metav1.ObjectMeta) bool {
	skipIt := false
	switch {
	case c.skipNamespaceRegexp.MatchString(meta.Namespace):
		log.Debugf("%s.%s skipped from meta.Namespace", meta.Name, meta.Namespace)
		skipIt = true
	case utils.LabelsSubSet(meta.Labels, c.SkipLabels):
		log.Debugf("%s.%s skipped from meta.Labels", meta.Name, meta.Namespace)
		skipIt = true
	}
	return skipIt
}

// updateState returns number of objects (updated, deleted)
func (c *Common) updateState(updater updater) (int, int) {
	objMeta := updater.Meta()
	fqName := fmt.Sprintf("%s.%s", objMeta.Name, objMeta.Namespace)
	updatedCount := 0
	deletedCount := 0
	annotations := objMeta.GetAnnotations()
	if valueStr, found := annotations[kubeCustodianAnnotationTime]; found {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			log.Errorf("%s: failed to convert %s to integer", fqName, valueStr)
		}
		expiredSecs := c.timeStamp - (value + c.tagTTL)
		log.Debugf("%s already has annotation %s: %s, will expire in %.2f hours",
			fqName, kubeCustodianAnnotationTime, valueStr, -float64(expiredSecs)/3600)
		if expiredSecs >= 0 {
			log.Debugf("%s%s TTL expired %d seconds ago, deleting",
				c.dryRunStr, fqName, expiredSecs)
			if !c.DryRun {
				if err := updater.Delete(c); err != nil {
					log.Errorf("failed to delete %s with error: %v", fqName, err)
				} else {
					deletedCount++
				}
			}
		}
	} else {
		timeStampStr := fmt.Sprintf("%d", c.timeStamp)
		log.Debugf("%s%s creating annotation %s: %s",
			c.dryRunStr, fqName, kubeCustodianAnnotationTime, timeStampStr)
		if !c.DryRun {
			metav1.SetMetaDataAnnotation(objMeta,
				kubeCustodianAnnotationTime, timeStampStr)
			if err := updater.Update(c); err != nil {
				log.Errorf("failed to update %s with error: %v", fqName, err)
			} else {
				updatedCount++
			}
		}
	}
	return updatedCount, deletedCount
}
