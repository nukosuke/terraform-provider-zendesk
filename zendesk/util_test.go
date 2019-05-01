package zendesk

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsValidFile(t *testing.T) {
	v := isValidFile()
	_, errs := v("testdata/street.jpg", "file_path")
	if len(errs) != 0 {
		t.Fatalf("is Valid returned an error")
	}

	_, errs = v("Missing", "file_path")
	if len(errs) == 0 {
		t.Fatalf("is Valid did not return an error for missing file")
	}

	_, errs = v("testdata", "file_path")
	if len(errs) == 0 {
		t.Fatalf("is Valid did not return an error for a directory")
	}
}

const AssigneeSystemFieldEnvVar = "TF_VAR_ASSIGNEE_TICKET_FIELD_ID"
func SystemFieldVariablePreCheck(t *testing.T) {
	if v := os.Getenv(AssigneeSystemFieldEnvVar); v == "" {
		t.Fatalf("%s must be set for acceptance tests", AssigneeSystemFieldEnvVar)
	}
}

func readExampleConfig(t *testing.T, filename string) string {
	dir, err := filepath.Abs("../examples")
	if err != nil {
		t.Fatalf("Failed to resolve fixture directory. Check the path: %s", err)
	}

	filepath := filepath.Join(dir, filename)
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		t.Fatalf("Failed to read fixture. %v", err)
	}

	return string(bytes)
}

func concatExampleConfig(t *testing.T, configs ...string) string {
	builder := new(strings.Builder)
	for _, config := range configs {
		_, err := fmt.Fprintln(builder, config)
		if err != nil {
			t.Fatalf("Encountered an error while concatonating config files: %v", err)
		}
	}

	return builder.String()
}
