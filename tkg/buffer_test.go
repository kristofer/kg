package tkg

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestAddRune2(t *testing.T) {
	gb := NewBuffer()
	//s := "fooLorem\nlite\nsed ut\naliqua.-\nhhh"
	//      012345678 91123 4567892 123456789 3123456
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)
	gb.DebugPrint()
	//gb.Insert("foo")
	gb.AddRune('f')
	gb.AddRune('o')
	gb.AddRune('o')

	gb.DebugPrint()

	t.Error("End of Buffer")

}
func TestCollapseGap(t *testing.T) {
	gb := NewBuffer()
	//s := "fooLorem\nlite\nsed ut\naliqua.-\nhhh"
	//      012345678 91123 4567892 123456789 3123456
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)
	gb.DebugPrint()
	//gb.Insert("foo")
	gb.AddRune('f')
	gb.AddRune('o')
	gb.AddRune('o')

	gb.DebugPrint()

	gb.CollapseGap()

	gb.DebugPrint()
	t.Error("End of Buffer")

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
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)
	// gb.Insert("foo")
	// gb.AddRune('k')
	// gb.AddRune('r')
	// gb.AddRune('i')
	// gb.AddRune('s')
	assert.Equal(t, 0, gb.LineStart(0))
	assert.Equal(t, 0, gb.LineStart(1))
	assert.Equal(t, 0, gb.LineStart(4))
	assert.Equal(t, 0, gb.LineStart(5))
	assert.Equal(t, 6, gb.LineStart(6))
	assert.Equal(t, 6, gb.LineStart(8))
	assert.Equal(t, 6, gb.LineStart(10))
	assert.Equal(t, 11, gb.LineStart(11))
	assert.Equal(t, 11, gb.LineStart(13))
	assert.Equal(t, 11, gb.LineStart(13))
	assert.Equal(t, 18, gb.LineStart(21))
	assert.Equal(t, 18, gb.LineStart(25))
	assert.Equal(t, 18, gb.LineStart(26))
	assert.Equal(t, 27, gb.LineStart(27))
	assert.Equal(t, 27, gb.LineStart(29))
	gb.DebugPrint()
}
func TestLineStart2(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    012345678 91123 4567892 123456789 3123456
	//s := "fooLorem\nlite\nsed ut\naliqua.-\nhhh"
	//      012345678 91123 4567892 123456789 3123456
	gb.SetText(s)
	//gb.Insert("foo")
	gb.AddRune('f')
	gb.AddRune('o')
	gb.AddRune('o')
	// gb.AddRune('s')
	assert.Equal(t, 0, gb.LineStart(0))
	assert.Equal(t, 0, gb.LineStart(1))
	assert.Equal(t, 0, gb.LineStart(4))
	assert.Equal(t, 0, gb.LineStart(5))
	assert.Equal(t, 0, gb.LineStart(6))
	assert.Equal(t, 0, gb.LineStart(8))
	assert.Equal(t, 9, gb.LineStart(9))
	assert.Equal(t, 9, gb.LineStart(10))
	assert.Equal(t, 9, gb.LineStart(11))
	assert.Equal(t, 9, gb.LineStart(13))
	assert.Equal(t, 14, gb.LineStart(14))
	assert.Equal(t, 21, gb.LineStart(21))
	assert.Equal(t, 21, gb.LineStart(25))
	assert.Equal(t, 21, gb.LineStart(26))
	assert.Equal(t, 21, gb.LineStart(27))
	assert.Equal(t, 21, gb.LineStart(29))
	assert.Equal(t, 30, gb.LineStart(30))
	assert.Equal(t, 30, gb.LineStart(32))
	assert.Equal(t, 30, gb.LineStart(33))
	assert.Equal(t, 30, gb.LineStart(35))
	assert.Equal(t, 30, gb.LineStart(410))
	gb.DebugPrint()
}

func TestLineEnd(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.SetText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	assert.Equal(t, 7, gb.LineEnd(0))
	assert.Equal(t, 7, gb.LineEnd(1))
	assert.Equal(t, 7, gb.LineEnd(2))
	assert.Equal(t, 7, gb.LineEnd(5))
	assert.Equal(t, 13, gb.LineEnd(11))
	assert.Equal(t, 13, gb.LineEnd(13))
	assert.Equal(t, 25, gb.LineEnd(21))
	assert.Equal(t, 34, gb.LineEnd(27))
}

