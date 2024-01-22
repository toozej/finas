package drun

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/toozej/finas/internal/dcmd"
	"github.com/toozej/finas/pkg/helpers"
)

func TestBuildRunArgs(t *testing.T) {
	expectedRunArgs := []string{
		"--name helloworld",
		"--rm",
		"hello-world",
	}

	dcmd := dcmd.DockerCmd{
		Name:           "helloworld",
		Image:          "hello-world",
		Tag:            "latest",
		Flags:          []string{"--rm"},
		Volumes:        []string{},
		BindMounts:     []string{},
		Networks:       []string{},
		PublishedPorts: []string{},
		Entrypoint:     "",
		Arguments:      []string{},
		Help:           "hello world example",
	}
	actualRunArgs := buildRunArgs(dcmd)

	if !cmp.Equal(actualRunArgs, expectedRunArgs, cmpopts.IgnoreUnexported()) {
		t.Errorf("Actual runArgs does not match expected. Got %v, expected %v", actualRunArgs, expectedRunArgs)
	}
}

func TestRunDockerCommand(t *testing.T) {
	t.Skip("TODO TestRunDockerCommand needs to be fixed")

	// generate "unique" name for helloworld container
	minLength := 5
	maxLength := 15
	max := big.NewInt(int64((maxLength - minLength + 1) + minLength))
	randomLength, _ := rand.Int(rand.Reader, max)
	randomLengthInt := randomLength.Int64()
	helloworldName := helpers.GenerateRandomString(int(randomLengthInt))

	// setup expected
	expectedDockerCmdString := fmt.Sprintf("docker run --name %s --rm hello-world", helloworldName)
	expectedDockerCmd := exec.Command("/bin/sh", "-c", expectedDockerCmdString) //#nosec nosemgrep: go.lang.security.audit.dangerous-exec-command.dangerous-exec-command
	expectedOutput := expectedDockerCmd.Stdout
	err := expectedDockerCmd.Run()
	if err != nil {
		t.Errorf("Error running expected Docker command: %v", err)
	}

	dcmd := dcmd.DockerCmd{
		Name:           helloworldName,
		Image:          "hello-world",
		Tag:            "latest",
		Flags:          []string{"--rm"},
		Volumes:        []string{},
		BindMounts:     []string{},
		Networks:       []string{},
		PublishedPorts: []string{},
		Entrypoint:     "",
		Arguments:      []string{},
		Help:           "hello world example",
	}
	actualOutput := os.Stdout
	RunDockerCommand(dcmd, []string{})

	if !cmp.Equal(actualOutput, expectedOutput, cmpopts.IgnoreUnexported()) {
		t.Errorf("Actual Docker command run does not match expected. Got %v, expected %v", actualOutput, expectedOutput)
	}

}
