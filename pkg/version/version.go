package version

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Version information. These will be filled in by the compiler.
var (
	Version = "local"
	Commit  = ""
	Branch  = ""
	BuiltAt = ""
	Builder = ""
)

// Info holds build information
type Info struct {
	Commit  string
	Version string
	Branch  string
	BuiltAt string
	Builder string
}

// Get creates an initialized Info object
func Get() (Info, error) {
	return Info{
		Commit:  Commit,
		Version: Version,
		Branch:  Branch,
		BuiltAt: BuiltAt,
		Builder: Builder,
	}, nil
}

// Command creates version command
func Command() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version.",
		Long:  `Print the version and build information.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			info, err := Get()
			if err != nil {
				return err
			}

			fmt.Println("Version: ", info.Version)
			fmt.Println("Git commit: ", info.Commit)
			fmt.Println("Built At: ", info.BuiltAt)
			if viper.GetBool("debug") {
				fmt.Println("Git Branch: ", info.Branch)
				fmt.Println("Builder: ", info.Builder)
			}

			return nil
		},
	}
}
