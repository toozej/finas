package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/toozej/finas/internal/ccmdgen"
	"github.com/toozej/finas/internal/dcmd"
	"github.com/toozej/finas/pkg/man"
	"github.com/toozej/finas/pkg/version"
)

var v *viper.Viper

var rootCmd = &cobra.Command{
	Use:     "finas",
	Aliases: []string{"f"},
	Short:   "FINAS Is Not A Shell",
	Long:    `FINAS is a canned Docker Run command tool for easy command recollection and reuse`,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("To see available commands, run 'finas help'")
	},
}

func main() {
	v = viper.GetViper()
	// setup debug flag
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug-level logging")
	if err := v.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Error("Error setting debug root-level flag")
	}

	// setup man and version sub-commands
	rootCmd.AddCommand(
		man.NewManCmd(),
		version.Command(),
	)

	// load Docker Run cmds
	cmdMap, err := dcmd.LoadCmds()
	if err != nil {
		log.Fatalf("Loading Docker Run commands failed: %s\n", err)
	}
	log.Debugf("From rootCmd: cmdMap contains %v", cmdMap)
	// Access and use the DockerCmd objects in the map
	for key, cmd := range cmdMap {
		log.Debug("Key:", key)
		log.Debug("Name:", cmd.Name)
		log.Debug("Image:", cmd.Image)
		log.Debug("Tag:", cmd.Tag)
		log.Debug("Flags:", cmd.Flags)
		log.Debug("Volumes:", cmd.Volumes)
		log.Debug("BindMounts:", cmd.BindMounts)
		log.Debug("Networks:", cmd.Networks)
		log.Debug("PublishedPorts:", cmd.PublishedPorts)
		log.Debug("Entrypoint:", cmd.Entrypoint)
		log.Debug("Arguments:", cmd.Arguments)
		log.Debug("Help:", cmd.Help)

		// generate a finas sub-command
		rootCmd.AddCommand(ccmdgen.NewFinasCommand(cmd))
	}

	// run main root command
	if err := rootCmd.Execute(); err != nil {
		log.Fatal("Error running rootCmd: ", err.Error())
	}
}
