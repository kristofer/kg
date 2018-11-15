package tkg

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
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

	fmt.Printf("%v %d %d\n", gb.GetText(), gb.BufferLen(), gb.ActualLen())
	fmt.Println(1, gb.PointForLine(1))
	fmt.Println(2, gb.PointForLine(2))
	fmt.Println(3, gb.PointForLine(3))
	fmt.Println(4, gb.PointForLine(4))
	j := gb.Point()
	fmt.Println("Cur: ", j)
	//gb.MoveGap(-10)
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

func TestLineStart(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. "
	gb.SetText(s)

	assert.Equal(t, 0, gb.LineStart(0))
	assert.Equal(t, 0, gb.LineStart(1))
	assert.Equal(t, 0, gb.LineStart(4))
	assert.Equal(t, 0, gb.LineStart(5))
	assert.Equal(t, 6, gb.LineStart(6))
	assert.Equal(t, 6, gb.LineStart(8))
	assert.Equal(t, 6, gb.LineStart(10))
	assert.Equal(t, 11, gb.LineStart(11))
	assert.Equal(t, 11, gb.LineStart(13))
}
func TestPointForLine(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)

	assert.Equal(t, 0, gb.PointForLine(0))
	assert.Equal(t, 0, gb.PointForLine(1))
	assert.Equal(t, 6, gb.PointForLine(2))
	assert.Equal(t, 11, gb.PointForLine(3))
	assert.Equal(t, 18, gb.PointForLine(4))
	assert.Equal(t, 27, gb.PointForLine(5))
	assert.Equal(t, 27, gb.PointForLine(6))
	assert.Equal(t, 27, gb.PointForLine(100))
}
func TestColumnForPoint(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)

	assert.Equal(t, 1, gb.ColumnForPoint(0))
	assert.Equal(t, 2, gb.ColumnForPoint(1))
	assert.Equal(t, 3, gb.ColumnForPoint(2))
	assert.Equal(t, 5, gb.ColumnForPoint(5))
	assert.Equal(t, 1, gb.ColumnForPoint(11))
	assert.Equal(t, 3, gb.ColumnForPoint(13))
	// assert.Equal(t, 27, gb.ColumnForPoint(5))
	// assert.Equal(t, 27, gb.ColumnForPoint(6))
	// assert.Equal(t, 27, gb.ColumnForPoint(100))
}
func TestLineForPoint(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)

	assert.Equal(t, 1, gb.LineForPoint(0))
	assert.Equal(t, 1, gb.LineForPoint(1))
	assert.Equal(t, 1, gb.LineForPoint(2))
	assert.Equal(t, 1, gb.LineForPoint(5))
	assert.Equal(t, 2, gb.LineForPoint(6))
	assert.Equal(t, 3, gb.LineForPoint(11))
	assert.Equal(t, 3, gb.LineForPoint(13))
	assert.Equal(t, 4, gb.LineForPoint(20))
	assert.Equal(t, 4, gb.LineForPoint(gb.BufferLen()-1))
	// assert.Equal(t, 27, gb.ColumnForPoint(5))
	// assert.Equal(t, 27, gb.ColumnForPoint(6))
	// assert.Equal(t, 27, gb.ColumnForPoint(100))
}
func TestPointForXY(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)

	//assert.Equal(t, "1, 1)", fmt.Sprint("%d, %d", gb.TestXYForPoint(0))
	assert.Equal(t, 0, gb.PointForXY(1, 1))
	assert.Equal(t, 1, gb.PointForXY(2, 1))
	assert.Equal(t, 3, gb.PointForXY(4, 1))
	assert.Equal(t, 4, gb.PointForXY(5, 1))
	assert.Equal(t, 4, gb.PointForXY(6, 1))
	assert.Equal(t, 4, gb.PointForXY(7, 1))
	assert.Equal(t, 4, gb.PointForXY(10, 1))
	assert.Equal(t, 4, gb.PointForXY(32, 1))
	assert.Equal(t, 6, gb.PointForXY(1, 2))
	assert.Equal(t, 11, gb.PointForXY(1, 3))
	assert.Equal(t, 13, gb.PointForXY(3, 3))
	assert.Equal(t, 18, gb.PointForXY(1, 4))
	assert.Equal(t, 19, gb.PointForXY(2, 4))
	assert.Equal(t, 20, gb.PointForXY(3, 4))
	assert.Equal(t, gb.BufferLen()-1, gb.PointForXY(6, 6))

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
		rch, _ := gb.RuneAt(k)
		fmt.Printf("%c", rch)
		k++
	}
	fmt.Println("|-")

}
