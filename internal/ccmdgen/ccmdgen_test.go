package ccmdgen

import (
	"reflect"
	"testing"

	"github.com/spf13/cobra"
	"github.com/toozej/finas/internal/dcmd"
	"github.com/toozej/finas/internal/drun"
)

func TestNewFinasCommand(t *testing.T) {
	dcmd := dcmd.DockerCmd{
		Name:  "helloworld",
		Image: "hello-world",
		Flags: []string{"--rm"},
		Help:  "hello world example",
	}

	expectedCmd := &cobra.Command{
		Use:   "helloworld",
		Short: "hello world example",
		Run: func(c *cobra.Command, args []string) {
			drun.RunDockerCommand(dcmd, args)
		},
	}

	actualCmd := NewFinasCommand(dcmd)

	// Use reflection to iterate through fields
	actualType := reflect.TypeOf(actualCmd).Elem()
	expectedType := reflect.TypeOf(expectedCmd).Elem()

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
