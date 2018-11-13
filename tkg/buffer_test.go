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
	gb.PrintPoint()
	gb.SetText(s)
	gb.PrintPoint()
	//gb.debugPrint()

	for i := 0; i < 11; i++ {
		gb.PointNext()
	}
	gb.PrintPoint()
	gb.Insert(r)
	gb.PrintPoint()

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
		gb.PrintPoint()
	}
	//gb.debugPrint()

	fmt.Printf("%v %d %d\n", gb.GetText(), gb.BufferLen(), gb.Len())
	fmt.Println(1, gb.PointForLine(1))
	fmt.Println(2, gb.PointForLine(2))
	fmt.Println(3, gb.PointForLine(3))
	fmt.Println(4, gb.PointForLine(4))
	j := gb.Point()
	fmt.Println("Cur: ", j)
	gb.MoveGap(-10)
	j = gb.Point()
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
	gb.PrintPoint()
	gb.SetText(s)
	for i := 0; i < 11; i++ {
		gb.PointNext()
	}
	gb.PrintPoint()
	gb.Insert(r)
	gb.PrintPoint()
	fmt.Printf("|%v|\n", gb.GetText())
}

func TestTextLines(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut\nlabore et dolore magna aliqua. "
	gb.SetText(s)
	fmt.Printf("%v\n", gb.GetText())
	fmt.Println("--- [1, 3)")
	fmt.Printf("%v\n", gb.GetTextForLines(1, 3))
	fmt.Println("--- [2, 3)")
	fmt.Printf("%v\n", gb.GetTextForLines(2, 5))

}

func TestRuneAt(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut\nlabore et dolore magna aliqua. "
	gb.SetText(s)

	k := 0
	// fmt.Printf("%c\n", gb.RuneAt(k))
	// fmt.Printf("%c\n", gb.RuneAt(k+1))
	// fmt.Printf("%c\n", gb.RuneAt(k+2))
	for k < gb.BufferLen() {
		fmt.Printf("%c", gb.RuneAt(k))
		k++
	}
	fmt.Println("|-")

}
