package gapbuffer

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGapBuffer(t *testing.T) {

	// if total != 10 {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	// }

	gb := NewGapBuffer()
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. "
	r := "[Ut enim ad minima]\n"
	u := "[veniam, quis nostrum exercitationem]"
	w := "Οὐχὶ ταὐτὰ παρίσταταί \nμοι γιγνώσκειν, ὦ ἄνδρες\n"
	gb.PrintCursor()
	gb.SetText(s)
	gb.PrintCursor()
	//gb.debugPrint()

	for i := 0; i < 11; i++ {
		gb.CursorNext()
	}
	gb.PrintCursor()
	gb.Insert(r)
	gb.PrintCursor()

	//gb.debugPrint()
	//gb.Backspace()
	//gb.debugPrint()

	for i := 0; i < 11; i++ {
		//gb.Delete()
		//gb.debugPrint()
		gb.Insert(u)
		//gb.debugPrint()
		gb.Insert(w)
		//gb.debugPrint()
	}
	// gb.debugPrint()
	// gb.Insert(u)
	// gb.debugPrint()
	// gb.Insert(w)
	// gb.debugPrint()
	for i := 0; i < 10; i++ {
		gb.Insert(strconv.Itoa(i))
		gb.PrintCursor()
	}
	//gb.debugPrint()

	fmt.Printf("%v\n", gb.GetText())

}
