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
	}

	actualPaths, err := cmdPaths("./testdata/dcmd/")
	if err != nil {
		t.Errorf("Error getting cmd paths: %v", err)
	}
	if !reflect.DeepEqual(actualPaths, expectedPaths) {
		t.Errorf("Actual cmd paths do not match expected. Got %v, expected %v", actualPaths, expectedPaths)
	}
}

func TestLoadCmds(t *testing.T) {
	// Set up test data
	expectedHelloWorld := DockerCmd{}
	expectedMd2Html := DockerCmd{}
	expectedCmdMap := make(map[string]DockerCmd)
	expectedCmdMap["Key1"] = expectedHelloWorld
	expectedCmdMap["Key2"] = expectedMd2Html

	// Call LoadCmds and check the result
	err := LoadCmds("./testdata/dcmd/")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}
	if !reflect.DeepEqual(CmdMap, expectedCmdMap) {
		t.Errorf("Loaded command map does not match expected. Got %v, expected %v", CmdMap, expectedCmdMap)
	}
}
