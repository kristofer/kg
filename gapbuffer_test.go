package kg

import "testing"

func TestGapBuffer(t *testing.T) {

	// if total != 10 {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	// }

	gb := NewGapBuffer()
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. "
	r := "[Ut enim ad minima veniam, quis nostrum exercitationem]"
	gb.SetText(s)

	gb.Insert(r)

	gb.debugPrint()
}
