package kg

import (
	"testing"
)

func TestEditor(t *testing.T) {
	edit := &Editor{}
	edit.StartEditor([]string{}, 0)
}
