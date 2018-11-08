package kg

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

// Create a new Buffer
func NewGapBuffer() *GapBuffer {
	return &GapBuffer{}
}

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

func (r *GapBuffer) gapStart() int {
	return r.preLen
}

func (r *GapBuffer) gapLen() int {
	return r.postStart() - r.preLen
}

func (r *GapBuffer) postStart() int {
	return len(r.data) - r.postLen
}

func (r *GapBuffer) Insert(s string) {
	if r.gapLen() < len(s) {
		newData := make([]rune, len(r.data)*2)

		copy(newData, r.data[:r.preLen])
		copy(newData[r.postStart()+len(r.data):],
			r.data[r.postStart():])

		r.data = newData
	}

	copy(r.data[r.gapStart():], []rune(s))
	r.preLen += len(s)
}

func (r *GapBuffer) Delete() {
	if r.postLen == 0 {
		return
	}

	r.postLen--
}

func (r *GapBuffer) Backspace() {
	if r.preLen == 0 {
		return
	}

	r.preLen--
}

func (r *GapBuffer) CursorNext() {
	if r.postLen == 0 {
		return
	}

	r.data[r.preLen] = r.data[r.postStart()]
	r.preLen++
	r.postLen--
}

func (r *GapBuffer) CursorPrevious() {
	if r.preLen == 0 {
		return
	}

	r.data[r.postStart()-1] = r.data[r.preLen-1]
	r.preLen--
	r.postLen++
}

func (r *GapBuffer) debugPrint() {
	for i := 0; i < len(r.data); i++ {
		if i >= r.gapStart() && i < r.gapStart()+r.gapLen() {
			fmt.Printf(" ")
		} else if i < r.preLen {
			preColor.Printf("%c", r.data[i])
		} else {
			postColor.Printf("%c", r.data[i])
		}
	}

	fmt.Printf("\n")
}


/* Enlarge gap by n chars, position of gap cannot change */
func (bp *Buffer)GrowGap(n Point) bool {
	//char_t *new;
	// var buflen, newlen, xgap, xegap Point
		
	// assert(bp->b_buf <= bp->b_gap);
	// assert(bp->b_gap <= bp->b_egap);
	// assert(bp->b_egap <= bp->b_ebuf);

	// xgap = bp->b_gap - bp->b_buf;
	// xegap = bp->b_egap - bp->b_buf;
	// buflen = bp->b_ebuf - bp->b_buf;
    
	// /* reduce number of reallocs by growing by a minimum amount */
	// n = (n < MIN_GAP_EXPAND ? MIN_GAP_EXPAND : n);
	// newlen = buflen + n * sizeof (char_t);

	// if (buflen == 0) {
	// 	if (newlen < 0 || MAX_SIZE_T < newlen)
	// 		fatal("%s: Failed to allocate required memory.\n");
	// 	new = (char_t*) malloc((size_t) newlen);
	// 	if (new == NULL)			
	// 		fatal("%s: Failed to allocate required memory.\n");	/* Cannot edit a file without a buffer. */
	// } else {
	// 	if (newlen < 0 || MAX_SIZE_T < newlen) {
	// 		msg("Failed to allocate required memory");
	// 		return (FALSE);
	// 	}
	// 	new = (char_t*) realloc(bp->b_buf, (size_t) newlen);
	// 	if (new == NULL) {
	// 		msg("Failed to allocate required memory");    /* Report non-fatal error. */
	// 		return (FALSE);
	// 	}
	// }

	// /* Relocate pointers in new buffer and append the new
	//  * extension to the end of the gap.
	//  */
	// bp->b_buf = new;
	// bp->b_gap = bp->b_buf + xgap;      
	// bp->b_ebuf = bp->b_buf + buflen;
	// bp->b_egap = bp->b_buf + newlen;
	// while (xegap < buflen--)
	// 	*--bp->b_egap = *--bp->b_ebuf;
	// bp->b_ebuf = bp->b_buf + newlen;

	// assert(bp->b_buf < bp->b_ebuf);          /* Buffer must exist. */
	// assert(bp->b_buf <= bp->b_gap);
	// assert(bp->b_gap < bp->b_egap);          /* Gap must grow only. */
	// assert(bp->b_egap <= bp->b_ebuf);
	// return (TRUE);
	return false
}

//point_t movegap(bp *Buffer, point_t offset)
func (bp *Buffer) MoveGap(offset Point) Point {

	// char_t *p = ptr(bp, offset);
	// while (p < bp->b_gap)
	// 	*--bp->b_egap = *--bp->b_gap;
	// while (bp->b_egap < p)
	// 	*bp->b_gap++ = *bp->b_egap++;
	// assert(bp->b_gap <= bp->b_egap);
	// assert(bp->b_buf <= bp->b_gap);
	// assert(bp->b_egap <= bp->b_ebuf);
	// return (pos(bp, bp->b_egap));
	return 0
}

/* Given a buffer offset, convert it to a pointer into the buffer */
//char_t * ptr(bp *Buffer, register point_t offset)
func (bp *Buffer) Ptr(offset Point) Point {
	if (offset < 0)
		return (bp->b_buf);
	return (bp->b_buf+offset + (bp->b_buf + offset < bp->b_gap ? 0 : bp->b_egap-bp->b_gap));
}

/* Given a pointer into the buffer, convert it to a buffer offset */
//point_t pos(bp *Buffer, register char_t *cp)
func (bp *Buffer) Pos(cp Point) Point {
	assert(bp->b_buf <= cp && cp <= bp->b_ebuf);
	return (cp - bp->b_buf - (cp < bp->b_egap ? 0 : bp->b_egap - bp->b_gap));
}

/* find the point for start of line ln */
func (bp *Buffer) LineToPoint(ln int) Point {
	// point_t end_p = pos(curbp, curbp->b_ebuf);
	// point_t p, start;

	// for (p=0, start=0; p < end_p; p++) {
	// 	if ( *(ptr(curbp, p)) == '\n') {
	// 		if (--ln == 0)
	// 			return start;
	// 		if (p + 1 < end_p) 
	// 			start = p + 1;
	// 	}
	// }
	return 0
}

/* scan buffer and fill in curline and lastline */
func (bp *Buffer) GetLineStats(curline int, lastline int) {
	// point_t end_p = pos(curbp, curbp->b_ebuf);
	// point_t p;
	// int line;
    
	// *curline = -1;
    
	// for (p=0, line=0; p < end_p; p++) {
	// 	line += (*(ptr(curbp,p)) == '\n') ? 1 : 0;
	// 	*lastline = line;
        
	// 	if (*curline == -1 && p == curbp->b_point) {
	// 		*curline = (*(ptr(curbp,p)) == '\n') ? line : line + 1;
	// 	}
	// }

	// *lastline = *lastline + 1;
	
	// if (curbp->b_point == end_p)
	// 	*curline = *lastline;
}