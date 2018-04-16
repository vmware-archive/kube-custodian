package cleaner

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	utils "github.com/jjo/kube-custodian/pkg/utils"
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

// Run is main entry point for this package
func (c Common) Run() {
	c.DeleteDeployments()
	c.DeleteStatefulSets()
	c.DeleteJobs()
	c.DeletePods()
}

func (c *Common) skipFromMeta(meta *metav1.ObjectMeta) bool {
	skipIt := false
	switch {
	case c.skipNamespaceRegexp.MatchString(meta.Namespace):
		log.Debugf("%s.%s skipped from meta.Namespace", meta.Name, meta.Namespace)
		skipIt = true
	case utils.LabelsSubSet(meta.Labels, c.SkipLabels):
		log.Debugf("%s.%s skipped from meta.Labels", meta.Name, meta.Labels)
		skipIt = true
	}
	return skipIt
}

func (c *Common) updateState(Update func() error, Delete func() error, objMeta *metav1.ObjectMeta) int {

	fqName := fmt.Sprintf("%s.%s", objMeta.Name, objMeta.Namespace)
	cnt := 0
	annotations := objMeta.GetAnnotations()
	if valueStr, found := annotations[kubeCustodianAnnotationTime]; found {
		value, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			log.Errorf("%s: failed to convert %s to integer", fqName, valueStr)
		}
		expiredSecs := c.timeStamp - (value + c.tagTTL)
		log.Debugf("%s already has annotation %s: %s, will expire in %.2f hours",
			fqName, kubeCustodianAnnotationTime, valueStr, -float64(expiredSecs)/3600)
		if expiredSecs > 0 {
			log.Debugf("%s%s TTL expired %d seconds ago, deleting",
				c.dryRunStr, fqName, expiredSecs)
			if !c.DryRun {
				if err := Delete(); err != nil {
					log.Errorf("failed to delete %s with error: %v", fqName, err)
				} else {
					cnt++
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
			if err := Update(); err != nil {
				log.Errorf("failed to update %s with error: %v", fqName, err)
			} else {
				cnt++
			}
		}
	}
	return cnt
}
