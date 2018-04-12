package cmd

import (
	"github.com/jjo/kube-custodian/pkg/cleaner"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagRequiredLabels = "required-labels"
)

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
		cleaner.DeleteJobs(client, dryRun, namespace, requiredLabels)
		cleaner.DeleteDeployments(client, dryRun, namespace, requiredLabels)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.PersistentFlags().StringSlice(flagRequiredLabels, []string{"created_by"}, "Labels required for resources to be skipped from scanning")
}
