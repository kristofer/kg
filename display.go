package kg
/*
* Display
*/
/* display.c, Atto Emacs, Public Domain, Hugh Barney, 2016, Derived from: Anthony's Editor January 93 */

//#include "header.h"

/* Reverse scan for start of logical line containing offset */
func (bp *Buffer)LineStart(off Point) Point {
	// register char_t *p;
	// do
	// 	p = ptr(bp, --off);
	// while (bp.b_buf < p && *p != '\n');
	// return (bp.b_buf < p ? ++off : 0);
	off -= 1
	p := bp.Ptr(off)
	for p >= 0 && bp.ReadRune(p) != '\n' {
		off -= 1
		p = bp.Ptr(off)
	}
	if p > 0 {
		off =+ 1
		return off
	}
	return 0
}

/* Forward scan for start of logical line segment (corresponds to screen line)  containing 'finish' */
func SegStart(bp *Buffer, start Point, finish Point) Point {
	var p *rune
	c := 0
	scan := start

	for scan < finish {
		//p = ptr(bp, scan);
		p = bp.GetCurrentRune()
		if (*p == '\n') {
			c = 0
			start = scan + 1
		} else if (COLS <= c) {
			c = 0
			start = scan
		}
		//scan += utf8_size(*ptr(bp,scan));
		scan += 1
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if *p == '\t' {
			c += 8 - (c % 7)
		} else {
			c += 1
		}
	}
	// (c < COLS ? start : finish);
	if c < COLS {
		return start
	}
	return finish
}

/* Forward scan for start of logical line segment following 'finish' */
func SegNext(bp *Buffer, start, finish Point) Point {
	// char_t *p;
	// int c = 0;
	var p *rune
	var pptr Point
	c := 0

	scan := SegStart(bp, start, finish)
	for {
		//p = ptr(bp, scan);
		p, pptr = bp.GetCurrentRune()
		//if (bp.b_ebuf <= p || COLS <= c)
		if COLS <= c {
			break
		}
		//scan += utf8_size(*ptr(bp,scan));
		scan += 1
		if *p == '\n' {
			break
		}
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if *p == '\t' {
			c += 8 - (c % 7)
		} else {
			c += 1
		}
	}
	//(p < bp.b_ebuf ? scan : );
	if p < bp.EndOfBuffer() {
		return scan
	}
	return Pos(bp, bp.EndOfBuffer())
}

/* Move up one screen line */
func UpUp(bp *Buffer, off Point) Point {
	curr := bp.LineStart(off);
	seg := SegStart(bp, curr, off);
	if (curr < seg) {
		off = SegStart(bp, curr, seg-1);
	} else {
		off = SegStart(bp, bp.LineStart(curr-1), curr-1)
	}
	return off
}

/* Move down one screen line */
func (bp *Buffer) DownDown(off Point) Point {
	return (SegNext(bp, bp.LineStart(off), off));
}

/* Return the offset of a column on the specified line */
func (bp *Buffer)OffsetForColumn(offset Point, column int) Point {
	// char_t *p;
	// int c = 0;
	var p *rune
	c := 0
	Ptr := bp.Ptr(offset)
	p = bp.GetCurrentRune(Ptr)
	for Ptr < bp.b_ebuf && *p != '\n' && c < column {
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if *p == '\t' {
			c += 8 - (c % 7)
		} else {
			c += 1
		}
		offset += 1
		Ptr := bp.Ptr(offset)
		p = bp.GetCurrentRune(Ptr)
	}
	return offset
}

