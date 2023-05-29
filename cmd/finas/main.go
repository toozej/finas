package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/toozej/finas/internal/dcmd"
	"github.com/toozej/finas/pkg/man"
	"github.com/toozej/finas/pkg/version"
)

func main() {
	// load Docker Run cmds
	if err := dcmd.LoadCmds(); err != nil {
		panic(fmt.Errorf("invalid Docker Run commands: %s", err))
	}

	command := &cobra.Command{
		Use:   "finas",
		Short: "FINAS Is Not A Shell",
		Long:  `FINAS is a canned Docker Run command tool for easy command recollection and reuse`,
		Run: func(cmd *cobra.Command, args []string) {
			for filePath, command := range dcmd.CmdMap {
				fmt.Println("File Path:", filePath)
				fmt.Println("Image Name:", command.Name)
				fmt.Println("Tag:", command.Tag)
				fmt.Println("Flags:", command.Flags)
				fmt.Println("Volumes:", command.Volumes)
				fmt.Println("Bind Mounts:", command.BindMounts)
				fmt.Println("Networks:", command.Networks)
				fmt.Println("Published Ports:", command.PublishedPorts)
				fmt.Println()
			}
		},
	}

	command.AddCommand(
		man.NewManCmd(),
		version.Command(),
	)

	if err := command.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

}
