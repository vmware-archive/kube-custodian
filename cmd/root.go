package cmd

import (
	"bytes"
	goflag "flag"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	flagDryRun        = "dry-run"
	flagAllNamespaces = "all-namespaces"
	flagVerbose       = "verbose"
)

var clientConfig clientcmd.ClientConfig
var overrides clientcmd.ConfigOverrides

// Gently taken from github.com/ksonnet/kubecfg
func init() {
	rootCmd.PersistentFlags().Bool(flagDryRun, false, "Dry run")
	rootCmd.PersistentFlags().CountP(flagVerbose, "v", "Increase verbosity")
	rootCmd.PersistentFlags().Bool(flagAllNamespaces, false, "All namespaces")

	// The "usual" clientcmd/kubectl flags
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	kflags := clientcmd.RecommendedConfigOverrideFlags("")
	rootCmd.PersistentFlags().StringVar(&loadingRules.ExplicitPath, "kubeconfig", "", "Path to a kube config. Only required if out-of-cluster")
	rootCmd.MarkPersistentFlagFilename("kubeconfig")
	clientcmd.BindOverrideFlags(&overrides, rootCmd.PersistentFlags(), kflags)
	clientConfig = clientcmd.NewInteractiveDeferredLoadingClientConfig(loadingRules, &overrides, os.Stdin)

	rootCmd.PersistentFlags().Set("logtostderr", "true")
}

var rootCmd = &cobra.Command{
	Use:           "kube-custodian",
	Short:         "Cleanup stuff",
	SilenceErrors: true,
	SilenceUsage:  true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		goflag.CommandLine.Parse([]string{})
		flags := cmd.Flags()

		// Need either --all-namespaces or --namespace=...
		allNS, _ := flags.GetBool(flagAllNamespaces)
		namespace, _ := flags.GetString("namespace")
		if allNS == (len(namespace) > 0) {
			log.Fatalf("Cowardly refusing to use a default namespace, provide --namespace=<NS> xor --all-namespaces")
		}

		out := cmd.OutOrStderr()
		log.SetOutput(out)

		logFmt := NewLogFormatter(out)
		log.SetFormatter(logFmt)

		verbosity, err := flags.GetCount(flagVerbose)
		if err != nil {
			return err
		}
		log.SetLevel(logLevel(verbosity))

		return nil
	},
}

// Execute is main() entry point
func Execute() error {
	if err := rootCmd.Execute(); err != nil {
		log.Errorf("rootCmd.Execute(): %v", err)
		return err
	}
	return nil
}

// NewKubeClient creates a new kubernetes Clientset from already setup clientConfig
func NewKubeClient() *kubernetes.Clientset {
	c, err := clientConfig.ClientConfig()

	if err != nil {
		log.Fatalf("NewKubeClient: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(c)
	if err != nil {
		log.Fatalf("NewKubeClient: failed to create clienset: %v", err)
	}

	return clientset
}

type logFormatter struct {
	escapes  *terminal.EscapeCodes
	colorise bool
}

func logLevel(verbosity int) log.Level {
	switch verbosity {
	case 0:
		return log.InfoLevel
	default:
		return log.DebugLevel
	}
}

// NewLogFormatter creates a new log.Formatter customised for writer
func NewLogFormatter(out io.Writer) log.Formatter {
	var ret = logFormatter{}
	if f, ok := out.(*os.File); ok {
		ret.colorise = terminal.IsTerminal(int(f.Fd()))
		ret.escapes = terminal.NewTerminal(f, "").Escape
	}
	return &ret
}

func (f *logFormatter) Format(e *log.Entry) ([]byte, error) {
	buf := bytes.Buffer{}
	if f.colorise {
		buf.Write(f.levelEsc(e.Level))
		fmt.Fprintf(&buf, "%-5s ", strings.ToUpper(e.Level.String()))
		buf.Write(f.escapes.Reset)
	}

	buf.WriteString(strings.TrimSpace(e.Message))
	buf.WriteString("\n")

	return buf.Bytes(), nil
}
func (f *logFormatter) levelEsc(level log.Level) []byte {
	switch level {
	case log.DebugLevel:
		return []byte{}
	case log.WarnLevel:
		return f.escapes.Yellow
	case log.ErrorLevel, log.FatalLevel, log.PanicLevel:
		return f.escapes.Red
	default:
		return f.escapes.Blue
	}
}
