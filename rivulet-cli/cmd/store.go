package cmd

import (
	"io"
	"os"
	"time"

	"github.com/kurocifer/rivulet/utils"
	"github.com/spf13/cobra"
)

var addr string
var nodes []string
var key string

// storeCmd represents the store command
var storeCmd = &cobra.Command{
	Use:   "store",
	Short: "Stores a file and replicates to other nodes on the network",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, err := os.Open(args[0])
		if err != nil {
			return err
		}

		return storeCmdEx(file)
	},
}

func storeCmdEx(r io.Reader) error {
	server := utils.MakeServer(addr, nodes...)
	go server.Start()
	time.Sleep(2 * time.Second)

	err := server.Store(key, r)
	if err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.AddCommand(storeCmd)

	storeCmd.Flags().StringVarP(&addr, "port", "p", ":8080", "Port on which server should run")
	storeCmd.Flags().StringSliceVarP(&nodes, "nodes", "n", []string{}, "Address of nodes to connect to")
	storeCmd.Flags().StringVarP(&key, "key", "k", "", "File key (will be used to identify the file)")
}
