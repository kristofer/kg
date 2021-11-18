package kg

import (
	"errors"
	"fmt"
)

/*
 * Buffer is where all the operations on the main rune array are implemented.
 * Because of the Gap, all the indexing around it should be done by these routines.
 */

// Buffer main struct
type Buffer struct {
	Point   int
	postLen int
	data    []rune

	Next       *Buffer
	Mark       int
	OrigPoint  int    /* b_cpoint the original current point, used for multiple window displaying */
	PageStart  int    /*  start of page */
	PageEnd    int    /*  end of page */
	Reframe    bool   /*  force a reframe of the display */
	WinCount   int    /* b_nt count of windows referencing this buffer */
	TextSize   int    /*  current size of text being edited (not including gap) */
	PrevSize   int    /* b_psize previous size */
	PointRow   int    /* b_row Point row */
	PointCol   int    /* b_col Point col */
	Filename   string // b_fname[NAME_MAX + 1]; /* filename */
	Buffername string //[b_bnameSTRBUF_S];   /* buffer name */
	Flags      byte   /* char b_flags buffer flags */
	modified   bool
}

// MarkModified xxx
func (bp *Buffer) MarkModified() {
	bp.modified = true
	bp.TextSize = bp.Point + bp.postLen
}

// NewBuffer - Create a new Buffer
func NewBuffer() *Buffer {
	nb := Buffer{}
	nb.setText("\n")
	return &nb
}

// setText xxx
func (bp *Buffer) setText(s string) {
	bp.data = []rune(s)
	bp.Point = 0
	bp.postLen = len(bp.data)
	bp.TextSize = bp.Point + bp.postLen
}

// getText  xxx
func (bp *Buffer) getText() string {
	//bp.TextSize = bp.Point + bp.postLen
	ret := make([]rune, bp.Point+bp.postLen)
	copy(ret, bp.data)
	copy(ret[bp.Point:], bp.data[bp.postStart():])
	return string(ret)
}

// RuneAt finally reliable!! (well, maybe not)
func (bp *Buffer) RuneAt(pt int) (rune, error) {
	//log.Println("RuneAt pt = ", pt)
	if pt >= len(bp.data) {
		return 0, errors.New("beyond data buffer in RuneAt")
	}
	if pt < 0 {
		//return '\u0000', errors.New("negative buffer pointer in RuneAt")
		pt = 0
	}
	if npt := bp.dataPointForBufferPoint(pt); npt < len(bp.data) {
		return bp.data[npt], nil
	}
	return 0, errors.New("ran over end of data buffer in RuneAt")
}

func (bp *Buffer) dataPointForBufferPoint(pt int) int {
	npt := 0
	if pt < bp.Point {
		npt = pt
	}
	if pt >= bp.Point && pt < len(bp.data) {
		npt = pt + bp.gapLen()
	}
	return npt
}

// AddRune add a run to the buffer
func (bp *Buffer) AddRune(ch rune) {
	if bp.gapLen() == 0 {
		_ = bp.GrowGap(gapchunk)
	}
	bp.data[bp.Point] = ch
	bp.Point++
	bp.MarkModified()
}

// SetPoint set the current point to np
func (bp *Buffer) SetPoint(np int) {
	bp.CollapseGap()
	//bp.MoveGap(np - bp.Point)
	// move gap <-(left) by np chars
	gs := bp.gapStart()
	//log.Printf("gap start %d len %d new pt %d dist %d\n", gs, bp.gapLen(), np, gs-np)
	f := 0
	for i := gs - np; i > 0; i-- {
		bp.data[bp.postStart()-1] = bp.data[bp.Point-1]
		bp.Point--
		bp.postLen++
		f++
	}
	//log.Printf("shuffled %d\n", f)

	if bp.PageEnd < bp.Point {
		bp.Reframe = true
	}
}

// setCursor xxx
func (bp *Buffer) setCursor() {
	x, y := bp.XYForPoint(bp.Point)
	bp.PointRow = y
	bp.PointCol = x
}

func (bp *Buffer) gapStart() int {
	return bp.Point
}

func (bp *Buffer) gapLen() int {
	return bp.postStart() - bp.Point
}

func (bp *Buffer) postStart() int {
	return len(bp.data) - bp.postLen
}

// CollapseGap moves the gap to the end of the buffer
func (bp *Buffer) CollapseGap() {
	//for i := bp.Point; bp.postLen > 0; i++ {
	for bp.postLen > 0 {
		bp.data[bp.Point] = bp.data[len(bp.data)-bp.postLen]
		bp.Point++
		bp.postLen--
	}
}

// Insert adds the string, growing the gap if needed.
func (bp *Buffer) Insert(s string) {
	if bp.gapLen() < len(s) {
		newGap := len(s) + 32
		_ = bp.GrowGap(newGap)
	}
	copy(bp.data[bp.gapStart():], []rune(s))
	bp.Point += len(s)
	bp.MarkModified()
}

