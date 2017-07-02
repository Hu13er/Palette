package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitlab.com/NagByte/Palette/common"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:    ", common.Version)
		fmt.Println("Build Time: ", common.BuildTime)
		fmt.Println("Commit:     ", common.Commit)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
