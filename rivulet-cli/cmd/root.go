package cmd

import (
	"fmt"
	"os"

	"github.com/kurocifer/rivulet/rivulet-cli/pkg/daemon"
	"github.com/kurocifer/rivulet/rivulet-cli/pkg/utils"
	"github.com/spf13/cobra"
)

var action string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "rvt-cli",
	Short: "A decentralized and distributed file storage system",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 && !cmd.Flags().Changed("action") {
			action = "status"
			defer cmd.HelpFunc()(cmd, args)
		} else {
			println(action)
		}

		utils.CreateWorkDir()

		err := rootCmdEx()
		if err != nil {
			fmt.Println(err)
			cmd.HelpFunc()(cmd, args)
		}
	},
}

func rootCmdEx() error {
	switch action {
	case "start":
		daemon.StartDaemon()

	case "stop":
		daemon.StopDaemon()

	case "status":
		daemon.GetDaemonStatus()

	default:
		return fmt.Errorf("unknown action '%s'", action)
	}

	return nil
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&action, "action", "a", "", "Specify what you want to do with the daemon: 'start', 'stop', 'status'")
}
