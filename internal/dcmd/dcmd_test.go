package dcmd

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
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
	if !cmp.Equal(actualPaths, expectedPaths) {
		t.Errorf("Actual cmd paths do not match expected. Got %v, expected %v", actualPaths, expectedPaths)
	}
}

func TestLoadCmds(t *testing.T) {
	// Set up test data
	expectedHelloWorld := DockerCmd{
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
		UserArguments:  []string{},
		Help:           "hello world example",
	}
	expectedMd2Html := DockerCmd{
		Name:           "md2html",
		Image:          "mbentley/grip",
		Tag:            "latest",
		Flags:          []string{"-i", "--rm"},
		Volumes:        []string{},
		BindMounts:     []string{"${PWD}:/data", "~/.grip:/.grip"},
		Networks:       []string{},
		PublishedPorts: []string{"8080:8080"},
		Entrypoint:     "",
		Arguments:      []string{"--context=username/repo README.md 0.0.0.0:8080"},
		UserArguments:  []string{},
		Help:           "Display Markdown file as HTML webpage",
	}
	expectedCmdMap := make(map[string]DockerCmd)
	expectedCmdMap["helloworld"] = expectedHelloWorld
	expectedCmdMap["md2html"] = expectedMd2Html

	// Call LoadCmds and check the result
	actualCmdMap, err := LoadCmds("./testdata/dcmd")
	if err != nil {
		t.Errorf("Error loading config: %v", err)
	}

	// Use reflection to iterate through fields
	actualType := reflect.TypeOf(actualCmdMap).Elem()
	expectedType := reflect.TypeOf(expectedCmdMap).Elem()

	for i := 0; i < actualType.NumField(); i++ {
		actualField := actualType.Field(i)
		expectedField := expectedType.Field(i)

		// Compare exported fields
		if actualField.Name != expectedField.Name {
			t.Errorf("Field names do not match: Expected %s, Actual %s", expectedField.Name, actualField.Name)
		}

		actualValue := reflect.ValueOf(actualField).Interface()
		expectedValue := reflect.ValueOf(expectedField).Interface()

		if !reflect.DeepEqual(actualValue, expectedValue) {
			t.Errorf("Field values do not match: Expected %v, Actual %v", expectedValue, actualValue)
		}
	}

}
