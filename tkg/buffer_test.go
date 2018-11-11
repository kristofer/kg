package tkg

import (
	"fmt"
	"strconv"
	"testing"
)

func TestBuffer(t *testing.T) {

	// if total != 10 {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	// }

	gb := NewBuffer()
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

	for i := 0; i < 4; i++ {
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

	fmt.Printf("%v %d %d\n", gb.GetText(), gb.BufferLen(), gb.Len())
	fmt.Println(1, gb.IntForLine(1))
	fmt.Println(2, gb.IntForLine(2))
	fmt.Println(3, gb.IntForLine(3))
	fmt.Println(4, gb.IntForLine(4))
	j := gb.Cursor()
	fmt.Println("Cur: ", j)
	gb.MoveGap(-10)
	j = gb.Cursor()
	fmt.Println("Cur: ", j)
	l1, l2 := gb.GetLineStats()
	fmt.Println("Lines", l1, l2)
}

func printIdxForLine(ln int) {
	fmt.Println()
}

func TestBufferGrow(t *testing.T) {

	// if total != 10 {
	// 	t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
	// }

	gb := NewBuffer()
	s := "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. "
	r := "[Ut enim ad minima]"
	//u := "[veniam, quis nostrum exercitationem]"
	//w := "Οὐχὶ ταὐτὰ παρίσταταί \nμοι γιγνώσκειν, ὦ ἄνδρες\n"
	gb.PrintCursor()
	gb.SetText(s)
	for i := 0; i < 11; i++ {
		gb.CursorNext()
	}
	gb.PrintCursor()
	gb.Insert(r)
	gb.PrintCursor()
	fmt.Printf("|%v|\n", gb.GetText())
}