// getTextForLines return string for [l1, l2) (l2 not included)
func (bp *Buffer) getTextForLines(l1, l2 int) string {
	pt1 := bp.PointForLine(l1)
	ret := make([]rune, bp.PointForLine(l2)-pt1)
	for i, j := pt1, 0; j < len(ret); i++ {
		rch, err := bp.RuneAt(i)
		checkErr(err)
		ret[j] = rch
		j++
	}
	return string(ret)
}

// GrowGap makes the gap bigger by n
// not sure why I need this.
func (bp *Buffer) GrowGap(n int) bool {
	newData := make([]rune, len(bp.data)+n)
	copy(newData, bp.data[:bp.Point])
	copy(newData[bp.postStart()+n:],
		bp.data[bp.postStart():])
	bp.data = newData
	bp.TextSize = bp.Point + bp.postLen
	return true
}

// MoveGap moves the gap to a Point
func (bp *Buffer) MoveGap(offset int) int {

	if offset < 0 {
		if bp.postLen == 0 {
			return 0
		}
		for i := 0; i < offset; i++ {
			bp.data[bp.Point] = bp.data[len(bp.data)-bp.postLen]
			bp.Point++
			bp.postLen--
		}
	}
	if offset > 0 {
		if bp.Point == 0 {
			return 0
		}
		for i := offset; i < 0; i++ {
			bp.data[bp.postStart()-1] = bp.data[bp.Point-1]
			bp.Point--
			bp.postLen++
		}
	}

	return offset
}

// Remove extent runes starting at from point
func (bp *Buffer) Remove(from int, extent int) {
	bp.SetPoint(from)
	for k := 0; k < extent; k++ {
		bp.Delete()
	}
}

//LineStart xxx
func (bp *Buffer) LineStart(point int) int {
	if point > len(bp.data)-bp.gapLen() {
		point = len(bp.data) - bp.gapLen()
	}
	sp := point - 1
	p, err := bp.RuneAt(sp)
	checkErr(err)
	if p == '\n' {
		sp++
		return sp
	}
	for x := sp; x > 0; x-- {
		if x == 0 {
			return 0
		}
		p, err = bp.RuneAt(x)
		checkErr(err)
		if p == '\n' {
			x++
			return x
		}
	}
	return 0
}

// LineEnd find the point at end of this line
func (bp *Buffer) LineEnd(point int) int {
	if point < 0 {
		return 0
	}
	ep := len(bp.data) - bp.gapLen()
	for {
		if point >= ep {
			return ep - 1
		}
		p, err := bp.RuneAt(point)
		checkErr(err)
		if p == '\n' {
			return point
		}
		point++
	}
}

// LineLenAtPoint length of line at point
func (bp *Buffer) LineLenAtPoint(point int) int {
	if point >= len(bp.data) {
		point = len(bp.data) - 1
	}
	if point < 0 {
		point = 0
	}
	start := bp.LineStart(point) - 1
	end := bp.LineEnd(point)
	return end - start
}

// PointForLine return point for beginning of line ln
func (bp *Buffer) PointForLine(ln int) int {
	if ln <= 1 {
		return 0
	}
	lines := 0
	for pt := 0; pt < bp.TextSize; pt++ {
		etch, err := bp.RuneAt(pt)
		checkErr(err)
		if etch == '\n' {
			lines++
		}
		if lines == ln {
			return bp.LineStart(pt)
		}
	}
	return bp.LineEnd(bp.TextSize) // -1
}

// LineForPoint returns the line number of point (origin = 1)
func (bp *Buffer) LineForPoint(point int) (line int) {
	line = 1
	pt := 0
	if point >= bp.TextSize {
		point = bp.TextSize - 1
	}
	doIncr := false
	for pt = 0; pt <= point; pt++ {
		if doIncr {
			line++
			doIncr = false
		}
		etch, err := bp.RuneAt(pt)
		checkErr(err)
		if etch == '\n' {
			doIncr = true
		}
	}
	return
}

// ColumnForPoint returns the column (origin = 1) of pt
func (bp *Buffer) ColumnForPoint(point int) (column int) {
	if point >= bp.TextSize {
		point = bp.TextSize - 1
	}
	return point - bp.LineStart(point) + 1
}

// XYForPoint returns the cursor location for a pt in the buffer
func (bp *Buffer) XYForPoint(pt int) (x, y int) {
	x = bp.ColumnForPoint(pt)
	if bp.TextSize = bp.Point + bp.postLen; pt >= bp.TextSize {
		x = bp.ColumnForPoint(bp.LineEnd(pt))
	}
	y = bp.LineForPoint(pt)
	return
}

// PointForXY returns the Point location for X, Y in the buffer
func (bp *Buffer) PointForXY(x, y int) (finalpt int) {
	lpt := bp.PointForLine(y)
	finalpt = lpt + x - 1
	if finalpt < 0 {
		finalpt = 0
	}
	return finalpt
}