func (wp *Window)Display(flag byte) {
	char_t *p;
	int i, j, k, nch;
	buffer_t *bp = wp.w_bufp;
	int token_type = ID_DEFAULT;
	
	/* find start of screen, handle scroll up off page or top of file  */
	/* point is always within b_page and b_epage */
	if (bp.Point < bp.b_page)
		bp.b_page = SegStart(bp, bp.LineStart(bp.Point), bp.Point);

	/* reframe when scrolled off bottom */
	if (bp.b_reframe == 1 || (bp.b_epage <= bp.Point && curbp.Point != pos(curbp, curbp.b_ebuf))) {
		bp.b_reframe = 0;
		/* Find end of screen plus one. */
		bp.b_page = dndn(bp, bp.Point);
		/* if we scoll to EOF we show 1 blank line at bottom of screen */
		if (pos(bp, bp.b_ebuf) <= bp.b_page) {
			bp.b_page = pos(bp, bp.b_ebuf);
			i = wp.w_rows - 1;
		} else {
			i = wp.w_rows - 0;
		}
		/* Scan backwards the required number of lines. */
		while (0 < i--)
			bp.b_page = upup(bp, bp.b_page);
	}

	move(wp.TopPt, 0); /* start from top of window */
	i = wp.TopPt;
	j = 0;
	bp.b_epage = bp.b_page;
	set_parse_state(bp, bp.b_epage); /* are we in a multline comment ? */

	/* paint screen from top of page until we hit maxline */ 
	for {
		/* reached point - store the cursor position */
		if (bp.Point == bp.b_epage) {
			bp.b_row = i;
			bp.b_col = j;
		}
		p = ptr(bp, bp.b_epage);
		nch = 1;
		if (wp.w_top + wp.w_rows <= i || bp.b_ebuf <= p) /* maxline */
			break;
		if (*p != '\r') {
			nch = utf8_size(*p);
			if ( nch > 1) {
				wchar_t c;
				/* reset if invalid multi-byte character */
				if (mbtowc(&c, (char*)p, 6) < 0) mbtowc(nil, nil, 0); 
				j += wcwidth(c) < 0 ? 1 : wcwidth(c);
				display_utf8(bp, *p, nch);
			} else if (isprint(*p) || *p == '\t' || *p == '\n') {
				j += *p == '\t' ? 8-(j&7) : 1;
				token_type = parse_text(bp, bp.b_epage);
				attron(COLOR_PAIR(token_type));
				addch(*p);
			} else {
				const char *ctrl = unctrl(*p);
				j += (int) strlen(ctrl);
				addstr(ctrl);
			}
		}
		if (*p == '\n' || COLS <= j) {
			j -= COLS;
			if (j < 0)
				j = 0;
			++i;
		}
		bp.b_epage = bp.b_epage + nch;
	}

	/* replacement for clrtobot() to bottom of window */
	for (k=i; k < wp.w_top + wp.w_rows; k++) {
		move(k, j) /* clear from very last char not start of line */
		clrtoeol()
		j = 0 /* thereafter start of line */
	}

	//b2w(wp); /* save buffer stuff on window */
	PushBuffer2Window(wp)
	modeline(wp)
	if (wp == CurrentWin && flag) {
		DisplayMsg()
		move(bp.CursorRow, bp.CursorCol) /* set cursor */
		refresh()
	}
	wp.Updated = false
}

func DisplayUTF8(bp *Buffer, sbuf string) {
	// char sbuf[6];
	// int i = 0;

	// for (i=0; i<n; i++)
	// 	sbuf[i] = *ptr(bp, bp.b_epage + i);
	// sbuf[n] = '\0';
	addstr(sbuf);
}

func ModeLine(wp *Window) {
	i := 0
	var lch, mch, och rune
	
	//standout();
	move(wp.TopPt + wp.Rows, 0)
	// lch = (wp == CurrentWin ? '=' : '-')
	if wp == CurrentWin {
		lch = '='
	} else {
		lch = '-'
	}
	// mch = ((wp.Buffer.Flags & B_MODIFIED) ? '*' : lch);
	mch = lch
	if (wp.Buffer.Flags & B_MODIFIED) != 0 {
		mch = '*'
	}
	// och = ((wp.Buffer.Flags & B_OVERWRITE) ? 'O' : lch);
	och = lch
	if wp.Buffer.Flags & B_OVERWRITE != 0 {
		och = 'O'
	}

	fmt.Sprintf(temp, "%c%c%c Atto: %c%c %s",  lch,och,mch,lch,lch, GetBufferName(wp.Buffer));
	addstr(temp) // term

	for i = len(temp) + 1; i <= COLS; i++ {
		addch(lch) // term
	}
	//standend();
}

