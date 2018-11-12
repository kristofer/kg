package tkg

import (
	"fmt"
	"strings"
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
	Reframe    int    /* b_reframe force a reframe of the display */
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

var RootBuffer *Buffer = nil
var CurrentBuffer *Buffer = nil

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
	pt1 := r.IntForLine(l1)
	pt2 := r.IntForLine(l2)
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
	//fmt.Println(p, r.preLen, r.postLen, r.postStart())
	if p <= r.preLen && r.preLen != 0 {
		return r.data[p]
	}
	if p >= r.postLen || p <= len(r.data) {
		return r.data[r.postStart()+(p-r.preLen)]
	}
	return '\u2318'
}
func (r *Buffer) BufferLen() int {
	return r.preLen + r.postLen
}

func (r *Buffer) Len() int {
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

// IntForLine return point for line ln
func (r *Buffer) IntForLine(ln int) int {
	ep := len(r.data) - 1
	sp := 0
	for p := 0; p < ep; p++ {
		if p == r.preLen {
			p = r.postStart()
		}
		if r.data[p] == '\n' {
			ln--
			if ln == 0 {
				return sp
			}
			if (p + 1) < ep {
				sp = p + 1
			}
		}
	}
	return ep
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

/* Buffer lists manipulation */
/* Find a buffer by filename or create if requested */
func FindBuffer(fname string, cflag bool) *Buffer {
	var bp *Buffer
	var sb *Buffer

	bp = RootBuffer
	for bp != nil {
		if strings.Compare(fname, bp.Filename) == 0 || strings.Compare(fname, bp.Buffername) == 0 {
			return bp
		}
		bp = bp.Next
	}

	if cflag != false {
		// if ((bp = (buffer_t *) malloc (sizeof (buffer_t))) == nil)
		// 	return (0);
		bp = NewBuffer()

		//BufferInit(bp)
		//assert(bp != nil);

		/* find the place in the list to insert this buffer */
		if RootBuffer == nil {
			RootBuffer = bp
		} else if strings.Compare(RootBuffer.Filename, fname) > 0 {
			/* insert at the begining */
			bp.Next = RootBuffer
			RootBuffer = bp
		} else {
			for sb = RootBuffer; sb.Next != nil; sb = sb.Next {
				if strings.Compare(sb.Next.Filename, fname) > 0 {
					break
				}
			}
			/* and insert it */
			bp.Next = sb.Next
			sb.Next = bp
		}
	}
	return bp
}

// DeleteBuffer unlink from the list of buffers, free associated memory,
// assumes buffer has been saved if modified
func DeleteBuffer(bp *Buffer) bool {
	//editor := bp.CurrentWindow.Editor
	var sb *Buffer

	/* we must have switched to a different buffer first */
	//assert(bp != CurrentBuffer)
	if bp != CurrentBuffer {
		/* if buffer is the head buffer */
		if bp == RootBuffer {
			RootBuffer = bp.Next
		} else {
			/* find place where the bp buffer is next */
			for sb = RootBuffer; sb.Next != bp && sb.Next != nil; sb = sb.Next {
			}
			if sb.Next == bp || sb.Next == nil {
				sb.Next = bp.Next
			}
		}

		/* now we can delete */
		//free(bp.BufferStart);
		//bp.BufferStart = nil
		//free(bp);
		bp = nil
	} else {
		return false
	}
	return true
}

// NextBuffer returns next buffer after current
func NextBuffer(CurrentWindow *Window) {
	editor := CurrentWindow.Editor
	if editor.CurrentBuffer != nil && editor.RootBuffer != nil {
		CurrentWindow.DisassociateBuffer()
		if CurrentBuffer.Next != nil {
			CurrentBuffer = CurrentBuffer.Next

		} else {
			CurrentBuffer = RootBuffer
		}
		CurrentWindow.AssociateBuffer(CurrentBuffer)
	}
}

// GetBufferName returns buffer name
func GetBufferName(bp *Buffer) string {
	if bp.Filename != "" {
		return bp.Filename
	}
	return bp.Buffername
}

// CountBuffers how many buffers in list
func CountBuffers() int {
	var bp *Buffer
	i := 0

	for bp = RootBuffer; bp != nil; bp = bp.Next {
		i++
	}
	return i
}

// ModifiedBuffers true is any buffers modified
func ModifiedBuffers() bool {
	var bp *Buffer

	for bp = RootBuffer; bp != nil; bp = bp.Next {
		if bp.modified == true {
			return true
		}
	}
	return false
}
