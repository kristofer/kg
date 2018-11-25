package kg

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

// Buffer main struct
type Buffer struct {
	data    []rune
	preLen  int
	postLen int
	Next    *Buffer /* b_next Link to next buffer_t */
	Mark    int     /* b_mark the mark */
	//Point      int     /* b_point the point */
	OrigPoint int /* b_cpoint the original current point, used for mutliple window displaying */
	PageStart int /* b_page start of page */
	PageEnd   int /* b_epage end of page */
	// FirstLine  int    /* b_page start of page */
	// LastLine   int    /* b_epage end of page */
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

// MarkModified xxx
func (bp *Buffer) MarkModified() {
	bp.modified = true
}

// NewBuffer - Create a new Buffer
func NewBuffer() *Buffer {
	nb := Buffer{}
	nb.data = []rune("\n")
	nb.preLen = 0
	nb.postLen = len(nb.data)
	return &nb
}

// SetText xxx
func (bp *Buffer) SetText(s string) {
	bp.data = []rune(s)
	bp.preLen = 0
	bp.postLen = len(bp.data)
}

// GetText  xxx
func (bp *Buffer) GetText() string {
	ret := make([]rune, bp.preLen+bp.postLen)
	copy(ret, bp.data)
	copy(ret[bp.preLen:], bp.data[bp.postStart():])
	return string(ret)
}

func (bp *Buffer) logBufferEOB(pt int) {
	if bp.EndOfBuffer(pt) == true {
		pc, file, no, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		if ok && details != nil {
			log.Printf(">>Called from %s\n>> %s Ln# %d\n", details.Name(), file, no)
		}
		log.Println(">>Setting Point to EOB", pt, bp.BufferLen())
	}
}

// RuneAt finally reliable!!
func (bp *Buffer) RuneAt(pt int) (rune, error) {
	bp.logBufferEOB(pt)
	if pt >= len(bp.data) {
		return 0, errors.New("Beyond data buffer in RuneAt")
	}
	if pt < 0 {
		return '\u0000', errors.New("negative buffer pointer in RuneAt")
	}
	if npt := bp.dataPointForBufferPoint(pt); npt < len(bp.data) {
		return bp.data[npt], nil
	}
	return 0, errors.New("Ran over end of data buffer in RuneAt")
}

func (bp *Buffer) dataPointForBufferPoint(pt int) int {
	npt := 0
	if pt < bp.preLen {
		npt = pt
	}
	if pt >= bp.preLen && pt < len(bp.data) {
		npt = pt + bp.gapLen()
	}
	return npt
}

// AddRune add a run to the buffer
func (bp *Buffer) AddRune(ch rune) {
	// if bp.data == nil {
	// 	bp.SetText(string(ch))
	// 	return
	// }
	if bp.gapLen() == 0 {
		_ = bp.GrowGap(gapchunk)
	}
	bp.data[bp.preLen] = ch
	bp.preLen++
	bp.MarkModified()
}

// Point return point
func (bp *Buffer) Point() int {
	return bp.preLen
}

// SetPoint set the current point to np
func (bp *Buffer) SetPoint(np int) {
	bp.logBufferEOB(np)
	bp.CollapseGap()
	// move gap <-(left) by np chars
	gs := bp.gapStart()
	for i := gs - np; i > 0; i-- {
		bp.data[bp.postStart()-1] = bp.data[bp.preLen-1]
		bp.preLen--
		bp.postLen++
	}
	if bp.PageEnd < bp.preLen {
		//log.Println("reframing!")
		bp.Reframe = true
	}
}

//SetPointAndCursor xxx
func (bp *Buffer) SetPointAndCursor(np int) {
	bp.SetPoint(np)
	bp.setCursor()
}

// setCursor xxx
func (bp *Buffer) setCursor() {
	x, y := bp.XYForPoint(bp.preLen)
	bp.PointRow = y
	bp.PointCol = x
}

// PrintPoint print Point point
func (bp *Buffer) PrintPoint() {
	fmt.Println("C: ", bp.Point())
}

// BufferLen length of buffer
func (bp *Buffer) BufferLen() int {
	return bp.preLen + bp.postLen
}

// EndOfBuffer xxx
func (bp *Buffer) EndOfBuffer(pt int) bool {
	return pt >= (bp.preLen + bp.postLen)
}

// ActualLen length of buffer plus gap
func (bp *Buffer) ActualLen() int {
	return len(bp.data)
}

func (bp *Buffer) gapStart() int {
	return bp.preLen
}

func (bp *Buffer) gapLen() int {
	return bp.postStart() - bp.preLen
}

func (bp *Buffer) postStart() int {
	return len(bp.data) - bp.postLen
}

// CollapseGap moves the gap to the end of the buffer for replacement
func (bp *Buffer) CollapseGap() {
	for i := bp.preLen; bp.postLen > 0; i++ {
		bp.data[bp.preLen] = bp.data[len(bp.data)-bp.postLen]
		bp.preLen++
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
	bp.preLen += len(s)
	//fmt.Println("G", len(bp.data)-bp.postLen-bp.preLen, "S", bp.preLen, "E", bp.postLen)
	bp.MarkModified()
}

// GetTextForLines return string for [l1, l2) (l2 not included)
func (bp *Buffer) GetTextForLines(l1, l2 int) string {
	pt1 := bp.PointForLine(l1)
	pt2 := bp.PointForLine(l2)
	//fmt.Println(pt1, pt2)
	ret := make([]rune, pt2-pt1)
	j := 0
	for i := pt1; j < len(ret); i++ {
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

	copy(newData, bp.data[:bp.preLen])

	copy(newData[bp.postStart()+n:],
		bp.data[bp.postStart():])

	bp.data = newData
	return true
}

// MoveGap moves the gap to a Point
func (bp *Buffer) MoveGap(offset int) int {

	if offset < 0 {
		if bp.postLen == 0 {
			return 0
		}
		for i := 0; i < offset; i++ {
			bp.data[bp.preLen] = bp.data[len(bp.data)-bp.postLen]
			bp.preLen++
			bp.postLen--
		}
	}
	if offset > 0 {
		if bp.preLen == 0 {
			return 0
		}
		for i := offset; i < 0; i++ {
			bp.data[bp.postStart()-1] = bp.data[bp.preLen-1]
			bp.preLen--
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
	if p == '\n' {
		//fmt.Println("Newline", x)
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
			//fmt.Println("Newline", x)
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
	for pt := 0; pt < bp.BufferLen(); pt++ {
		etch, err := bp.RuneAt(pt)
		checkErr(err)
		if etch == '\n' {
			lines++
		}
		if lines == ln {
			return bp.LineStart(pt)
		}
	}
	return bp.LineEnd(bp.BufferLen()) // -1
}

// LineForPoint returns the line number of point (o = 1)
func (bp *Buffer) LineForPoint(point int) (line int) {
	line = 1
	pt := 0
	if point >= bp.BufferLen() {
		point = bp.BufferLen() - 1
	}
	doIncr := false
	for pt = 1; pt <= point; pt++ {
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

// ColumnForPoint returns the column (o = 1) of pt
func (bp *Buffer) ColumnForPoint(point int) (column int) {
	if point >= bp.BufferLen() {
		point = bp.BufferLen() - 1
	}
	start := bp.LineStart(point)
	return point - start + 1

}

// XYForPoint returns the cursor location for a pt in the buffer
func (bp *Buffer) XYForPoint(pt int) (x, y int) {
	x = bp.ColumnForPoint(pt)
	if bp.EndOfBuffer(pt) {
		x = bp.ColumnForPoint(bp.LineEnd(pt))
	}
	y = bp.LineForPoint(pt)
	return
}

// PointForXY returns the Point location for X, Y in the buffer
func (bp *Buffer) PointForXY(x, y int) (finalpt int) {
	//10, 1
	lpt := bp.PointForLine(y)
	c := x - 1
	finalpt = lpt + c //bp.DataPointForBufferPoint(lpt + c)
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
		if scan >= bp.BufferLen() {
			return bp.BufferLen()
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
		if scan >= bp.BufferLen() {
			return bp.BufferLen()
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
	if scan < bp.BufferLen() {
		return scan
	}
	return bp.BufferLen()
}

// Delete remove a rune forward
func (bp *Buffer) Delete() {
	if bp.postLen == 0 {
		return
	}

	bp.postLen--
}

// Backspace remove a rune backward
func (bp *Buffer) Backspace() {
	if bp.preLen == 0 {
		return
	}

	bp.preLen--
}

// PointUp move point up one line
func (bp *Buffer) PointUp() {
	c1 := bp.ColumnForPoint(bp.Point())
	l1 := bp.LineStart(bp.Point())
	l2 := bp.LineStart(l1 - 1)
	l2l := bp.LineLenAtPoint(l2)
	npt := l2 + c1 - 1
	if l2l < c1 {
		npt = l2 + l2l - 1
	}
	if npt < bp.PageStart {
		bp.Reframe = true
	}
	bp.SetPointAndCursor(npt)
}

// PointDown move point down one line
func (bp *Buffer) PointDown() {
	c1 := bp.ColumnForPoint(bp.Point())
	l1 := bp.LineEnd(bp.Point())
	l2 := bp.LineStart(l1 + 1)
	l2l := bp.LineLenAtPoint(l2)
	npt := l2 + c1 - 1
	if l2l < c1 {
		npt = l2 + l2l - 1
	}
	if npt > bp.PageEnd {
		bp.Reframe = true
	}
	bp.SetPointAndCursor(npt)
}

// PointNext move point left one
func (bp *Buffer) PointNext() {
	// this is from the END OF BUFFER nonsense I had to fix.
	if bp.postLen <= 1 { //== 0 {
		return
	}
	bp.data[bp.preLen] = bp.data[bp.postStart()]
	bp.preLen++
	bp.postLen--
}

// PointPrevious move point right one
func (bp *Buffer) PointPrevious() {
	if bp.preLen == 0 {
		return
	}

	bp.data[bp.postStart()-1] = bp.data[bp.preLen-1]
	bp.preLen--
	bp.postLen++
	//bp.setCursor()
	bp.logBufferEOB(bp.preLen)
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
	//bp := e.CurrentBuffer
	return bp.SegNext(bp.LineStart(pt), pt, cc)
}

// GetLineStats scan buffer and fill in curline and lastline
func (bp *Buffer) GetLineStats() (curline int, lastline int) {
	pt := bp.Point()
	_, curline = bp.XYForPoint(pt)
	_, lastline = bp.XYForPoint(bp.BufferLen())
	return curline, lastline
}

// DebugPrint xxx
func (bp *Buffer) DebugPrint() {
	fmt.Printf("*********(gap)\n")
	for i := 0; i < len(bp.data); i++ {
		if i >= bp.gapStart() && i < bp.gapStart()+bp.gapLen() {
			fmt.Printf("@")
		} else if i < bp.preLen {
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