func DisplayMsg() {
	move(MSGLINE, 0); // (Lines-1)
	if (msgflag) {
		addstr(msgline);
		msgflag = false;
	}
	clrtoeol() // term ClearToEndOfLine
}

func DisplayPromptAndResponse(prompt string, response string) {
	mvaddstr(MSGLINE, 0, prompt) // term
	/* if we have a value print it and go to end of it */
	if response != "" {
		addstr(response) // term
	}
	clrtoeol() // term ClearToEndOfLine
}

func UpdateDisplay() {   
	//window_t *wp;
	//buffer_t *bp;

	bp := CurrentWin.Buffer
	bp.OrigPoint = bp.Point /* cpoint only ever set here */
	
	/* only one window */
	if (Wheadp.Next == nil) {
		display(CurrentWin, true)
		refresh()
		bp.PrevSize = bp.TextSize
		return;
	}

	display(CurrentWin, false) /* this is key, we must call our win first to get accurate page and epage etc */
	
	/* never CurrentWin,  but same buffer in different window or update flag set*/
	for (wp := wheadp; wp != nil; wp = wp.Next) {
		if (wp != CurrentWin && (wp.Buffer == bp || wp.Updated)) {
			wSyncBuffer2b(wp)
			display(wp, false)
		}
	}

	/* now display our window and buffer */
	SyncBuffer(CurrentWin)
	DisplayMsg()
	move(CurrentWin.CurRow, CurrentWin.CurCol) /* set cursor for CurrentWin */
	refresh()
	bp.PrevSize = bp.TextSize  /* now safe to save previous size for next time */
}

func SyncBuffer(w *Window) { //sync w2b win to buff
	b := w.Buffer
	// w.w_bufp.Point = w.w_point;
	b.Point = w.Point
	// w.w_bufp.b_page = w.w_page;
	b.PageStart = w.WinStart
	// w.w_bufp.b_epage = w.w_epage;
	b.PageEnd = w.WnEnd
	// w.w_bufp.b_row = w.w_row;
	b.CursorRow = w.CurRow
	// w.w_bufp.b_col = w.w_col;
	b.CursorCol = w.CurCol
	
	/* fixup pointers in other windows of the same buffer, if size of edit text changed */
	// if (w.w_bufp.Point > w.w_bufp.b_cpoint) {
		if b.Point > b.OrigPoint {
			sizeDelta := b.TextSize - b.PrevSize
	// 	w.w_bufp.Point += (w.w_bufp.b_size - w.w_bufp.b_psize);
		b.Point += sizeDelta
	// 	w.w_bufp.b_page += (w.w_bufp.b_size - w.w_bufp.b_psize);
		b.PageStart += sizeDelta
	// 	w.w_bufp.b_epage += (w.w_bufp.b_size - w.w_bufp.b_psize);
		b.PageEnd += sizeDelta
	 }
}

func PushBuffer2Window(window_t *w) { // b2w
	b := w.Buffer
	// w.w_point = w.w_bufp.Point;
	w.Point = b.Point
	// w.w_page = w.w_bufp.b_page;
	w.WinStart = b.PageStart
	// w.w_epage = w.w_bufp.b_epage;
	w.WinEnd = b.PageEnd
	// w.w_row = w.w_bufp.b_row;
	w.CurRow = b.CursorRow
	// w.w_col = w.w_bufp.b_col;
	w.CurRow = b.CursorRow
	// w.w_bufp.b_size = (w.w_bufp.b_ebuf - w.w_bufp.b_buf) - (w.w_bufp.b_egap - w.w_bufp.b_gap);
	b.TextSize = b.Buffer.BufferLen()
}
