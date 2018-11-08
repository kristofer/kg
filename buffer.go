package kg

/*
* Buffer
 */

func buffer_init(bp *Buffer) {
	bp.b_mark = NOMARK
	bp.b_point = 0
	bp.b_cpoint = 0
	bp.b_page = 0
	bp.b_epage = 0
	bp.b_reframe = 0
	bp.b_size = 0
	bp.b_psize = 0
	bp.b_flags = 0
	bp.b_cnt = 0
	bp.b_buf = nil
	bp.b_ebuf = nil
	bp.b_gap = nil
	bp.b_egap = nil
	bp.b_next = nil
	bp.b_fname = ""
}

/* Find a buffer by filename or create if requested */
func find_buffer(fname string, cflag bool) *Buffer {
	var bp *Buffer
	var sb *Buffer

	bp = Bheadp
	for bp != nil {
		if string.Compare(fname, bp.b_fname) == 0 || string.Compare(fname, bp.b_bname) == 0 {
			return bp
		}
		bp = bp.b_next
	}

	if cflag != false {
		// if ((bp = (buffer_t *) malloc (sizeof (buffer_t))) == nil)
		// 	return (0);
		bp = make(Buffer())

		buffer_init(bp)
		//assert(bp != nil);

		/* find the place in the list to insert this buffer */
		if Bheadp == nil {
			Bheadp = bp
		} else if string.Compare(Bheadp.b_fname, fname) > 0 {
			/* insert at the begining */
			bp.b_next = Bheadp
			Bheadp = bp
		} else {
			for sb = Bheadp; sb.b_next != nil; sb = sb.b_next {
				if string.Compare(sb.b_next.b_fname, fname) > 0 {
					break
				}
			}
			/* and insert it */
			bp.b_next = sb.b_next
			sb.b_next = bp
		}
	}
	return bp
}

/* unlink from the list of buffers, free associated memory, assumes buffer has been saved if modified */
func delete_buffer(bp *Buffer) bool {
	var sb *Buffer

	/* we must have switched to a different buffer first */
	assert(bp != Curbp)

	/* if buffer is the head buffer */
	if bp == Bheadp {
		Bheadp = bp.b_next
	} else {
		/* find place where the bp buffer is next */
		for sb = Bheadp; sb.b_next != bp && sb.b_next != nil; sb = sb.b_next {
		}
		if sb.b_next == bp || sb.b_next == nil {
			sb.b_next = bp.b_next
		}
	}

	/* now we can delete */
	//free(bp.b_buf);
	bp.b_buf = nil
	//free(bp);
	bp = nil
	return true
}

func next_buffer() {
	// assert(Curbp != nil);
	// assert(Bheadp != nil);
	if Curbp != nil && Bheadp != nil {
		disassociate_b(Curwp)
		//Curbp = (Curbp.b_next != nil ? Curbp.b_next : Bheadp);
		if Curbp.b_next != nil {
			Curbp = Curbp.b_next

		} else {
			Curbp = Bheadp
		}
		associate_b2w(Curbp, Curwp)
	}
}

func get_buffer_name(bp *Buffer) string {
	if bp.b_fname != nil && bp.b_fname != "" {
		return bp.b_fname
	}
	return bp.b_bname
}

func count_buffers() int {
	var bp *Buffer
	i := 0

	for bp = Bheadp; bp != nil; bp = bp.b_next {
		i++
	}
	return i
}

func modified_buffers() bool {
	var bp *Buffer

	for bp = Bheadp; bp != nil; bp = bp.b_next {
		if bp.b_flags & B_MODIFIED {
			return true
		}
	}
	return false
}
