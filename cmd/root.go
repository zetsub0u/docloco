package cmd

import (
	"github.com/spf13/cobra"
	"github.com/zetsub0u/docloco/config"
	"os"
)

// CLI Entrypoint
var RootCmd = &cobra.Command{
	Use:   "docloco",
	Short: "",
	Long:  ``,
}

func init() {
	var cfgFile string
	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "docloco", "Use specific config file.")
	RootCmd.ParseFlags(os.Args)

	// Load Configuration
	config.Store.Load(cfgFile)
}
