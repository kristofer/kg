package tkg

import (
	"errors"
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

func (r *Buffer) RuneAt(p int) (rune, error) {
	if p >= len(r.data) {
		return 0, errors.New("Beyond data buffer in RuneAt")
	}
	if p < 0 {
		return '\u0000', errors.New("negative buffer pointer in RuneAt")
	}
	if p <= r.preLen {
		return r.data[p], nil
	}
	if p > r.preLen && p < len(r.data) {
		p -= r.preLen
		npt := r.postStart() + p
		if npt >= len(r.data) {
			log.Println("pt len gap s e l", r.postStart()+p, len(r.data), r.gapStart(), r.postStart(), r.gapLen())
			return 0, errors.New("Ran over end of data buffer in RuneAt")
		}
		return r.data[npt], nil
	}
	if p < len(r.data)-1 {
		log.Println("RuneAt", p, r.data[p], r.preLen, r.postLen, r.postStart())
	} else {
		log.Println("RuneAt", p, '\uFFFD', r.preLen, r.postLen, r.postStart())
	}
	return 0, errors.New("error at end of RuneAt") //'\u2318'
}

// AddRune add a run to the buffer
func (r *Buffer) AddRune(ch rune) {
	if r.gapLen() == 0 {
		_ = r.GrowGap(CHUNK)
	}
	r.data[r.preLen] = ch
	r.preLen++
}

// Point return point
func (r *Buffer) Point() int {
	return r.preLen
}

// Return rune at Point
// func (r *Buffer) RuneForPoint() rune {
// 	return r.data[r.preLen]
// }

// SetPoint set thecurrent point to np
func (r *Buffer) SetPoint(np int) {
	// 	slade gap to end
	copy(r.data[r.preLen:], r.data[r.gapStart():])
	r.preLen += r.postLen
	r.postLen = 0
	// move gap left by np chars
	for i := np; i < 0; i-- {
		r.data[r.postStart()-1] = r.data[r.preLen-1]
		r.preLen--
		r.postLen++
	}
}

// PrintPoint print Point point
func (r *Buffer) PrintPoint() {
	fmt.Println("C: ", r.Point())
}

// BufferLen length of buffer
func (r *Buffer) BufferLen() int {
	return r.preLen + r.postLen
}

// ActualLen length of buffer plus gap
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

// CollapseGap moves the gap to the end of the buffer for replacement
func (r *Buffer) CollapseGap() {
	copy(r.data[r.preLen:], r.data[r.gapStart():])
	r.preLen += r.postLen
	r.postLen = 0

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

// MoveGap moves the gap to a Point
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
	p, err := r.RuneAt(point)
	if err != nil {
		panic(err)
	}
	for point > 0 {
		point--
		p, err = r.RuneAt(point)
		if err != nil {
			panic(err)
		}
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

// LineEnd find the point at end of this line
func (r *Buffer) LineEnd(point int) int {
	if point < 0 {
		return 0
	}
	ep := len(r.data)
	for {
		if point == r.preLen {
			point = r.postStart()
		}
		if point >= ep {
			return ep
		}
		p, err := r.RuneAt(point)
		if err != nil {
			panic(err)
		}
		if p == '\n' {
			return point
		}
		point++
	}
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

// XYForPoint returns the cursor location for a pt in the buffer
func (r *Buffer) XYForPoint(pt int) (x, y int) {
	x, y = 0, 0
	x = r.ColumnForPoint(pt)
	y = r.LineForPoint(pt)
	return
}

// PointForXY returns the Point location for X, Y in the buffer
func (bp *Buffer) PointForXY(x, y int) (finalpt int) {
	c := 1
	r := 1
	lch := bp.data[0] // last rune
	lpt := 0
	if (c == x) && (r == y) {
		return 0
	}
	ep := len(bp.data) - 1
	for pt := 1; pt < ep; pt++ {
		if pt == bp.preLen+1 { // jump over gap
			pt = bp.postStart()
		}
		if lch == '\n' {
			if (r-1 == y) && (c <= x) {
				return lpt
			}
			r++
			c = 1
		} else {
			c++
		}
		if (c == x) && (r == y) {
			return pt
		}
		lch = bp.data[pt]
		lpt = pt
	}
	return ep
}

// SegStart Forward scan for start of logical line segment
// (corresponds to screen line)  containing 'finish'
func (r *Buffer) SegStart(start, finish, limit int) int {
	//var p rune
	c := 0
	scan := start

	for scan < finish {
		//p = ptr(bp, scan);
		p, err := r.RuneAt(scan)
		if err != nil {
			panic(err)
		}

		if p == '\n' {
			c = 0
			start = scan + 1
		} else {
			if limit <= c {
				c = 0
				start = scan
			}
		}
		scan++
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if p == '\t' {
			c += 4 //8 - (c % 7)
		} else {
			c++
		}
	}
	// (c < COLS ? start : finish);
	if c < limit {
		return start
	}
	return finish
}

/* SegNext Forward scan for start of logical line segment following 'finish' */
func (r *Buffer) SegNext(start, finish, limit int) int {
	// char_t *p;
	// int c = 0;
	//bp := e.CurrentBuffer
	//var p rune
	//var pptr int
	c := 0

	scan := r.SegStart(start, finish, limit)
	for {
		p, err := r.RuneAt(scan)
		if err != nil {
			panic(err)
		}
		//if (bp.b_ebuf <= p || COLS <= c)
		if limit <= c {
			break
		}
		//scan += utf8_size(*ptr(bp,scan));
		scan++
		if p == '\n' {
			break
		}
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if p == '\t' {
			c += 4 //8 - (c % 7)
		} else {
			c++
		}
	}
	//(p < bp.b_ebuf ? scan : );
	if scan < r.BufferLen() {
		return scan
	}
	return r.BufferLen()
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

// PointUp move point up one line
func (r *Buffer) PointUp() {
	c1 := r.ColumnForPoint(r.Point())
	l1 := r.LineStart(r.Point())
	l1--
	l2 := r.LineStart(l1)
	r.SetPoint(l2 + c1)
}

// PointDown move point down one line
func (r *Buffer) PointDown() {

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