func TestLineLenAtPoint(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)
	// gb.Insert("foo")
	// gb.AddRune('k')
	// gb.AddRune('r')
	// gb.AddRune('i')
	// gb.AddRune('s')
	// gb.AddRune('\n')
	assert.Equal(t, 6, gb.LineLenAtPoint(0))
	assert.Equal(t, 6, gb.LineLenAtPoint(1))
	assert.Equal(t, 6, gb.LineLenAtPoint(2))
	assert.Equal(t, 6, gb.LineLenAtPoint(5))
	assert.Equal(t, 5, gb.LineLenAtPoint(6))
	assert.Equal(t, 5, gb.LineLenAtPoint(8))
	assert.Equal(t, 5, gb.LineLenAtPoint(10))
	assert.Equal(t, 7, gb.LineLenAtPoint(11))
	assert.Equal(t, 7, gb.LineLenAtPoint(13))
	assert.Equal(t, 9, gb.LineLenAtPoint(20))
	assert.Equal(t, 9, gb.LineLenAtPoint(22))
	assert.Equal(t, 9, gb.LineLenAtPoint(26))
	assert.Equal(t, 3, gb.LineLenAtPoint(27))
	assert.Equal(t, 3, gb.LineLenAtPoint(28))
	assert.Equal(t, 3, gb.LineLenAtPoint(gb.BufferLen()-1))
	// assert.Equal(t, 27, gb.ColumnForPoint(5))
	// assert.Equal(t, 27, gb.ColumnForPoint(6))
	// assert.Equal(t, 27, gb.ColumnForPoint(100))
}
func TestLineLenAtPoint2(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.SetText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	assert.Equal(t, 8, gb.LineLenAtPoint(0))
	assert.Equal(t, 8, gb.LineLenAtPoint(1))
	assert.Equal(t, 8, gb.LineLenAtPoint(2))
	assert.Equal(t, 8, gb.LineLenAtPoint(5))
	assert.Equal(t, 8, gb.LineLenAtPoint(6))
	gb.DebugPrint()

	assert.Equal(t, 6, gb.LineLenAtPoint(8))
	assert.Equal(t, 6, gb.LineLenAtPoint(10))
	assert.Equal(t, 6, gb.LineLenAtPoint(11))
	assert.Equal(t, 6, gb.LineLenAtPoint(13))
	assert.Equal(t, 7, gb.LineLenAtPoint(20))
	assert.Equal(t, 7, gb.LineLenAtPoint(22))
	assert.Equal(t, 9, gb.LineLenAtPoint(26))
	assert.Equal(t, 9, gb.LineLenAtPoint(27))
	assert.Equal(t, 9, gb.LineLenAtPoint(28))
	assert.Equal(t, 9, gb.LineLenAtPoint(34))
	assert.Equal(t, 3, gb.LineLenAtPoint(35))
	assert.Equal(t, 3, gb.LineLenAtPoint(36))
	assert.Equal(t, 3, gb.LineLenAtPoint(37))
	assert.Equal(t, 3, gb.LineLenAtPoint(gb.BufferLen()-1))
	// assert.Equal(t, 27, gb.ColumnForPoint(5))
	// assert.Equal(t, 27, gb.ColumnForPoint(6))
	// assert.Equal(t, 27, gb.ColumnForPoint(100))
}
func TestPointForLine(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.SetText(s)
	gb.DebugPrint()
	// gb.SetText(s)
	// gb.Insert("foo")
	// gb.AddRune('k')
	// gb.AddRune('r')
	// gb.AddRune('i')
	// gb.AddRune('s')
	// gb.AddRune('\n')
	assert.Equal(t, 0, gb.PointForLine(0))
	assert.Equal(t, 0, gb.PointForLine(1))
	assert.Equal(t, 6, gb.PointForLine(2))
	assert.Equal(t, 11, gb.PointForLine(3))
	assert.Equal(t, 18, gb.PointForLine(4))
	assert.Equal(t, 27, gb.PointForLine(5))
	assert.Equal(t, 27, gb.PointForLine(6))
	assert.Equal(t, 27, gb.PointForLine(100))
	gb.DebugPrint()
}
func TestPointForLine2(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.SetText(s)
	// gb.SetText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	gb.DebugPrint()
	assert.Equal(t, 0, gb.PointForLine(1))
	assert.Equal(t, 8, gb.PointForLine(2))
	assert.Equal(t, 14, gb.PointForLine(3))
	assert.Equal(t, 19, gb.PointForLine(4))
	assert.Equal(t, 26, gb.PointForLine(5))
	assert.Equal(t, 35, gb.PointForLine(6))
	assert.Equal(t, 37, gb.PointForLine(100))
}
func TestColumnForPoint(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.SetText(s)
	gb.SetText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
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
func TestPointForXY(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.SetText(s)
	gb.SetText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
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
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	// gb.PointNext()
	// gb.PointNext()
	// gb.PointNext()
	// gb.PointPrevious()
	fmt.Printf("*********\n")
	k := 0
	for k < gb.BufferLen()-1 {
		rch, err := gb.RuneAt(k)
		if err != nil {
			t.Errorf("k %d rch %c %s", k, rch, err)
		}
		if rch == '\n' {
			fmt.Printf("%c\n", 0x00B6)
		} else {
			fmt.Printf("%c", rch)
		}
		k++
	}
	fmt.Printf("\n*********\n")
	gb.DebugPrint()
	t.Error("End of Buffer")
}
