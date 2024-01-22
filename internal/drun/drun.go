package drun

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/toozej/finas/internal/dcmd"
)

// gather Docker command arguments and return them as a properly formatted `docker run ...` args slice
func buildRunArgs(dcmd dcmd.DockerCmd) []string {
	runArgs := []string{}

	runArgs = append(runArgs, fmt.Sprintf("--name %s", dcmd.Name))
	if dcmd.Flags != nil {
		runArgs = append(runArgs, dcmd.Flags...)
	}

	if dcmd.Volumes != nil {
		volumeArgs := []string{}
		for _, v := range dcmd.Volumes {
			volumeArgs = append(volumeArgs, fmt.Sprintf("--volume %s", v))
		}
		runArgs = append(runArgs, volumeArgs...)
	}

	if dcmd.BindMounts != nil {
		bindArgs := []string{}
		for _, v := range dcmd.BindMounts {
			bindArgs = append(bindArgs, fmt.Sprintf("--volume %s", v))
		}
		runArgs = append(runArgs, bindArgs...)
	}

	if dcmd.Networks != nil {
		networkArgs := []string{}
		for _, v := range dcmd.Networks {
			networkArgs = append(networkArgs, fmt.Sprintf("--network %s", v))
		}
		runArgs = append(runArgs, networkArgs...)
	}

	if dcmd.PublishedPorts != nil {
		portArgs := []string{}
		for _, v := range dcmd.PublishedPorts {
			portArgs = append(portArgs, fmt.Sprintf("--publish %s", v))
		}
		runArgs = append(runArgs, portArgs...)
	}

	if dcmd.Entrypoint != "" {
		runArgs = append(runArgs, fmt.Sprintf("--entrypoint %s", dcmd.Entrypoint))
	}

	runArgs = append(runArgs, dcmd.Image)

	if dcmd.Arguments != nil {
		runArgs = append(runArgs, dcmd.Arguments...)
	}

	return runArgs
}

// run the `docker run` command with args gathered from Docker command struct
// and user-supplied args as well
func RunDockerCommand(dcmd dcmd.DockerCmd, userArgs []string) {
	// setup debug logging if enabled
	v := viper.GetViper()
	if v.GetBool("debug") {
		log.Print("From RunDockerCommand: Debug is enabled!")
		log.SetLevel(log.DebugLevel)
	}

	// set up `docker run ...` args based on DockerCmd contents
	runArgs := buildRunArgs(dcmd)

	if v.GetBool("debug") {
		log.Debug("Args contains:")
		for _, v := range runArgs {
			log.Debugf("%s is type %T\n", v, v)
		}
		for _, v := range userArgs {
			log.Debugf("%s is type %T\n", v, v)
		}
	}

	// space-separate Args items to be fed into exec.Command()
	joinedRunArgs := strings.Join(append(append([]string{"docker", "run"}, runArgs...), userArgs...), " ")

	// setup actual `docker run ...` command
	dockerCmd := exec.Command("/bin/sh", "-c", joinedRunArgs) //#nosec nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	dockerCmd.Stdout = os.Stdout
	dockerCmd.Stderr = os.Stderr
	log.Debugf("dockerCmd is: %s", dockerCmd)

	// run `docker run ...` command and check for errors
	err := dockerCmd.Run()
	if err != nil {
		log.Fatalf("Error running Docker command: %v", err)
	}

	log.Infof("Successfully ran Docker command for container: %s\n", dcmd.Name)
}
