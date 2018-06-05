package main

import (
	"os"

	"github.com/bitnami-labs/kube-custodian/cmd"
	log "github.com/sirupsen/logrus"
)

var (
	version  string
	revision string
)

func main() {

	cmd.Version = version
	cmd.Revision = revision
	if err := cmd.Execute(); err != nil {
		// PersistentPreRunE may not have been run for early
		// errors, like invalid command line flags.
		logFmt := cmd.NewLogFormatter(log.StandardLogger().Out)
		log.SetFormatter(logFmt)
		log.Error(err.Error())

		switch err {
		case nil:
			os.Exit(0)
		default:
			os.Exit(1)
		}
	}
}
