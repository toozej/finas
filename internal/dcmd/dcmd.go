package dcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

type DockerCmd struct {
	Name           string   `json:"name"`
	Tag            string   `json:"tag"`
	Flags          []string `json:"flags,omitempty"`
	Volumes        []string `json:"volumes,omitempty"`
	BindMounts     []string `json:"bind_mounts,omitempty"`
	Networks       []string `json:"networks,omitempty"`
	PublishedPorts []string `json:"published_ports,omitempty"`
	Arguments      []string `json:"arguments,omitempty"`
}

// Map to store the loaded DockerCmd structs
var CmdMap map[string]DockerCmd

func cmdPaths(additionalPaths []string) ([]string, error) {
	// setup list of paths where Docker Run cmd .json files could exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error:", err)
		return []string{}, err
	}

	paths := []string{
		"./config",
		filepath.Join(homeDir, ".config/finas"),
		"/usr/local/etc/finas",
		"/etc/finas",
	}

	if len(additionalPaths) > 0 {
		paths = append(paths, additionalPaths...)
	}

	return paths, nil
}

// LoadCmds loads Docker Run cmds from JSON files
// located in cmdPaths in order of preference
func LoadCmds(additionalPaths ...string) error {
	v := viper.New()

	directories, _ := cmdPaths(additionalPaths)

	// Iterate over input directories
	for _, dir := range directories {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			fmt.Printf("Error reading directory '%s': %s\n", dir, err)
			continue
		}

		// Iterate over files in the current directory
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".json") {
				filePath := filepath.Join(dir, file.Name())

				// Set Viper to read the JSON file
				v.SetConfigFile(filePath)

				// Read the JSON file into the DockerImage struct
				var cmd DockerCmd
				err := v.ReadInConfig()
				if err != nil {
					fmt.Printf("Error reading JSON file '%s': %s\n", filePath, err)
					continue
				}

				err = v.Unmarshal(&cmd)
				if err != nil {
					fmt.Printf("Error parsing JSON file '%s': %s\n", filePath, err)
					continue
				}

				// Store the loaded struct in the map
				CmdMap[filePath] = cmd
			}
		}
	}

	return v.Unmarshal(&CmdMap)
}
