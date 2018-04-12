package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/jjo/kube-custodian/pkg/cleaner"
)

const (
	flagRequiredLabels = "required-labels"
	flagSystemNS       = "sys-namespaces-re"
)

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringSlice(flagRequiredLabels, []string{"created_by"}, "Labels required for resources to be skipped from scanning")
	deleteCmd.PersistentFlags().String(flagSystemNS, cleaner.SystemNS, "\"system\" namespaces to skip")
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete resources",
	Run: func(cmd *cobra.Command, args []string) {
		flags := cmd.Flags()
		requiredLabels, err := flags.GetStringSlice(flagRequiredLabels)
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

		if len(requiredLabels) < 1 {
			log.Fatal("At least one required-label is needed")
		}
		log.Infof("Required labels: %v ...", requiredLabels)

		sysNS, err := flags.GetString(flagSystemNS)
		if err != nil {
			log.Fatal(err)
		}

		cleaner.SetSystemNS(sysNS)
		cleaner.DeleteJobs(client, dryRun, namespace, requiredLabels)
		cleaner.DeleteDeployments(client, dryRun, namespace, requiredLabels)
	},
}
