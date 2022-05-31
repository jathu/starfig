package main

import (
	"fmt"
	"os"

	"github.com/jathu/starfig/internal/command"
	"github.com/jathu/starfig/internal/logging"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func safeExit(err error) {
	if err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}

// Populated via go build with flag: -ldflags "-X main.starfigVersion=<version>"
// This is the default.
var starfigVersion string = "dev"

func main() {
	logging.SetupLogger()
	logrus.Debug("starting starfig")

	var rootCmd = &cobra.Command{
		Use:   "starfig",
		Short: "Starfig is a programmatic config generator. It helps create static configs using Starlark, a deterministic Python-like language. Learn more at https://github.com/jathu/starfig.",
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		Run: func(cmd *cobra.Command, args []string) {
			safeExit(cmd.Help())
		},
	}

	var buildKeepGoing bool
	buildCmd := cobra.Command{
		Use:   "build [targets...]",
		Short: "Build config targets.",
		Long:  `Build config targets within the universe. The argument takes a list of build targets. The argument also allows building a whole package by using the spread operator. i.e. //... //example/...`,
		Run: func(cmd *cobra.Command, args []string) {
			safeExit(command.Build(args, buildKeepGoing))
		},
	}
	buildCmd.Flags().BoolVar(&buildKeepGoing, "keep-going", false, "Continue to build as many targets as possible even if there are errors.")
	rootCmd.AddCommand(&buildCmd)

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the starfig version.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(fmt.Sprintf("starfig-%s", starfigVersion))
		},
	})

	safeExit(rootCmd.Execute())
}
