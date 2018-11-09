package kg

import (
	"fmt"
)

/* window.c, Atto Emacs, Hugh Barney, Public Domain, 2015 */

//#include "header.h"

//int win_cnt = 0;
var winCount = 0

type Window struct {
	Next     *Window /* Next window */
	Buffer   *Buffer /* Buffer displayed in window */
	CP       Point
	Mark     Point
	w_page   Point
	w_epage  Point
	TopPt    Point /* Origin 0 top row of window */
	w_rows   Point /* no. of rows of text in window */
	w_row    int   /* cursor row */
	w_col    int   /* cursor col */
	w_update bool
	w_name   string //[STRBUF_S];
} //window_t;

func NewWindow() *Window {
	wp := new(Window) //(window_t *)malloc(sizeof(window_t));

	//assert(wp != NULL); /* call fatal instead XXX */
	wp.Next = nil
	wp.Buffer = nil
	wp.CP = 0
	wp.Mark = NOMARK
	wp.w_top = 0
	wp.w_rows = 0
	wp.w_update = FALSE
	//sprintf(wp->w_name, "W%d", ++win_cnt);
	winCount++
	wp.w_name = fmt.Sprintf("W%d", winCount)
	return &wp
}

func one_window(wp *Window) {
	wp.w_top = 0
	wp.w_rows = LINES - 2
	wp.Next = nil
}

func split_window() {
	var wp *Window
	var wp2 *Window
	ntru, ntrl := 0, 0

	if Curwp.w_rows < 3 {
		msg("Cannot split a %d line window", Curwp.w_rows)
		return
	}

	wp = new_window()
	associate_b2w(curwp.Buffer, wp)
	b2w(wp) /* inherit buffer settings */

	ntru = (curwp.w_rows - 1) / 2    /* Upper size */
	ntrl = (curwp.w_rows - 1) - ntru /* Lower size */

	/* Old is upper window */
	Curwp.w_rows = ntru
	wp.w_top = Curwp.w_top + ntru + 1
	wp.w_rows = ntrl

	/* insert it in the list */
	wp2 = curwp.Next
	Curwp.Next = wp
	wp.Next = wp2
	redraw() /* mark the lot for update */
}

func next_window() {
	Curwp.w_update = true /* make sure modeline gets updated */
	//Curwp = (Curwp.Next == nil ? Wheadp : Curwp.Next)
	if Curwp.Next == nil {
		Curwp = Wheadp
	} else {
		Curwp = Curwp.Next
	}
	Curbp = Curwp.Buffer

	if Curbp.b_cnt > 1 {
		w2b(Curwp) /* push win vars to buffer */
	}
}

func delete_other_windows() {
	if Wheadp.Next == nil {
		msg("Only 1 window")
		return
	}
	free_other_windows(Curwp)
}

func free_other_windows(winp *Window) {
	var wp *Window
	var next *Window
	wp = Wheadp
	next = wp
	for next != nil {
		next = wp.Next /* get next before a call to free() makes wp undefined */
		if wp != winp {
			disassociate_b(wp) /* this window no longer references its buffer */
			//free(wp);
		}
		wp = next
	}

	Wheadp = winp
	Curwp = winp
	one_window(winp)
}

func associate_b2w(bp *Buffer, wp *Window) {
	//assert(bp != NULL);
	//assert(wp != NULL);
	if bp != nil && wp != nil {
		wp.Buffer = bp
		bp.b_cnt++
	}
}

func disassociate_b(wp *Window) {
	// assert(wp != NULL);
	// assert(wp->Buffer != NULL);
	if wp != nil && wp.w_buf != nil {
		wp.Buffer.b_cnt--
	}
}