// SegStart Forward scan for start of logical line segment
// (corresponds to screen line)  containing 'finish'
func (bp *Buffer) SegStart(start, finish, limit int) int {
	//var p rune
	c := 0
	scan := start

	for scan < finish {
		//p = ptr(bp, scan);
		if scan >= bp.TextSize {
			return bp.TextSize
		}
		rch, err := bp.RuneAt(scan)
		checkErr(err)

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

// SegNext Forward scan for start of logical line segment following 'finish'
func (bp *Buffer) SegNext(start, finish, limit int) int {
	c := 0

	scan := bp.SegStart(start, finish, limit)
	for {
		if scan >= bp.TextSize {
			return bp.TextSize
		}
		rch, err := bp.RuneAt(scan)
		checkErr(err)
		if limit <= c {
			break
		}
		scan++
		if rch == '\n' {
			break
		}
		if rch == '\t' {
			c += 4 //8 - (c % 7)
		} else {
			c++
		}
	}
	if scan < bp.TextSize {
		return scan
	}
	return bp.TextSize
}

// Delete remove a rune forward
func (bp *Buffer) Delete() {
	if bp.postLen == 0 {
		return
	}
	bp.postLen--
	bp.MarkModified()
}

// Backspace remove a rune backward
func (bp *Buffer) Backspace() {
	if bp.Point == 0 {
		return
	}
	bp.Point--
	bp.MarkModified()
}

// PointUp move point up one line
func (bp *Buffer) PointUp() {
	c1 := bp.ColumnForPoint(bp.Point)
	l1 := bp.LineStart(bp.Point)
	l2 := bp.LineStart(l1 - 1)
	l2l := bp.LineLenAtPoint(l2)
	npt := l2 + c1 - 1
	if l2l < c1 {
		npt = l2 + l2l - 1
	}
	if npt < bp.PageStart {
		bp.Reframe = true
	}
	bp.SetPoint(npt)
	bp.setCursor()
}

// PointDown move point down one line
func (bp *Buffer) PointDown() {
	c1 := bp.ColumnForPoint(bp.Point)
	l1 := bp.LineEnd(bp.Point)
	l2 := bp.LineStart(l1 + 1)
	l2l := bp.LineLenAtPoint(l2)
	npt := l2 + c1 - 1
	if l2l < c1 {
		npt = l2 + l2l - 1
	}
	if npt > bp.PageEnd {
		bp.Reframe = true
	}
	bp.SetPoint(npt)
	bp.setCursor()
}

// PointNext move point left one
func (bp *Buffer) PointNext() {
	// this is from the END OF BUFFER nonsense I had to fix.
	if bp.postLen <= 1 { //== 0 {
		return
	}
	bp.data[bp.Point] = bp.data[bp.postStart()]
	bp.Point++
	bp.postLen--
}

// PointPrevious move point right one
func (bp *Buffer) PointPrevious() {
	if bp.Point == 0 {
		return
	}

	bp.data[bp.postStart()-1] = bp.data[bp.Point-1]
	bp.Point--
	bp.postLen++
}

// UpUp Move up one screen line
func (bp *Buffer) UpUp(pt, cc int) int {
	curr := bp.LineStart(pt)
	seg := bp.SegStart(curr, pt, cc)
	if curr < seg {
		pt = bp.SegStart(curr, seg-1, cc)
	} else {
		pt = bp.SegStart(bp.LineStart(curr-1), curr-1, cc)
	}
	return pt
}

// DownDown Move down one screen line
func (bp *Buffer) DownDown(pt, cc int) int {
	return bp.SegNext(bp.LineStart(pt), pt, cc)
}

// GetLineStats scan buffer and fill in curline and lastline
func (bp *Buffer) GetLineStats() (curline int, lastline int) {
	pt := bp.Point
	_, curline = bp.XYForPoint(pt)
	_, lastline = bp.XYForPoint(bp.TextSize)
	return curline, lastline
}

func (bp *Buffer) gotoLine(ln int) {
	pt := bp.PointForLine(ln)
	bp.SetPoint(pt)
}

// DebugPrint prints out a view of the buffer and the gap and so on.
func (bp *Buffer) DebugPrint() {
	fmt.Printf("*********(gap)\n")
	for i := 0; i < len(bp.data); i++ {
		if i >= bp.gapStart() && i < bp.gapStart()+bp.gapLen() {
			fmt.Printf("@")
		} else if i < bp.Point {
			if bp.data[i] == '\n' {
				fmt.Printf("%c\n", 0x00B6)
			} else {
				fmt.Printf("%c", bp.data[i])
			}
		} else {
			if bp.data[i] == '\n' {
				fmt.Printf("%c\n", 0x00B6)
			} else {
				fmt.Printf("%c", bp.data[i])
			}
		}
	}
	fmt.Printf("\n*********\n")
}
