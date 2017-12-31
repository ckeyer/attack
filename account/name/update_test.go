package name

import (
	"testing"
)

func TestDownload(t *testing.T) {
	err := Update()
	if err != nil {
		t.Fatal(err)
	}

	t.Error("...")
}
