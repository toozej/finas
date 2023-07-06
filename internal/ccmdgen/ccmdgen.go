package ccmdgen

import (
	"github.com/spf13/cobra"
	"github.com/toozej/finas/internal/dcmd"
	"github.com/toozej/finas/internal/drun"
)

func NewFinasCommand(dcmd dcmd.DockerCmd) *cobra.Command {
	cmd := &cobra.Command{
		Use:   dcmd.Name,
		Short: dcmd.Help,
		Run: func(c *cobra.Command, args []string) {
			drun.RunDockerCommand(dcmd, args)
		},
	}

	return cmd
}
