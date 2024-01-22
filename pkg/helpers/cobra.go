package helpers

import (
	"bytes"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func ExecuteCommandC(command *cobra.Command, args []string) (*cobra.Command, string, error) {
	// SetOut and SetErr https://github.com/spf13/cobra/issues/1484
	buf := new(bytes.Buffer)
	command.SetOut(buf)
	command.SetErr(buf)
	command.SetArgs(args)

	command, err := command.ExecuteC()
	if err != nil {
		log.Fatalf("There was an error executing the command: %v", err)
	}

	read := buf.String()
	return command, read, err

}
