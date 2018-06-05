package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time
	Version string
	// Revision also
	Revision string
)

var versionCmd = &cobra.Command{
	Use: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kube-custodian version %s, %s\n", Version, Revision)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
