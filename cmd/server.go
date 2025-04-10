package cmd

import (
	"github.com/maomaozgw/filament/cmd/serve"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "filamentinventory",
	Short: "filamentinventory is a web server for 3d printer filament management.",
	Long:  `filamentinventory is a web server for 3d printer filament management.`,
}

func init() {
	rootCmd.AddCommand(serve.Cmd)
}

func Execute() {
	rootCmd.Execute()
}
