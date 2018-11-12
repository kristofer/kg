package kg

import "fmt"

/*
* Display
 */
/* display.c, Atto Emacs, Public Domain, Hugh Barney, 2016, Derived from: Anthony's Editor January 93 */

//#include "header.h"

/* Reverse scan for start of logical line containing offset */
func (bp *Buffer) LineStart(off int) int {
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
		off = +1
		return off
	}
	return 0
}

/* Forward scan for start of logical line segment (corresponds to screen line)  containing 'finish' */
func SegStart(bp *Buffer, start int, finish int) int {
	var p *rune
	c := 0
	scan := start

	for scan < finish {
		//p = ptr(bp, scan);
		p = bp.GetCurrentRune()
		if *p == '\n' {
			c = 0
			start = scan + 1
		} else if COLS <= c {
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
func SegNext(bp *Buffer, start, finish int) int {
	// char_t *p;
	// int c = 0;
	var p *rune
	var pptr int
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
func UpUp(bp *Buffer, off int) int {
	curr := bp.LineStart(off)
	seg := SegStart(bp, curr, off)
	if curr < seg {
		off = SegStart(bp, curr, seg-1)
	} else {
		off = SegStart(bp, bp.LineStart(curr-1), curr-1)
	}
	return off
}

/* Move down one screen line */
func (bp *Buffer) DownDown(off int) int {
	return (SegNext(bp, bp.LineStart(off), off))
}

/* Return the offset of a column on the specified line */
func (bp *Buffer) OffsetForColumn(offset int, column int) int {
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

//func (wp *Window) Display(flag byte) {
//char_t *p;
// i, j, k, nch := 0;

// bp := wp.Buffer

// token_type := ID_DEFAULT

// /* find start of screen, handle scroll up off page or top of file  */
// /* int is always within b_page and b_epage */
// if (bp.int < bp.b_page)
// 	bp.b_page = SegStart(bp, bp.LineStart(bp.int), bp.int)

// /* reframe when scrolled off bottom */
// if (bp.b_reframe == 1 || (bp.b_epage <= bp.int && curbp.int != pos(curbp, curbp.b_ebuf))) {
// 	bp.b_reframe = 0;
// 	/* Find end of screen plus one. */
// 	bp.b_page = dndn(bp, bp.int);
// 	/* if we scoll to EOF we show 1 blank line at bottom of screen */
// 	if (pos(bp, bp.b_ebuf) <= bp.b_page) {
// 		bp.b_page = pos(bp, bp.b_ebuf);
// 		i = wp.w_rows - 1;
// 	} else {
// 		i = wp.w_rows - 0;
// 	}
// 	/* Scan backwards the required number of lines. */
// 	while (0 < i--)
// 		bp.b_page = upup(bp, bp.b_page);
// }

// move(wp.TopPt, 0); /* start from top of window */
// i = wp.TopPt;
// j = 0;
// bp.b_epage = bp.b_page;
// set_parse_state(bp, bp.b_epage); /* are we in a multline comment ? */

// /* paint screen from top of page until we hit maxline */
// for {
// 	/* reached int - store the cursor position */
// 	if (bp.int == bp.b_epage) {
// 		bp.b_row = i;
// 		bp.b_col = j;
// 	}
// 	p = ptr(bp, bp.b_epage);
// 	nch = 1;
// 	if (wp.w_top + wp.w_rows <= i || bp.b_ebuf <= p) /* maxline */
// 		break;
// 	if (*p != '\r') {
// 		nch = utf8_size(*p);
// 		if ( nch > 1) {
// 			wchar_t c;
// 			/* reset if invalid multi-byte character */
// 			if (mbtowc(&c, (char*)p, 6) < 0) mbtowc(nil, nil, 0);
// 			j += wcwidth(c) < 0 ? 1 : wcwidth(c);
// 			display_utf8(bp, *p, nch);
// 		} else if (isprint(*p) || *p == '\t' || *p == '\n') {
// 			j += *p == '\t' ? 8-(j&7) : 1;
// 			token_type = parse_text(bp, bp.b_epage);
// 			attron(COLOR_PAIR(token_type));
// 			addch(*p);
// 		} else {
// 			const char *ctrl = unctrl(*p);
// 			j += (int) strlen(ctrl);
// 			addstr(ctrl);
// 		}
// 	}
// 	if (*p == '\n' || COLS <= j) {
// 		j -= COLS;
// 		if (j < 0)
// 			j = 0;
// 		++i;
// 	}
// 	bp.b_epage = bp.b_epage + nch;
// }

// /* replacement for clrtobot() to bottom of window */
// for (k=i; k < wp.w_top + wp.w_rows; k++) {
// 	move(k, j) /* clear from very last char not start of line */
// 	clrtoeol()
// 	j = 0 /* thereafter start of line */
// }

// //b2w(wp); /* save buffer stuff on window */
// PushBuffer2Window(wp)
// modeline(wp)
// if (wp == CurrentWin && flag) {
// 	DisplayMsg()
// 	move(bp.CursorRow, bp.CursorCol) /* set cursor */
// 	refresh()
// }
// wp.Updated = false
//}

func DisplayUTF8(bp *Buffer, sbuf string) {
	// char sbuf[6];
	// int i = 0;

	// for (i=0; i<n; i++)
	// 	sbuf[i] = *ptr(bp, bp.b_epage + i);
	// sbuf[n] = '\0';
	addstr(sbuf)
}

func ModeLine(wp *Window) {
	i := 0
	var lch, mch, och rune

	//standout();
	move(wp.TopPt+wp.Rows, 0)
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
	if wp.Buffer.Flags&B_OVERWRITE != 0 {
		och = 'O'
	}

	fmt.Sprintf(temp, "%c%c%c Atto: %c%c %s", lch, och, mch, lch, lch, GetBufferName(wp.Buffer))
	addstr(temp) // term

	for i = len(temp) + 1; i <= COLS; i++ {
		addch(lch) // term
	}
	//standend();
}

func DisplayMsg() {
	move(MSGLINE, 0) // (Lines-1)
	if msgflag {
		addstr(msgline)
		msgflag = false
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
