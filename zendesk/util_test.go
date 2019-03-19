package zendesk

import (
	"io/ioutil"
	"path/filepath"
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
