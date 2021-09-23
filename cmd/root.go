package cmd

import (
	"csi-nfs/pkg/nfs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"time"
)

var (
	debug    bool
	name     string
	nodeID   string
	endpoint string
	version  string

	// prepare  some  parameters  to init GRPC  server
	// or prepare for other store type like lvm and cloud or disk
	parameter1 string
	parameter2 int
	parameter3 time.Duration
)

var rootCmd = &cobra.Command{
	Use:        "csi-nfs",
	Aliases:    nil,
	SuggestFor: nil,
	Short:      "CSI PROJECT ARCHE",
	Long:       "",
	Example:    "",
	// NewCSIDriver
	Run: func(cmd *cobra.Command, args []string) {
		nfs.NewCSIDriver(name, version, nodeID, endpoint, parameter1, parameter2, parameter3).Run()

		logrus.Info("start the NewCSIDriver...")
	},
}

//  init the log formatter
func initLog() {
	if debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

}

// init  parameter by env
func init() {
	// Initialize the log component after parsing command line parameters
	cobra.OnInitialize(initLog)

	// The debug option is used to control whether logrus outputs debug logs
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug log")

	rootCmd.PersistentFlags().StringVar(&version, "1.0.1", "", "csi version ")

	// The parameters necessary for running the CSI plug-in must be set by the user
	rootCmd.PersistentFlags().StringVar(&nodeID, "nodeid", "", "csi node id")
	_ = rootCmd.MarkPersistentFlagRequired("nodeid")
	rootCmd.PersistentFlags().StringVar(&endpoint, "endpoint", "", "csi endpoint")
	_ = rootCmd.MarkPersistentFlagRequired("endpoint")

	// Users generally do not need to modify the csi plugin name, so we hide the `--name` option
	rootCmd.PersistentFlags().StringVar(&name, "name", "csi-archetype", "csi name")
	_ = rootCmd.PersistentFlags().MarkHidden("name")
	_ = rootCmd.MarkPersistentFlagRequired("name")

	// The user needs to add some parameters according to actual needs to ensure
	// the successful initialization of the CSI plug-in
	rootCmd.PersistentFlags().StringVar(&parameter1, "parameter1", "", "csi parameter1")
	rootCmd.PersistentFlags().IntVar(&parameter2, "parameter2", 10, "csi parameter2")
	rootCmd.PersistentFlags().DurationVar(&parameter3, "parameter3", 10*time.Second, "csi parameter3")

}

func Excute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Failed to start %v", err)
	}

}
