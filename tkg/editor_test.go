package tkg

import (
	"testing"
)

func TestEditor(t *testing.T) {
	edit := &Editor{}
	edit.StartEditor(testy, len(testy))
}
