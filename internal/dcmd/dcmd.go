package dcmd

import (
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

type DockerCmd struct {
	Name           string   `json:"name,omitempty" mapstructure:"name"`
	Image          string   `json:"image" mapstructure:"image"`
	Tag            string   `json:"tag" mapstructure:"tag"`
	Flags          []string `json:"flags,omitempty" mapstructure:"flags"`
	Volumes        []string `json:"volumes,omitempty" mapstructure:"volumes"`
	BindMounts     []string `json:"bind_mounts,omitempty" mapstructure:"bind_mounts"`
	Networks       []string `json:"networks,omitempty" mapstructure:"networks"`
	PublishedPorts []string `json:"published_ports,omitempty" mapstructure:"published_ports"`
	Entrypoint     string   `json:"entrypoint,omitempty" mapstructure:"entrypoint"`
	Arguments      []string `json:"arguments,omitempty" mapstructure:"arguments"`
	UserArguments  []string `json:"omitempty" mapstructure:"user_arguments"`
	Help           string   `json:"help" mapstructure:"help"`
}

// Map to store the loaded DockerCmd structs
var CmdMap map[string]DockerCmd

func cmdPaths(additionalPaths []string) ([]string, error) {
	// setup list of paths where Docker Run cmd .json files could exist
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Error getting user's home directory:", err)
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
func LoadCmds(additionalPaths ...string) (map[string]DockerCmd, error) {
	v := viper.GetViper()
	CmdMap := make(map[string]DockerCmd)
	if v.GetBool("debug") {
		log.Print("From LoadCmds: Debug is enabled!")
		log.SetLevel(log.DebugLevel)
	}

	directories, _ := cmdPaths(additionalPaths)

	// Iterate over input directories
	for _, dir := range directories {
		files, err := os.ReadDir(dir)
		if err != nil {
			if os.IsNotExist(err) {
				log.Debugf("Directory not found, continuing: %s\n", dir)
				continue
			} else {
				log.Errorf("Error reading directory '%s': %s\n", dir, err)
				continue
			}
		}

		// Iterate over files in the current directory
		for _, file := range files {
			if strings.HasSuffix(file.Name(), ".json") {
				filePath := filepath.Join(dir, file.Name())
				fileName := path.Base(filePath)
				fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]

				// Set Viper to read the JSON file
				v.SetConfigType("json")
				v.SetConfigFile(filePath)

				// Read the JSON file into the DockerImage struct
				var cmd DockerCmd
				err := v.ReadInConfig()
				if err != nil {
					log.Errorf("Error reading JSON file '%s': %s\n", filePath, err)
					continue
				}

				err = v.Unmarshal(&cmd)
				if err != nil {
					log.Errorf("Error parsing JSON file '%s': %s\n", filePath, err)
					continue
				}

				// set the dcmd's Name field to the JSON filename without extension
				cmd.Name = fileNameWithoutExt

				// Store the loaded struct in the map at key based on filename
				CmdMap[fileNameWithoutExt] = cmd
			}
		}
	}

	if v.GetBool("debug") {
		printCmdMap(CmdMap)
	}

	return CmdMap, nil
}

func printCmdMap(cmdMap map[string]DockerCmd) {
	log.Debugf("CmdMap currently contains: %v\n", cmdMap)
	for key, cmd := range cmdMap {
		// Get the type of the struct
		t := reflect.TypeOf(cmd)

		// Iterate over the fields
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			value := reflect.ValueOf(cmd).Field(i).Interface()

			log.Debugf("Key: %s, Field: %s, Value: %v\n", key, field.Name, value)
		}
	}
}
