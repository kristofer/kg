package kg

import (
	"testing"
)

func TestEditor(t *testing.T) {
	edit := &Editor{}
	testy := []string{"kg"} // , "docs/UTF8.txt"}
	edit.StartEditor(testy, len(testy))
}
