package kg

import "strings"

/*
* Buffer
 */

type Buffer struct {
	Next      *Buffer /* b_next Link to next buffer_t */
	Mark      int     /* b_mark the mark */
	Point     int     /* b_point the point */
	OrigPoint int     /* b_cpoint the original current point, used for mutliple window displaying */
	PageStart int     /* b_page start of page */
	PageEnd   int     /* b_epage end of page */
	Reframe   int     /* b_reframe force a reframe of the display */
	WinCount  int     /* b_cnt count of windows referencing this buffer */
	TextSize  int     /* b_size current size of text being edited (not including gap) */
	PrevSize  int     /* b_psize previous size */
	// BufferStart *string /* b_buf start of buffer */
	// BufferEnd   *string /* b_ebuf end of buffer */
	// GapStart    *string /* b_gap start of gap */
	// GapEnd      *string /* b_egap end of gap */
	Buffer     *GapBuffer /* actual buffer*/
	CursorRow  int        /* b_row cursor row */
	CursorCol  int        /* b_col cursor col */
	Filename   string     // b_fname[NAME_MAX + 1]; /* filename */
	Buffername string     //[b_bnameSTRBUF_S];   /* buffer name */
	Flags      byte       /* char b_flags buffer flags */
}

func BufferInit(bp *Buffer) {
	bp.Buffer = NewGapBuffer()
	bp.Mark = NOMARK
	bp.Point = 0
	bp.OrigPoint = 0
	bp.PageStart = 0
	bp.PageEnd = 0
	bp.Reframe = 0
	bp.TextSize = 0
	bp.PrevSize = 0
	bp.Flags = 0
	bp.WinCount = 0
	// bp.BufferStart = nil
	// bp.BufferEnd = nil
	// bp.GapStart = nil
	// bp.GapEnd = nil
	bp.Next = nil
	bp.Filename = ""
}

func (bp *Buffer) GetCurrentRune() (*rune, int) {
	return bp.Buffer.GetCurrentRune(), bp.Point
}
func (bp *Buffer) GetCurrentRune(arb int) (*rune, int) {
	return &(bp.Buffer[arb]), arb
}
func (bp *Buffer) EndOfBuffer() int {
	return bp.TextSize - 1
}

/* Find a buffer by filename or create if requested */
func FindBuffer(fname string, cflag bool) *Buffer {
	var bp *Buffer
	var sb *Buffer

	bp = Bheadp
	for bp != nil {
		if strings.Compare(fname, bp.Filename) == 0 || strings.Compare(fname, bp.Buffername) == 0 {
			return bp
		}
		bp = bp.Next
	}

	if cflag != false {
		// if ((bp = (buffer_t *) malloc (sizeof (buffer_t))) == nil)
		// 	return (0);
		bp = make(Buffer())

		BufferInit(bp)
		//assert(bp != nil);

		/* find the place in the list to insert this buffer */
		if Bheadp == nil {
			Bheadp = bp
		} else if string.Compare(Bheadp.filename, fname) > 0 {
			/* insert at the begining */
			bp.Next = Bheadp
			Bheadp = bp
		} else {
			for sb = Bheadp; sb.Next != nil; sb = sb.Next {
				if string.Compare(sb.Next.filename, fname) > 0 {
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

/* unlink from the list of buffers, free associated memory, assumes buffer has been saved if modified */
func DeleteBuffer(bp *Buffer) bool {
	var sb *Buffer

	/* we must have switched to a different buffer first */
	//assert(bp != Curbp)
	if bp != Curbp {
		/* if buffer is the head buffer */
		if bp == Bheadp {
			Bheadp = bp.Next
		} else {
			/* find place where the bp buffer is next */
			for sb = Bheadp; sb.Next != bp && sb.Next != nil; sb = sb.Next {
			}
			if sb.Next == bp || sb.Next == nil {
				sb.Next = bp.Next
			}
		}

		/* now we can delete */
		//free(bp.BufferStart);
		bp.BufferStart = nil
		//free(bp);
		bp = nil
	} else {
		return false
	}
	return true
}

func NextBuffer() {
	// assert(Curbp != nil);
	// assert(Bheadp != nil);
	if Curbp != nil && Bheadp != nil {
		disassociate_b(Curwp)
		//Curbp = (Curbp.Next != nil ? Curbp.Next : Bheadp);
		if Curbp.Next != nil {
			Curbp = Curbp.Next

		} else {
			Curbp = Bheadp
		}
		associate_b2w(Curbp, Curwp)
	}
}

func GetBufferName(bp *Buffer) string {
	if bp.filename != nil && bp.filename != "" {
		return bp.filename
	}
	return bp.Buffername
}

func CountBuffers() int {
	var bp *Buffer
	i := 0

	for bp = Bheadp; bp != nil; bp = bp.Next {
		i++
	}
	return i
}

func ModifiedBuffers() bool {
	var bp *Buffer

	for bp = Bheadp; bp != nil; bp = bp.Next {
		if bp.Flags & B_MODIFIED {
			return true
		}
	}
	return false
}

/* Enlarge gap by n chars, position of gap cannot change */
func (bp *Buffer) GrowGap(n int) bool {
	//char_t *new;
	// var buflen, newlen, xgap, xegap int

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

	// /* Relocate inters in new buffer and append the new
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

//int_t movegap(bp *Buffer, int_t offset)
func (bp *Buffer) MoveGap(offset int) int {

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

/* Given a buffer offset, convert it to a inter into the buffer */
//char_t * ptr(bp *Buffer, register int_t offset)
func (bp *Buffer) Ptr(offset int) int {
	// if (offset < 0)
	// 	return (bp->b_buf);
	// return (bp->b_buf+offset + (bp->b_buf + offset < bp->b_gap ? 0 : bp->b_egap-bp->b_gap));
	return 0
}

/* Given a inter into the buffer, convert it to a buffer offset */
//int_t pos(bp *Buffer, register char_t *cp)
func (bp *Buffer) Pos(cp int) int {
	// assert(bp->b_buf <= cp && cp <= bp->b_ebuf);
	// return (cp - bp->b_buf - (cp < bp->b_egap ? 0 : bp->b_egap - bp->b_gap));
	return 0
}

/* find the int for start of line ln */
func (bp *Buffer) LineToint(ln int) int {
	// int_t end_p = pos(curbp, curbp->b_ebuf);
	// int_t p, start;

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
	// int_t end_p = pos(curbp, curbp->b_ebuf);
	// int_t p;
	// int line;

	// *curline = -1;

	// for (p=0, line=0; p < end_p; p++) {
	// 	line += (*(ptr(curbp,p)) == '\n') ? 1 : 0;
	// 	*lastline = line;

	// 	if (*curline == -1 && p == curbp->b_int) {
	// 		*curline = (*(ptr(curbp,p)) == '\n') ? line : line + 1;
	// 	}
	// }

	// *lastline = *lastline + 1;

	// if (curbp->b_int == end_p)
	// 	*curline = *lastline;
}
