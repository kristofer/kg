package kg

import (
	"fmt"
)

/* window.c, Atto Emacs, Hugh Barney, Public Domain, 2015 */

//#include "header.h"

//int win_cnt = 0; Global
var winCount = 0

type Window struct {
	Next   *Window /* w_next Next window */
	Buffer *Buffer /* w_bufp Buffer displayed in window */
	Point  int     // w_point
	//
	Mark     int    // w_mark
	WinStart int    // w_page
	WinEnd   int    // w_epage
	TopPt    int    /* w_top Origin 0 top row of window  on screen */
	Rows     int    /* w_rows no. of rows of text in window */
	CurRow   int    /* w_row cursor row */
	CurCol   int    /* w_col cursor col */
	Updated  bool   // int w_update
	Name     string // w_name[STRBUF_S];
} //window_t;

func NewWindow() *Window {
	wp := new(Window) //(window_t *)malloc(sizeof(window_t));

	//assert(wp != NULL); /* call fatal instead XXX */
	wp.Next = nil
	wp.Buffer = nil
	wp.Point = 0
	wp.Mark = NOMARK
	wp.TopPt = 0
	wp.Rows = 0
	wp.Updated = FALSE
	//sprintf(wp->Name, "W%d", ++win_cnt);
	winCount++
	wp.Name = fmt.Sprintf("W%d", winCount)
	return &wp
}

func (wp *Window) OneWindow() {
	wp.TopPt = 0
	wp.Rows = LINES - 2
	wp.Next = nil
}

func (wp *Window) SplitWindow() {
	//var wp *Window
	var wp2 *Window
	ntru, ntrl := 0, 0

	if Curwp.Rows < 3 {
		msg("Cannot split a %d line window", Curwp.Rows)
		return
	}

	wp = NewWindow()
	bp.AssociateBuffer(curwp.Buffer)
	b2w(wp) /* inherit buffer settings */

	ntru = (curwp.Rows - 1) / 2    /* Upper size */
	ntrl = (curwp.Rows - 1) - ntru /* Lower size */

	/* Old is upper window */
	Curwp.Rows = ntru
	wp.TopPt = Curwp.TopPt + ntru + 1
	wp.Rows = ntrl

	/* insert it in the list */
	wp2 = Curwp.Next
	Curwp.Next = wp
	wp.Next = wp2
	redraw() /* mark the lot for update */
}

func (wp *Window) NextWindow() {
	Curwp.Updated = true /* make sure modeline gets updated */
	//Curwp = (Curwp.Next == nil ? Wheadp : Curwp.Next)
	if Curwp.Next == nil {
		Curwp = Wheadp
	} else {
		Curwp = Curwp.Next
	}
	Curbp = Curwp.Buffer

	if Curbp.WinCount > 1 {
		w2b(Curwp) /* push win vars to buffer */
	}
}

func (wp *Window) DeleteOtherWindows() {
	if Wheadp.Next == nil {
		msg("Only 1 window")
		return
	}
	FreeOtherWindows(wp)
}

func (winp *Window) FreeOtherWindows() {
	var wp *Window
	var next *Window
	wp = Wheadp
	next = wp
	for next != nil {
		next = wp.Next /* get next before a call to free() makes wp undefined */
		if wp != winp {
			DisassociateBuffer(wp) /* this window no longer references its buffer */
			//free(wp);
		}
		wp = next
	}

	Wheadp = winp
	Curwp = winp
	OneWindow(winp)
}

func (bp *Buffer) AssociateBuffer(wp *Window) {
	//assert(bp != NULL);
	//assert(wp != NULL);
	if bp != nil && wp != nil {
		wp.Buffer = bp
		bp.b_cnt++
	}
}

func (wp *Window) DisassociateBuffer() {
	// assert(wp != NULL);
	// assert(wp->Buffer != NULL);
	if wp != nil && wp.Buffer != nil {
		wp.Buffer.WinCount--
		wp.Buffer = nil
	}
}
