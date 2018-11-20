package tkg

import (
	"errors"
	"fmt"
	"log"
	"testing"
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
	FirstLine  int    /* b_page start of page */
	LastLine   int    /* b_epage end of page */
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

// RuneAt finally have a reliable!!
func (r *Buffer) RuneAt(pt int) (rune, error) {
	if pt >= len(r.data) {
		return 0, errors.New("Beyond data buffer in RuneAt")
	}
	if pt < 0 {
		return '\u0000', errors.New("negative buffer pointer in RuneAt")
	}
	if npt := r.dataPointForBufferPoint(pt); npt < len(r.data) {
		return r.data[npt], nil
	}
	return 0, errors.New("Ran over end of data buffer in RuneAt")
}

func (r *Buffer) dataPointForBufferPoint(pt int) int {
	npt := 0
	if pt < r.preLen {
		npt = pt
	}
	if pt >= r.preLen && pt < len(r.data) {
		npt = pt + r.gapLen()
	}
	return npt
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

// SetPoint set the current point to np
func (r *Buffer) SetPoint(np int) {
	// 	slade gap to end
	r.CollapseGap()
	// move gap <-(left) by np chars
	gs := r.gapStart()
	for i := gs - np; i > 0; i-- {
		r.data[r.postStart()-1] = r.data[r.preLen-1]
		r.preLen--
		r.postLen++
	}
	if r.PageEnd < r.preLen {
		log.Println("reframing!")
		r.Reframe = true
	}
}
func (r *Buffer) SetPointAndCursor(np int) {
	r.SetPoint(np)
	x, y := r.XYForPoint(np)
	r.PointRow = y
	r.PointCol = x
}

// PrintPoint print Point point
func (r *Buffer) PrintPoint() {
	fmt.Println("C: ", r.Point())
}

// BufferLen length of buffer
func (r *Buffer) BufferLen() int {
	return r.preLen + r.postLen
}
func (r *Buffer) EndOfBuffer(pt int) bool {
	return pt >= (r.preLen + r.postLen - 1)
}

// ActualLen length of buffer plus gap
func (r *Buffer) ActualLen() int {
	return len(r.data)
}

func (r *Buffer) gapStart() int {
	return r.preLen
}
func (r *Buffer) GapStart() int {
	return r.preLen
}

func (r *Buffer) gapLen() int {
	return r.postStart() - r.preLen
}
func (r *Buffer) GapLen() int {
	return r.postStart() - r.preLen
}

func (r *Buffer) postStart() int {
	return len(r.data) - r.postLen
}

// CollapseGap moves the gap to the end of the buffer for replacement
func (r *Buffer) CollapseGap() {
	for i := r.preLen; r.postLen > 0; i++ {
		r.data[r.preLen] = r.data[len(r.data)-r.postLen]
		r.preLen++
		r.postLen--
	}
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

// GetTextForLines return string for [l1, l2) (l2 not included)
func (r *Buffer) GetTextForLines(l1, l2 int) string {
	pt1 := r.PointForLine(l1)
	pt2 := r.PointForLine(l2)
	//fmt.Println(pt1, pt2)
	ret := make([]rune, pt2-pt1)
	j := 0
	for i := pt1; j < len(ret); i++ {
		rch, err := r.RuneAt(i)
		if err != nil {
			panic(err)
		}
		ret[j] = rch
		j++
	}
	return string(ret)
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

	if offset < 0 {
		if r.postLen == 0 {
			return 0
		}
		for i := 0; i < offset; i++ {
			r.data[r.preLen] = r.data[len(r.data)-r.postLen]
			r.preLen++
			r.postLen--
		}
	}
	if offset > 0 {
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
	if point > len(r.data)-r.gapLen() {
		point = len(r.data) - r.gapLen()
	}
	sp := point - 1
	p, err := r.RuneAt(sp)
	if p == '\n' {
		//fmt.Println("Newline", x)
		sp++
		return sp
	}
	for x := sp; x > 0; x-- {
		if x == 0 {
			return 0
		}
		p, err = r.RuneAt(x)
		if err != nil {
			panic(err)
		}
		if p == '\n' {
			//fmt.Println("Newline", x)
			x++
			return x
		}
	}
	return 0
}

// LineEnd find the point at end of this line
func (r *Buffer) LineEnd(point int) int {
	if point < 0 {
		return 0
	}
	ep := len(r.data) - r.gapLen()
	for {
		if point >= ep {
			return ep - 1
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

// LineLenAtPoint length of line at point
func (r *Buffer) LineLenAtPoint(point int) int {
	if point >= len(r.data) {
		point = len(r.data) - 1
	}
	if point < 0 {
		point = 0
	}
	start := r.LineStart(point) - 1
	end := r.LineEnd(point)
	return end - start
}

// PointForLine return point for beginning of line ln
func (r *Buffer) PointForLine(ln int) int {
	if ln <= 1 {
		return 0
	}
	lines := 0
	for pt := 0; pt < r.BufferLen(); pt++ {
		etch, err := r.RuneAt(pt)
		if err != nil {
			panic(err)
		}
		if etch == '\n' {
			lines++
		}
		if lines == ln {
			return r.LineStart(pt)
		}
	}
	return r.LineEnd(r.BufferLen() - 1)
}

// LineForPoint returns the line number of point (o = 1)
func (r *Buffer) LineForPoint(point int) (line int) {
	line = 1
	pt := 0
	if point >= r.BufferLen() {
		point = r.BufferLen() - 1
	}
	doIncr := false
	for pt = 1; pt <= point; pt++ {
		if doIncr {
			line++
			doIncr = false
		}
		etch, err := r.RuneAt(pt)
		if err != nil {
			panic(err)
		}
		if etch == '\n' {
			//line++
			doIncr = true
		}
	}
	// if pt == r.BufferLen() {
	// 	line--
	// }
	return
}

// ColumnForPoint returns the column (o = 1) of pt
func (r *Buffer) ColumnForPoint(point int) (column int) {
	if point >= r.BufferLen() {
		point = r.BufferLen() - 1
	}
	start := r.LineStart(point)
	return point - start + 1

}

// XYForPoint returns the cursor location for a pt in the buffer
func (r *Buffer) XYForPoint(pt int) (x, y int) {
	x = r.ColumnForPoint(pt)
	if r.EndOfBuffer(pt) {
		x = r.ColumnForPoint(r.LineEnd(pt))
	}
	y = r.LineForPoint(pt)
	return
}

// PointForXY returns the Point location for X, Y in the buffer
func (bp *Buffer) PointForXY(x, y int) (finalpt int) {
	//10, 1
	lpt := bp.PointForLine(y)
	c := x - 1
	finalpt = lpt + c //bp.DataPointForBufferPoint(lpt + c)
	return finalpt
}

// SegStart Forward scan for start of logical line segment
// (corresponds to screen line)  containing 'finish'
func (r *Buffer) SegStart(start, finish, limit int) int {
	//var p rune
	c := 0
	scan := start

	for scan < finish {
		//p = ptr(bp, scan);
		if scan >= r.BufferLen() {
			return r.BufferLen()
		}
		rch, err := r.RuneAt(scan)
		if err != nil {
			panic(err)
		}

		if rch == '\n' {
			c = 0
			start = scan + 1
		} else {
			if limit <= c {
				c = 0
				start = scan
			}
		}
		scan++
		if rch == '\t' {
			c += 4
		} else {
			c++
		}
	}
	if c < limit {
		return start
	}
	return finish
}

/* SegNext Forward scan for start of logical line segment following 'finish' */
func (r *Buffer) SegNext(start, finish, limit int) int {
	c := 0

	scan := r.SegStart(start, finish, limit)
	for {
		if scan >= r.BufferLen() {
			return r.BufferLen()
		}
		rch, err := r.RuneAt(scan)
		if err != nil {
			panic(err)
		}
		//if (bp.b_ebuf <= p || COLS <= c)
		if limit <= c {
			break
		}
		//scan += utf8_size(*ptr(bp,scan));
		scan++
		if rch == '\n' {
			break
		}
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if rch == '\t' {
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
	npt := l2 + c1 - 1
	if npt < r.PageStart {
		r.Reframe = true
	}
	r.SetPointAndCursor(npt)
}

// PointDown move point down one line
func (r *Buffer) PointDown() {
	c1 := r.ColumnForPoint(r.Point())
	l1 := r.LineEnd(r.Point())
	l2 := r.LineStart(l1 + 1)
	//fmt.Printf("PointDown c1 %d, l1 %d, l2 %d)\n", c1, l1, l2)
	npt := l2 + c1 - 1
	if npt > r.PageEnd {
		r.Reframe = true
	}
	r.SetPointAndCursor(npt)
	//fmt.Printf("Point %d (%d,%d)\n", r.Point(), r.PointCol, r.PointRow)
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

// UpUp Move up one screen line
func (bp *Buffer) UpUp(pt, cc int) int {
	//bp := e.CurrentBuffer
	curr := bp.LineStart(pt)
	seg := bp.SegStart(curr, pt, cc)
	if curr < seg {
		pt = bp.SegStart(curr, seg-1, cc)
	} else {
		pt = bp.SegStart(bp.LineStart(curr-1), curr-1, cc)
	}
	// x, y := bp.XYForPoint(pt)
	// if (y - 1) >= 1 {
	// 	pt = bp.PointForXY(x, y-1)
	// }
	return pt
}

// DownDown Move down one screen line
func (bp *Buffer) DownDown(pt, cc int) int {
	//bp := e.CurrentBuffer
	return bp.SegNext(bp.LineStart(pt), pt, cc)
	// x, y := bp.XYForPoint(pt)
	// npt := bp.PointForXY(x, y+1)
	// log.Printf("npt %d pt %d x %d y %d\n", npt, pt, x, y)
	// return npt
}

/* GetLineStats scan buffer and fill in curline and lastline */
func (r *Buffer) GetLineStats() (curline int, lastline int) {
	pt := r.Point()
	_, curline = r.XYForPoint(pt)
	_, lastline = r.XYForPoint(r.BufferLen())
	return curline, lastline
}

func (r *Buffer) DebugPrint() {
	fmt.Printf("*********(gap)\n")
	for i := 0; i < len(r.data); i++ {
		if i >= r.gapStart() && i < r.gapStart()+r.gapLen() {
			fmt.Printf("@")
		} else if i < r.preLen {
			if r.data[i] == '\n' {
				fmt.Printf("%c\n", 0x00B6)
			} else {
				fmt.Printf("%c", r.data[i])
			}
		} else {
			if r.data[i] == '\n' {
				fmt.Printf("%c\n", 0x00B6)
			} else {
				fmt.Printf("%c", r.data[i])
			}
		}
	}
	fmt.Printf("\n*********\n")
}
func (r *Buffer) DebugError(t *testing.T) {
	t.Error("|-")
	t.Errorf("*********\n")
	for i := 0; i < len(r.data); i++ {
		if i >= r.gapStart() && i < r.gapStart()+r.gapLen() {
			t.Errorf("@")
		} else if i < r.preLen {
			t.Errorf("%c", r.data[i])
		} else {
			t.Errorf("%c", r.data[i])
		}
	}

	t.Errorf("\n")
}
