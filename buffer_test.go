package kg

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
	gb.setText(s)
	for i := 0; i < 11; i++ {
		gb.PointNext()
	}
	gb.Insert(r)
	fmt.Printf("|%v|\n", gb.getText())
}

func TestAddRune2(t *testing.T) {
	gb := NewBuffer()
	//s := "fooLorem\nlite\nsed ut\naliqua.-\nhhh"
	//      012345678 91123 4567892 123456789 3123456
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.setText(s)
	gb.DebugPrint()
	//gb.Insert("foo")
	gb.AddRune('f')
	gb.AddRune('o')
	gb.AddRune('o')

	gb.DebugPrint()

	//t.Error("End of Buffer")

}
func TestCollapseGap(t *testing.T) {
	gb := NewBuffer()
	//s := "fooLorem\nlite\nsed ut\naliqua.-\nhhh"
	//      012345678 91123 4567892 123456789 3123456
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.setText(s)
	gb.DebugPrint()
	//gb.Insert("foo")
	gb.AddRune('f')
	gb.AddRune('o')
	gb.AddRune('o')

	gb.DebugPrint()

	gb.CollapseGap()

	gb.DebugPrint()
	//t.Error("End of Buffer")

}

func TestTextLines(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut\nlabore et dolore magna aliqua. "
	//    0123456789112345678921234567 89312345678941234567895123456 7896123456789612345678971234567898123456789
	gb.setText(s)
	//fmt.Printf("%v\n", gb.getText())
	fmt.Println("--- [1, 3)")
	fmt.Printf("%v\n", gb.getTextForLines(1, 3))
	fmt.Println("--- [1, 2)")
	fmt.Printf("%v\n", gb.getTextForLines(1, 2))
	fmt.Println("--- [2, 3)")
	fmt.Printf("%v\n", gb.getTextForLines(2, 3))
	fmt.Println("--- [2, 5)")
	fmt.Printf("%v\n", gb.getTextForLines(2, 5))
	gb.DebugPrint()
	//t.Error("force print")
}

