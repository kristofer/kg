package tkg

import (
	"testing"
)

func TestEditor(t *testing.T) {
	edit := &Editor{}
	testy := []string{"tkg", "docs/UTF8.txt"}
	edit.StartEditor(testy, len(testy))
}
