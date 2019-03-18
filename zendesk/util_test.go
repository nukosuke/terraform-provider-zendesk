package zendesk

import "testing"

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
