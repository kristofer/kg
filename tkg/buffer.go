package tkg

import (
	"fmt"
	"log"
)

/*
 * Buffer
 */

type Buffer struct {
	data    []rune
	preLen  int
	postLen int
	Next    *Buffer /* b_next Link to next buffer_t */
	Mark    int     /* b_mark the mark */
	//Point      int     /* b_point the point */
	OrigPoint  int    /* b_cpoint the original current point, used for mutliple window displaying */
	PageStart  int    /* b_page start of page */
	PageEnd    int    /* b_epage end of page */
	Reframe    bool   /* b_reframe force a reframe of the display */
	WinCount   int    /* b_cnt count of windows referencing this buffer */
	TextSize   int    /* b_size current size of text being edited (not including gap) */
	PrevSize   int    /* b_psize previous size */
	PointRow   int    /* b_row Point row */
	PointCol   int    /* b_col Point col */
	Filename   string // b_fname[NAME_MAX + 1]; /* filename */
	Buffername string //[b_bnameSTRBUF_S];   /* buffer name */
	Flags      byte   /* char b_flags buffer flags */
	modified   bool
}

// var RootBuffer *Buffer = nil
// var CurrentBuffer *Buffer = nil

func (r *Buffer) MarkModified() {
	r.modified = true
}

// NewBuffer - Create a new Buffer
func NewBuffer() *Buffer {
	return &Buffer{}
}

func (r *Buffer) SetText(s string) {
	r.data = []rune(s)
	r.preLen = 0
	r.postLen = len(r.data)
}

func (r *Buffer) GetText() string {
	ret := make([]rune, r.preLen+r.postLen)

	copy(ret, r.data)
	copy(ret[r.preLen:], r.data[r.postStart():])

	return string(ret)
}

func (r *Buffer) GetTextForLines(l1, l2 int) string {
	pt1 := r.PointForLine(l1)
	pt2 := r.PointForLine(l2)
	fmt.Println(pt1, pt2)
	ret := make([]rune, pt2-pt1)

	for i, j := pt1, 0; i < len(ret); i, j = i+1, j+1 {
		if i == r.preLen {
			i = r.postStart()
		}
		ret[j] = r.data[i]
	}
	return string(ret)
}

func (r *Buffer) RuneAt(p int) rune {
	//log.Println("RuneAt", p, r.preLen, r.postLen, r.postStart())
	// if p < 0 {
	// 	return '\uFFFD' //'\u2318
	// }
	if p <= 0 {
		return r.data[0]
	}
	// if p > len(r.data) {
	// 	return '\u2318'
	// }
	if p <= r.preLen && r.preLen != 0 {
		return r.data[p]
	}
	if p >= r.postLen || p <= len(r.data)-1 {
		return r.data[r.postStart()+(p-r.preLen)]
	}
	if p < len(r.data) {
		log.Println("RuneAt", p, r.data[p], r.preLen, r.postLen, r.postStart())
	} else {
		log.Println("RuneAt", p, '\uFFFD', r.preLen, r.postLen, r.postStart())
	}
	return '\uFFFD' //'\u2318'
}

func (r *Buffer) AddRune(ch rune) {
	if r.gapLen() == 0 {
		_ = r.GrowGap(32)
	}

	//copy(r.data[r.gapStart():], []rune(s))
	r.data[r.preLen] = ch
	r.preLen++
}
func (r *Buffer) BufferLen() int {
	return r.preLen + r.postLen
}

func (r *Buffer) ActualLen() int {
	return len(r.data)
}

func (r *Buffer) gapStart() int {
	return r.preLen
}

func (r *Buffer) gapLen() int {
	return r.postStart() - r.preLen
}

func (r *Buffer) postStart() int {
	return len(r.data) - r.postLen
}

// Insert adds the string, growing the gap if needed.
func (r *Buffer) Insert(s string) {
	if r.gapLen() < len(s) {
		newGap := len(s) + 32
		_ = r.GrowGap(newGap)
	}

	copy(r.data[r.gapStart():], []rune(s))
	r.preLen += len(s)
	//fmt.Println("G", len(r.data)-r.postLen-r.preLen, "S", r.preLen, "E", r.postLen)
}

// GrowGap makes the gap bigger by n
// not sure why I need this.
func (r *Buffer) GrowGap(n int) bool {
	newData := make([]rune, len(r.data)+n)

	copy(newData, r.data[:r.preLen])

	copy(newData[r.postStart()+n:],
		r.data[r.postStart():])

	r.data = newData
	return true
}

