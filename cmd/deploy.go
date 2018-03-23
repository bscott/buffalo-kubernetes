package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var deployCmd = &cobra.Command{
	Use:     "deploy",
	Aliases: []string{"d"},
	Short:   "deploy to heroku using docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		return pushContainer()
	},
}

func pushContainer() error {
	c := exec.Command("kubectl", "run", "buffpush")
	fmt.Println(strings.Join(c.Args, " "))
	c.Stdin = os.Stdin
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	err := c.Run()
	if err != nil {
		return errors.WithStack(err)
	}
	return runMigrations()
}

func init() {
	herokuCmd.AddCommand(deployCmd)
}
