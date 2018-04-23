package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jjo/kube-custodian/pkg/cleaner"
)

const (
	flagSkipLabels      = "skip-labels"
	flagSkipNamespaceRe = "skip-namespace-re"
	flagTagTTL          = "tag-ttl"
	flagTagForDeletion  = "tag-for-deletion"
	flagDeleteTagged    = "delete-tagged"
	flagCleanUpTag      = "cleanup-tag"
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringSlice(flagSkipLabels, cleaner.CommonDefaults.SkipLabels, "Labels required for resources to be skipped from scanning")
	runCmd.PersistentFlags().String(flagSkipNamespaceRe, cleaner.CommonDefaults.SkipNamespaceRE, "Regex of namespaces to skip, typically 'system' ones and alike")
	runCmd.PersistentFlags().String(flagTagTTL, cleaner.CommonDefaults.TagTTL, "Time to live after marked, before deletion")
	runCmd.PersistentFlags().Bool(flagTagForDeletion, true, "Tag resources for later deletion")
	runCmd.PersistentFlags().Bool(flagDeleteTagged, true, "Delete tagged resources, after their Tag TTL has passed")
	runCmd.PersistentFlags().Bool(flagCleanUpTag, false, "Untag resources from later deletion")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Scan Kubernetes objects, mark for deletion (via annotation), delete those already \"expired\"",
	Run: func(cmd *cobra.Command, args []string) {
		var err error
		flags := cmd.Flags()
		c := cleaner.Common{}

		c.SkipLabels, err = flags.GetStringSlice(flagSkipLabels)
		if err != nil {
			log.Fatal(err)
		}
		c.Namespace, err = flags.GetString("namespace")
		if err != nil {
			log.Fatal(err)
		}
		c.TagTTL, err = flags.GetString(flagTagTTL)
		if err != nil {
			log.Fatal(err)
		}
		c.DryRun, err = flags.GetBool(flagDryRun)
		if err != nil {
			log.Fatal(err)
		}

		if len(c.SkipLabels) < 1 {
			log.Fatal("At least one skip-labels is needed")
		}
		log.Debugf("Skipping workloads with labels: %v ...", c.SkipLabels)

		c.SkipNamespaceRE, err = flags.GetString(flagSkipNamespaceRe)
		if err != nil {
			log.Fatal(err)
		}
		c.Init(NewKubeClient(cmd))
		_, _, errCnt := c.Run()
		os.Exit(errCnt)
	},
}