// MoveGap moves the gap forward by offset runes
func (r *Buffer) MoveGap(offset int) int {

	if offset > 0 {
		if r.postLen == 0 {
			return 0
		}
		for i := 0; i < offset; i++ {
			r.data[r.preLen] = r.data[len(r.data)-r.postLen]
			r.preLen++
			r.postLen--
		}
	}
	if offset < 0 {
		if r.preLen == 0 {
			return 0
		}
		for i := offset; i < 0; i++ {
			r.data[r.postStart()-1] = r.data[r.preLen-1]
			r.preLen--
			r.postLen++
		}
	}

	return offset
}
func (r *Buffer) LineStart(point int) int {
	if point < 0 {
		return 0
	}
	p := r.RuneAt(point)
	for point >= 0 {
		point--
		p = r.RuneAt(point)
		if p == '\n' {
			point++
			return point
		}
	}
	if point <= 0 {
		return 0
	}
	return point
}

// PointForLine return point for beginning of line ln
func (r *Buffer) PointForLine(ln int) int {
	if ln < 1 {
		return 0
	}
	ep := len(r.data) - 1
	sp := 0
	for pt := 0; pt < ep; pt++ {
		if pt == r.preLen {
			pt = r.postStart()
		}
		if r.data[pt] == '\n' {
			ln--
			if ln == 0 {
				return sp
			}
			if (pt + 1) < ep {
				sp = pt + 1
			}
		}
		if (pt + 1) == ep {
			return sp
		}
	}
	return ep
}

func (r *Buffer) ColumnForPoint(point int) (column int) {
	i := point
	if r.data[i] == '\n' {
		return point
	}
	for i > 0 {
		if r.data[i] == '\n' {
			return point - i
		}
		i--
	}

	return point - i + 1
}

func (r *Buffer) LineForPoint(point int) (line int) {
	line = 1
	pt := 0
	for pt = 1; pt <= point; pt++ {
		if r.data[pt-1] == '\n' {
			line++
		}
	}
	if pt == r.BufferLen() {
		line--
	}
	return
}

// ColRowForPoint returns the cursor location for a pt in the buffer
func (r *Buffer) XYForPoint(pt int) (x, y int) {
	x, y = 0, 0
	x = r.ColumnForPoint(pt)
	y = r.LineForPoint(pt)
	return
}

// Delete remove a rune forward
func (r *Buffer) Delete() {
	if r.postLen == 0 {
		return
	}

	r.postLen--
}

// Backspace remove a rune backward
func (r *Buffer) Backspace() {
	if r.preLen == 0 {
		return
	}

	r.preLen--
}

// Point return point
func (r *Buffer) Point() int {
	return r.preLen
}
func (r *Buffer) RuneForPoint() rune {
	return r.data[r.preLen]
}

func (r *Buffer) SetPoint(np int) {

}

// PrintPoint print Point point
func (r *Buffer) PrintPoint() {
	fmt.Println("C: ", r.Point())
}

// PointNext move point left one
func (r *Buffer) PointNext() {
	if r.postLen == 0 {
		return
	}

	r.data[r.preLen] = r.data[r.postStart()]
	r.preLen++
	r.postLen--
}

// PointPrevious move point right one
func (r *Buffer) PointPrevious() {
	if r.preLen == 0 {
		return
	}

	r.data[r.postStart()-1] = r.data[r.preLen-1]
	r.preLen--
	r.postLen++
}

/* GetLineStats scan buffer and fill in curline and lastline */
func (r *Buffer) GetLineStats() (curline int, lastline int) {
	ep := len(r.data) - 1
	line := 0
	curline = -1
	for p := 0; p < ep; p++ {
		if p == r.preLen {
			p = r.postStart()
		}
		if r.data[p] == '\n' {
			line++
			lastline = line
		}
		if curline == -1 && p == r.preLen {
			if r.data[p] == '\n' {
				curline = line
			} else {
				curline = line + 1
			}
		}
	}
	return curline, lastline
}

func (r *Buffer) debugPrint() {
	fmt.Printf("*********\n")
	for i := 0; i < len(r.data); i++ {
		if i >= r.gapStart() && i < r.gapStart()+r.gapLen() {
			fmt.Printf("@")
		} else if i < r.preLen {
			fmt.Printf("%c", r.data[i])
		} else {
			fmt.Printf("%c", r.data[i])
		}
	}

	fmt.Printf("\n")
}
