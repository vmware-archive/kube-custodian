package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jjo/kube-custodian/pkg/cleaner"
)

const (
	flagSkipLabels      = "skip-labels"
	flagSkipNamespaceRe = "skip-namespace-re"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringSlice(flagSkipLabels, cleaner.SkipMetaDefault.Labels, "Labels required for resources to be skipped from scanning")
	deleteCmd.PersistentFlags().String(flagSkipNamespaceRe, cleaner.SkipMetaDefault.NamespaceRE, "Regex of namespaces to skip, typically 'system' ones and alike")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		skipLabels, err := flags.GetStringSlice(flagSkipLabels)
		if err != nil {
			log.Fatal(err)
		}
		namespace, err := flags.GetString("namespace")
		if err != nil {
			log.Fatal(err)
		}
		dryRun, err := flags.GetBool(flagDryRun)
		if err != nil {
			log.Fatal(err)
		}
		client := NewKubeClient(cmd)

		if len(skipLabels) < 1 {
			log.Fatal("At least one skip-labels is needed")
		}
		log.Infof("Required labels: %v ...", skipLabels)

		skipNSRe, err := flags.GetString(flagSkipNamespaceRe)
		if err != nil {
			log.Fatal(err)
		}

		cleaner.SetSkipMeta(skipNSRe, skipLabels)
		cleaner.DeleteDeployments(client, dryRun, namespace)
		cleaner.DeleteStatefulSets(client, dryRun, namespace)
		cleaner.DeleteJobs(client, dryRun, namespace)
		cleaner.DeletePods(client, dryRun, namespace)
	},
}