func TestLineStart(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	gb.setText(s)
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
	gb.setText(s)
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

func TestLineForPoint(t *testing.T) {
	gb := NewBuffer()
	//    111111 22222 3333333 444444444 555
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    012345 67891 1234567 892123456 7893123456
	//s := "fooLorem\nlite\nsed ut\naliqua.-\nhhh"
	//      012345678 91123 4567892 123456789 3123456
	gb.setText(s)
	//gb.Insert("foo")
	// gb.AddRune('f')
	// gb.AddRune('o')
	// gb.AddRune('o')
	// gb.AddRune('s')
	assert.Equal(t, 1, gb.LineForPoint(0))
	assert.Equal(t, 1, gb.LineForPoint(1))
	assert.Equal(t, 1, gb.LineForPoint(4))
	assert.Equal(t, 1, gb.LineForPoint(5))
	assert.Equal(t, 2, gb.LineForPoint(6))
	assert.Equal(t, 2, gb.LineForPoint(8))
	assert.Equal(t, 2, gb.LineForPoint(9))
	assert.Equal(t, 2, gb.LineForPoint(10))
	assert.Equal(t, 3, gb.LineForPoint(11))
	assert.Equal(t, 3, gb.LineForPoint(13))
	assert.Equal(t, 3, gb.LineForPoint(14))
	assert.Equal(t, 4, gb.LineForPoint(21))
	assert.Equal(t, 4, gb.LineForPoint(25))
	assert.Equal(t, 4, gb.LineForPoint(26))
	assert.Equal(t, 5, gb.LineForPoint(27))
	assert.Equal(t, 5, gb.LineForPoint(28))
	assert.Equal(t, 5, gb.LineForPoint(29))
	assert.Equal(t, 5, gb.LineForPoint(30))
	assert.Equal(t, 5, gb.LineForPoint(31))
	assert.Equal(t, 5, gb.LineForPoint(32))
	assert.Equal(t, 5, gb.LineForPoint(33))
	// assert.Equal(t, 30, gb.LineForPoint(33))
	// assert.Equal(t, 30, gb.LineForPoint(35))
	// assert.Equal(t, 30, gb.LineForPoint(410))
	gb.DebugPrint()
}
func TestLineForPoint2(t *testing.T) {
	gb := NewBuffer()
	//    111111 22222 3333333 444444444 555
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    012345 67891 1234567 892123456 7893123456
	//      111111111 22222 3333333 444444444 555
	//s := "fooLorem\nlite\nsed ut\naliqua.-\nhhh"
	//      012345678 91123 4567892 123456789 3123456
	gb.setText(s)
	gb.Insert("foo")
	// gb.AddRune('f')
	// gb.AddRune('o')
	// gb.AddRune('o')
	// gb.AddRune('s')
	assert.Equal(t, 1, gb.LineForPoint(0))
	assert.Equal(t, 1, gb.LineForPoint(1))
	assert.Equal(t, 1, gb.LineForPoint(4))
	assert.Equal(t, 1, gb.LineForPoint(5))
	assert.Equal(t, 1, gb.LineForPoint(6))
	assert.Equal(t, 1, gb.LineForPoint(8))
	assert.Equal(t, 2, gb.LineForPoint(9))
	assert.Equal(t, 2, gb.LineForPoint(10))
	assert.Equal(t, 2, gb.LineForPoint(11))
	assert.Equal(t, 2, gb.LineForPoint(13))
	assert.Equal(t, 3, gb.LineForPoint(14))
	assert.Equal(t, 4, gb.LineForPoint(21))
	assert.Equal(t, 4, gb.LineForPoint(25))
	assert.Equal(t, 4, gb.LineForPoint(26))
	assert.Equal(t, 4, gb.LineForPoint(27))
	assert.Equal(t, 4, gb.LineForPoint(28))
	assert.Equal(t, 4, gb.LineForPoint(29))
	assert.Equal(t, 5, gb.LineForPoint(30))
	assert.Equal(t, 5, gb.LineForPoint(31))
	assert.Equal(t, 5, gb.LineForPoint(32))
	assert.Equal(t, 5, gb.LineForPoint(33))
	assert.Equal(t, 5, gb.LineForPoint(100))
	// assert.Equal(t, 30, gb.LineForPoint(33))
	// assert.Equal(t, 30, gb.LineForPoint(35))
	// assert.Equal(t, 30, gb.LineForPoint(410))
	gb.DebugPrint()
}

func TestLineEnd(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
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
	gb.setText(s)
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
	assert.Equal(t, 3, gb.LineLenAtPoint(gb.TextSize-1))
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
	gb.setText(s)
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
	assert.Equal(t, 3, gb.LineLenAtPoint(gb.TextSize-1))
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
	gb.setText(s)
	gb.DebugPrint()
	// gb.setText(s)
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
	assert.Equal(t, 29, gb.PointForLine(5))
	assert.Equal(t, 29, gb.PointForLine(6))
	assert.Equal(t, 29, gb.PointForLine(100))
	gb.DebugPrint()
}
func TestPointForLine2(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	// gb.setText(s)
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
	assert.Equal(t, 37, gb.PointForLine(6))
	assert.Equal(t, 37, gb.PointForLine(100))
}
func TestColumnForPoint(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	gb.setText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	assert.Equal(t, 1, gb.ColumnForPoint(0))
	assert.Equal(t, 2, gb.ColumnForPoint(1))
	assert.Equal(t, 3, gb.ColumnForPoint(2))
	assert.Equal(t, 8, gb.ColumnForPoint(7))
	assert.Equal(t, 1, gb.ColumnForPoint(8))
	assert.Equal(t, 4, gb.ColumnForPoint(11))
	assert.Equal(t, 6, gb.ColumnForPoint(13))
	assert.Equal(t, 1, gb.ColumnForPoint(14))
	assert.Equal(t, 1, gb.ColumnForPoint(19))
	assert.Equal(t, 2, gb.ColumnForPoint(27))
	assert.Equal(t, 3, gb.ColumnForPoint(28))
	assert.Equal(t, 4, gb.ColumnForPoint(29))
	assert.Equal(t, 5, gb.ColumnForPoint(30))
	assert.Equal(t, 2, gb.ColumnForPoint(36))
	assert.Equal(t, 3, gb.ColumnForPoint(37))
	assert.Equal(t, 3, gb.ColumnForPoint(38))
	assert.Equal(t, 3, gb.ColumnForPoint(39))
	// assert.Equal(t, 27, gb.ColumnForPoint(100))
}

func TestPointForXY(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	gb.DebugPrint()
	//assert.Equal(t, "1, 1)", fmt.Sprint("%d, %d", gb.TestXYForPoint(0))
	assert.Equal(t, 0, gb.PointForXY(1, 1))
	assert.Equal(t, 1, gb.PointForXY(2, 1))
	assert.Equal(t, 3, gb.PointForXY(4, 1))
	assert.Equal(t, 4, gb.PointForXY(5, 1))
	assert.Equal(t, 5, gb.PointForXY(6, 1))
	assert.Equal(t, 6, gb.PointForXY(7, 1))
	//assert.Equal(t, 4, gb.PointForXY(10, 1))
	//assert.Equal(t, 4, gb.PointForXY(32, 1))
	assert.Equal(t, 8, gb.PointForXY(1, 2))
	assert.Equal(t, 14, gb.PointForXY(1, 3))
	assert.Equal(t, 16, gb.PointForXY(3, 3))
	assert.Equal(t, 19, gb.PointForXY(1, 4))
	assert.Equal(t, 20, gb.PointForXY(2, 4))
	assert.Equal(t, 21, gb.PointForXY(3, 4))
	//assert.Equal(t, gb.TextSize-1, gb.PointForXY(6, 6))

}

func TestXYForPoint(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	gb.DebugPrint()
	x, y := 0, 0
	//assert.Equal(t, "1, 1)", fmt.Sprint("%d, %d", gb.TestXYForPoint(0))
	x, y = gb.XYForPoint(0)
	assert.Equal(t, 1, x)
	assert.Equal(t, 1, y)
	// assert.Equal(t, 1, gb.PointForXY(2, 1))
	x, y = gb.XYForPoint(1)
	assert.Equal(t, 2, x)
	assert.Equal(t, 1, y)
	// assert.Equal(t, 3, gb.PointForXY(4, 1))
	x, y = gb.XYForPoint(3)
	assert.Equal(t, 4, x)
	assert.Equal(t, 1, y)
	// assert.Equal(t, 4, gb.PointForXY(5, 1))
	x, y = gb.XYForPoint(4)
	assert.Equal(t, 5, x)
	assert.Equal(t, 1, y)
	// assert.Equal(t, 5, gb.PointForXY(6, 1))
	x, y = gb.XYForPoint(5)
	assert.Equal(t, 6, x)
	assert.Equal(t, 1, y)
	// assert.Equal(t, 14, gb.PointForXY(1, 3))
	x, y = gb.XYForPoint(14)
	assert.Equal(t, 1, x)
	assert.Equal(t, 3, y)
	// assert.Equal(t, 16, gb.PointForXY(3, 3))
	x, y = gb.XYForPoint(16)
	assert.Equal(t, 3, x)
	assert.Equal(t, 3, y)
	// assert.Equal(t, 19, gb.PointForXY(1, 4))
	x, y = gb.XYForPoint(19)
	assert.Equal(t, 1, x)
	assert.Equal(t, 4, y)
	// assert.Equal(t, 20, gb.PointForXY(2, 4))
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	// assert.Equal(t, 21, gb.PointForXY(3, 4))
	x, y = gb.XYForPoint(21)
	assert.Equal(t, 3, x)
	assert.Equal(t, 4, y)
	//assert.Equal(t, gb.TextSize-1, gb.PointForXY(6, 6))
	x, y = gb.XYForPoint(35)
	assert.Equal(t, 1, x)
	assert.Equal(t, 6, y)
	x, y = gb.XYForPoint(36)
	assert.Equal(t, 2, x)
	assert.Equal(t, 6, y)
	x, y = gb.XYForPoint(37)
	assert.Equal(t, 3, x)
	assert.Equal(t, 6, y)
	x, y = gb.XYForPoint(38)
	assert.Equal(t, 3, x)
	assert.Equal(t, 6, y)
	x, y = gb.XYForPoint(gb.TextSize)
	assert.Equal(t, 3, x)
	assert.Equal(t, 6, y)
	x, y = gb.XYForPoint(gb.TextSize + 1)
	assert.Equal(t, 3, x)
	assert.Equal(t, 6, y)

}
func TestSetPoint(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua.-\nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLXorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 8911234 5678 9212345 678931234 56789412345
	gb.setText(s)
	gb.DebugPrint()
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')

	assert.Equal(t, 8, gb.gapStart())
	assert.Equal(t, 27, gb.gapLen())
	gb.PointNext()
	assert.Equal(t, 9, gb.gapStart())
	assert.Equal(t, 27, gb.gapLen())
	gb.AddRune('X')
	assert.Equal(t, 10, gb.gapStart())
	assert.Equal(t, 26, gb.gapLen())
	gb.DebugPrint()

	gb.SetPoint(5)
	assert.Equal(t, 5, gb.gapStart())
	assert.Equal(t, 26, gb.gapLen())
	gb.SetPoint(8)
	assert.Equal(t, 8, gb.gapStart())
	assert.Equal(t, 26, gb.gapLen())
	gb.DebugPrint()
	gb.SetPoint(15)
	assert.Equal(t, 15, gb.gapStart())
	assert.Equal(t, 26, gb.gapLen())
	gb.SetPoint(10)
	assert.Equal(t, 10, gb.gapStart())
	assert.Equal(t, 26, gb.gapLen())
	gb.AddRune('X')
	gb.AddRune('X')
	gb.DebugPrint()
	assert.Equal(t, 12, gb.gapStart())
	assert.Equal(t, 24, gb.gapLen())
	gb.DebugPrint()
	gb.SetPoint(0)
	assert.Equal(t, 0, gb.gapStart())
	assert.Equal(t, 24, gb.gapLen())
	gb.SetPoint(8)
	assert.Equal(t, 8, gb.gapStart())
	assert.Equal(t, 24, gb.gapLen())
	gb.DebugPrint()
	gb.SetPoint(15)
	assert.Equal(t, 15, gb.gapStart())
	assert.Equal(t, 24, gb.gapLen())
	gb.SetPoint(10)
	assert.Equal(t, 10, gb.gapStart())
	assert.Equal(t, 24, gb.gapLen())

	gb.SetPoint(36)
	assert.Equal(t, 36, gb.gapStart())
	assert.Equal(t, 24, gb.gapLen())
	gb.SetPoint(23)
	assert.Equal(t, 23, gb.gapStart())
	assert.Equal(t, 24, gb.gapLen())
	gb.Backspace()
	gb.Backspace()
	gb.SetPoint(23)
	assert.Equal(t, 23, gb.gapStart())
	assert.Equal(t, 26, gb.gapLen())
	gb.SetPoint(23)
	assert.Equal(t, 23, gb.gapStart())
	assert.Equal(t, 26, gb.gapLen())
	gb.Insert("01234567890123456789")
	assert.Equal(t, 43, gb.gapStart())
	assert.Equal(t, 6, gb.gapLen())
	gb.SetPoint(gb.TextSize - 1)
	assert.Equal(t, 58, gb.gapStart())
	assert.Equal(t, 6, gb.gapLen())
	gb.SetPoint(0)
	assert.Equal(t, 0, gb.gapStart())
	assert.Equal(t, 6, gb.gapLen())

}

func TestRuneAt(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut\nlabore et dolore magna aliqua. "
	gb.setText(s)
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
	for k < gb.TextSize-1 {
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
	//t.Error("End of Buffer")
}

func TestSegStart(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	gb.setText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	assert.Equal(t, 8, gb.SegStart(0, 11, 50))
	assert.Equal(t, 0, gb.SegStart(0, 1, 50))
	assert.Equal(t, 0, gb.SegStart(0, 3, 50))
	assert.Equal(t, 14, gb.SegStart(0, 16, 50))
	assert.Equal(t, 14, gb.SegStart(9, 16, 50))
	assert.Equal(t, 19, gb.SegStart(9, 25, 50))
	assert.Equal(t, 26, gb.SegStart(9, 31, 50))
	assert.Equal(t, 35, gb.SegStart(32, 36, 50))
	assert.Equal(t, 38, gb.SegStart(0, 39, 50))
	assert.Equal(t, 38, gb.SegStart(35, 41, 50))

	//t.Error("end of test")
}

func TestSegNext(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	gb.setText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	assert.Equal(t, 14, gb.SegNext(0, 11, 50))
	assert.Equal(t, 8, gb.SegNext(0, 1, 50))
	assert.Equal(t, 8, gb.SegNext(0, 3, 50))
	assert.Equal(t, 19, gb.SegNext(0, 16, 50))
	assert.Equal(t, 19, gb.SegNext(9, 16, 50))
	assert.Equal(t, 26, gb.SegNext(9, 25, 50))
	assert.Equal(t, 35, gb.SegNext(9, 31, 50))
	assert.Equal(t, 38, gb.SegNext(32, 36, 50))
	assert.Equal(t, 38, gb.SegNext(0, 39, 50))
	assert.Equal(t, 38, gb.SegNext(35, 41, 50))

	//t.Error("end of test")
}

func TestUpUp(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	gb.setText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	assert.Equal(t, 0, gb.UpUp(0, 50))
	assert.Equal(t, 0, gb.UpUp(1, 50))
	assert.Equal(t, 0, gb.UpUp(3, 50))
	assert.Equal(t, 0, gb.UpUp(9, 50))
	assert.Equal(t, 8, gb.UpUp(16, 50))
	assert.Equal(t, 14, gb.UpUp(25, 50))
	assert.Equal(t, 19, gb.UpUp(31, 50))
	assert.Equal(t, 26, gb.UpUp(36, 50))
	assert.Equal(t, 35, gb.UpUp(39, 50))
	assert.Equal(t, 35, gb.UpUp(41, 50))

	//t.Error("end of test")
}

func TestDownDown(t *testing.T) {
	gb := NewBuffer()
	s := "Lorem\nlite\nsed ut\naliqua. \nhhh"
	//    01234 56789 1123456 789212345 67893
	//s := "fookris\nLorem\nlite\nsed ut\naliqua. \nhhh"
	//    	01234567 891123 45678 9212345 678931234 56789412345
	gb.setText(s)
	gb.setText(s)
	gb.Insert("foo")
	gb.AddRune('k')
	gb.AddRune('r')
	gb.AddRune('i')
	gb.AddRune('s')
	gb.AddRune('\n')
	assert.Equal(t, 8, gb.DownDown(0, 50))
	assert.Equal(t, 8, gb.DownDown(1, 50))
	assert.Equal(t, 8, gb.DownDown(3, 50))
	assert.Equal(t, 14, gb.DownDown(9, 50))
	assert.Equal(t, 19, gb.DownDown(16, 50))
	assert.Equal(t, 26, gb.DownDown(19, 50))
	assert.Equal(t, 26, gb.DownDown(25, 50))
	assert.Equal(t, 35, gb.DownDown(31, 50))
	assert.Equal(t, 38, gb.DownDown(36, 50))
	assert.Equal(t, 38, gb.DownDown(39, 50))
	assert.Equal(t, 38, gb.DownDown(41, 50))

	//t.Error("end of test")
}
