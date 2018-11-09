package kg

/*
* Buffer
 */

type Buffer struct {
	Next      *Buffer /* b_next Link to next buffer_t */
	Mark      Point   /* b_mark the mark */
	Point     Point   /* b_point the point */
	OrigPoint Point   /* b_cpoint the original current point, used for mutliple window displaying */
	PageStart Point   /* b_page start of page */
	PageEnd   Point   /* b_epage end of page */
	Reframe   Point   /* b_reframe force a reframe of the display */
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

/* Find a buffer by filename or create if requested */
func find_buffer(fname string, cflag bool) *Buffer {
	var bp *Buffer
	var sb *Buffer

	bp = Bheadp
	for bp != nil {
		if string.Compare(fname, bp.filename) == 0 || string.Compare(fname, bp.Buffername) == 0 {
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
