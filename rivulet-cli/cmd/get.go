package cmd

import (
	"time"

	"github.com/kurocifer/rivulet/utils"
	"github.com/spf13/cobra"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Gets the file specified by key from local storage, if it's not found, it fetches it from nodes",
	RunE: func(cmd *cobra.Command, args []string) error {
		server := utils.MakeServer(addr, nodes...)
		go server.Start()
		time.Sleep(2 * time.Second)
		_, err := server.Get(key)
		return err
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// I know I know I'm repeating this and all...
	getCmd.Flags().StringVarP(&addr, "port", "p", ":8080", "Port on which server should run")
	getCmd.Flags().StringSliceVarP(&nodes, "nodes", "n", []string{}, "Address of nodes to connect to")
	getCmd.Flags().StringVarP(&key, "key", "k", "", "File key (will be used to identify the file)")
}
