package gapbuffer

import (
	"fmt"
)

/*
 * GapBuffer
 */

type GapBuffer struct {
	data    []rune
	preLen  int
	postLen int
}

// NewGapBuffer - Create a new Buffer
func NewGapBuffer() *GapBuffer {
	return &GapBuffer{}
}

// func (r *GapBuffer) RuneAt(pt int) *rune {
// 	return &(r.data[pt])
// }

func (r *GapBuffer) SetText(s string) {
	r.data = []rune(s)
	r.preLen = 0
	r.postLen = len(r.data)
}

func (r *GapBuffer) GetText() string {
	ret := make([]rune, r.preLen+r.postLen)

	copy(ret, r.data)
	copy(ret[r.preLen:], r.data[r.postStart():])

	return string(ret)
}

func (r *GapBuffer) BufferLen() int {
	return r.preLen + r.postLen
}

func (r *GapBuffer) Len() int {
	return len(r.data)
}

func (r *GapBuffer) gapStart() int {
	return r.preLen
}

func (r *GapBuffer) gapLen() int {
	return r.postStart() - r.preLen
}

func (r *GapBuffer) postStart() int {
	return len(r.data) - r.postLen
}

// Insert adds the string, growing the gap if needed.
func (r *GapBuffer) Insert(s string) {
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
func (r *GapBuffer) GrowGap(n int) bool {
	newData := make([]rune, len(r.data)+n)

	copy(newData, r.data[:r.preLen])

	copy(newData[r.postStart()+n:],
		r.data[r.postStart():])

	r.data = newData
	return true
}

// MoveGap moves the gap forward by offset runes
func (r *GapBuffer) MoveGap(offset int) int {

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
func (r *GapBuffer) IntForLine(ln int) int {
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
func (r *GapBuffer) Delete() {
	if r.postLen == 0 {
		return
	}

	r.postLen--
}

// Backspace remove a rune backward
func (r *GapBuffer) Backspace() {
	if r.preLen == 0 {
		return
	}

	r.preLen--
}

// Cursor return point
func (r *GapBuffer) Cursor() int {
	return r.preLen
}

// PrintCursor print cursor point
func (r *GapBuffer) PrintCursor() {
	fmt.Println("C: ", r.Cursor())
}

// CursorNext move point left one
func (r *GapBuffer) CursorNext() {
	if r.postLen == 0 {
		return
	}

	r.data[r.preLen] = r.data[r.postStart()]
	r.preLen++
	r.postLen--
}

// CursorPrevious move point right one
func (r *GapBuffer) CursorPrevious() {
	if r.preLen == 0 {
		return
	}

	r.data[r.postStart()-1] = r.data[r.preLen-1]
	r.preLen--
	r.postLen++
}

/* GetLineStats scan buffer and fill in curline and lastline */
func (r *GapBuffer) GetLineStats() (curline int, lastline int) {
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

func (r *GapBuffer) debugPrint() {
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
