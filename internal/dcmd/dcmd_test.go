package dcmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestCmdPaths(t *testing.T) {
	homeDir, _ := os.UserHomeDir()

	expectedPaths := []string{
		"./config",
		filepath.Join(homeDir, ".config/finas"),
		"/usr/local/etc/finas",
		"/etc/finas",
		"./testdata/dcmd",
	}

	p := []string{"./testdata/dcmd"}
	actualPaths, err := cmdPaths(p)
	if err != nil {
		t.Errorf("Error getting cmd paths: %v", err)
	}
	if !reflect.DeepEqual(actualPaths, expectedPaths) {
		t.Errorf("Actual cmd paths do not match expected. Got %v, expected %v", actualPaths, expectedPaths)
	}
}

func TestLoadCmds(t *testing.T) {
	// Set up test data
	expectedHelloWorld := DockerCmd{
		Name:           "hello-world",
		Tag:            "latest",
		Flags:          []string{},
		Volumes:        []string{},
		BindMounts:     []string{},
		Networks:       []string{},
		PublishedPorts: []string{},
		Arguments:      []string{},
	}
	expectedMd2Html := DockerCmd{
		Name:           "mbentley/grip",
		Tag:            "latest",
		Flags:          []string{"-it", "--rm"},
		Volumes:        []string{},
		BindMounts:     []string{"${PWD}:/data", "~/.grip:/.grip"},
		Networks:       []string{},
		PublishedPorts: []string{"8080:8080"},
		Arguments:      []string{"--context=username/repo README.md 0.0.0.0:8080"},
	}
	expectedCmdMap := make(map[string]DockerCmd)
	expectedCmdMap["helloworld"] = expectedHelloWorld
	expectedCmdMap["md2html"] = expectedMd2Html

	// Call LoadCmds and check the result
	err := LoadCmds("./testdata/dcmd")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	if !reflect.DeepEqual(CmdMap, expectedCmdMap) {
		t.Errorf("Loaded command map does not match expected. Got %v, expected %v", CmdMap, expectedCmdMap)
	}
}
