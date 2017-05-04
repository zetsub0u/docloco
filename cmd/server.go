package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zetsub0u/docloco/docloco"
)

var cmdServer = &cobra.Command{
	Use:   "server",
	Short: "Run the docloco server",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		runServer()
	},
}

func init() {
	RootCmd.AddCommand(cmdServer)
}

func runServer() {
	docloco.RunServer()
}
