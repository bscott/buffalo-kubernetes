package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// herokuCmd represents the heroku command
var herokuCmd = &cobra.Command{
	Use:     "kubernetes",
	Aliases: []string{"h"},
	Short:   "Tools for deploying Buffalo to kubernetes",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("kubernetes called")
	},
}

func init() {
	RootCmd.AddCommand(kubernetesCmd)
}
