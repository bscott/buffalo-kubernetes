package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/fatih/structs"
	"github.com/gobuffalo/makr"
	"github.com/markbates/going/randx"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/AlecAivazis/survey.v1"
)

// availableCmd represents the available command
var setupCmd = &cobra.Command{
	Use:     "setup",
	Aliases: []string{"s"},
	Short:   "setup kubernetes for deploying with docker",
	RunE: func(cmd *cobra.Command, args []string) error {
		if setup.Interactive {
			err := Interactive()
			if err != nil {
				return errors.WithStack(err)
			}
		}
		return setup.Run()
	},
}

func Interactive() error {
	qs := []*survey.Question{
		{
			Name: "AppName",
			Prompt: &survey.Input{
				Message: "What would you like to name this app in Kubernetes?",
				Help:    "A blank response will default to 'BuffaloApp'.",
			},
		},
	}
	err := survey.Ask(qs, &setup)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

var nameSpaces = []string{"default", "buffalo"}
var setup = Setup{}

func init() {
	setupCmd.Flags().StringVarP(&setup.AppName, "app-name", "a", "", "the name for the kubernetes namespace")
	setupCmd.Flags().BoolVarP(&setup.Interactive, "interactive", "i", false, "use the interactive mode")
	herokuCmd.AddCommand(setupCmd)
}

type Setup struct {
	AppName       string
	Interactive   bool
}

func (s Setup) Run() error {
	g := makr.New()
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return validateGit()
		},
	})
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return installKubeCLI()
		},
	})
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return installHelm()
		},
	})
	
	g.Add(makr.Func{
		Runner: func(root string, data makr.Data) error {
			return initializeHostVar()
		},
	})
	if s.Database != "" {
		// Add DB Helm here
	}
	return g.Run(".", structs.Map(s))
	}

	return nil
}

func installKubeCLI() error {
	if _, err := exec.LookPath("kubectl"); err != nil {
		if runtime.GOOS == "darwin" {
			if _, err := exec.LookPath("brew"); err == nil {
				c := exec.Command("brew", "install", "kubectl")
				c.Stdin = os.Stdin
				c.Stderr = os.Stderr
				c.Stdout = os.Stdout
				return c.Run()
			}
		}
		return errors.New("Kubectl is not installed. https://kubernetes.io/docs/tasks/tools/install-kubectl/")
	}
	fmt.Println("--> Kubectl is installed")
	return nil
}

func installHelm() error {
	if _, err := exec.LookPath("helm"); err != nil {
		if runtime.GOOS == "darwin" {
			if _, err := exec.LookPath("brew"); err == nil {
				c := exec.Command("brew", "install", "kubernetes-helm")
				c.Stdin = os.Stdin
				c.Stderr = os.Stderr
				c.Stdout = os.Stdout
				return c.Run()
			}
		}
		return errors.New("Kubernetes Helm is not installed. https://github.com/kubernetes/helm/blob/master/docs/install.md")
	}
	fmt.Println("--> Kubernetes Helm is installed")
	return nil
}