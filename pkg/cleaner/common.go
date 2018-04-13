package cleaner

import (
	"regexp"

	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	utils "github.com/jjo/kube-custodian/pkg/utils"
)

type skipMetaType struct {
	NamespaceRE     string
	Labels          []string
	NamespaceRegexp *regexp.Regexp
}

var SkipMetaDefault = &skipMetaType{
	NamespaceRE: "kube-.*|.*(-system|monitoring|logging|ingress)",
	Labels:      []string{"created_by"},
}

var skipMeta skipMetaType

func init() {
	SetSkipMeta("", nil)
}

func skipFromMeta(meta *metav1.ObjectMeta) bool {
	skipIt := false
	switch {
	case skipMeta.NamespaceRegexp.MatchString(meta.Namespace):
		log.Debugf("%s.%s skipped from meta.Namespace", meta.Name, meta.Namespace)
		skipIt = true
	case utils.LabelsSubSet(meta.Labels, skipMeta.Labels):
		log.Debugf("%s.%s skipped from meta.Labels", meta.Name, meta.Labels)
		skipIt = true
	}
	return skipIt
}

// SetSkipMeta uses defaults if called as ("", nil)
func SetSkipMeta(namespaceRe string, labels []string) {
	skipMeta = *SkipMetaDefault
	if namespaceRe != "" {
		skipMeta.NamespaceRE = namespaceRe
	}
	if labels != nil {
		skipMeta.Labels = labels
	}
	skipMeta.NamespaceRegexp = regexp.MustCompile(skipMeta.NamespaceRE)
}
