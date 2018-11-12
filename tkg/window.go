package tkg

import (
	"fmt"
)

var winCount = 0

// Window main type
type Window struct {
	Editor *Editor
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

func NewWindow(e *Editor) *Window {
	wp := &Window{} // new(Window) //(window_t *)malloc(sizeof(window_t));

	//assert(wp != NULL); /* call fatal instead XXX */
	wp.Editor = e
	wp.Next = nil
	wp.Buffer = nil
	wp.Point = 0
	wp.Mark = NOMARK
	wp.TopPt = 0
	wp.Rows = 0
	wp.Updated = false
	//sprintf(wp->Name, "W%d", ++win_cnt);
	winCount++
	wp.Name = fmt.Sprintf("W%d", winCount)
	return wp
}

func (wp *Window) OneWindow() {
	wp.TopPt = 0
	wp.Rows = wp.Editor.Lines - 2
	wp.Next = nil
}

func (wp *Window) WindowResize() {
	wp.Editor.CurrentWindow.OneWindow()
}

func (wp *Window) SplitWindow() {
	var editor = wp.Editor
	//var wp *Window
	var wp2 *Window
	ntru, ntrl := 0, 0

	if editor.CurrentWindow.Rows < 3 {
		editor.msg("Cannot split a %d line window", editor.CurrentWindow.Rows)
		return
	}

	wp = NewWindow(editor)
	wp.AssociateBuffer(editor.CurrentWindow.Buffer)
	//b2w(wp) /* inherit buffer settings */

	ntru = (editor.CurrentWindow.Rows - 1) / 2    /* Upper size */
	ntrl = (editor.CurrentWindow.Rows - 1) - ntru /* Lower size */

	/* Old is upper window */
	editor.CurrentWindow.Rows = ntru
	wp.TopPt = editor.CurrentWindow.TopPt + ntru + 1
	wp.Rows = ntrl

	/* insert it in the list */
	wp2 = editor.CurrentWindow.Next
	editor.CurrentWindow.Next = wp
	wp.Next = wp2
	//redraw() /* mark the lot for update */
}

// NextWindow
func (wp *Window) NextWindow() {
	var editor = wp.Editor
	editor.CurrentWindow.Updated = true /* make sure modeline gets updated */
	//Curwp = (Curwp.Next == nil ? Wheadp : Curwp.Next)
	if editor.CurrentWindow.Next == nil {
		editor.CurrentWindow = editor.RootWindow
	} else {
		editor.CurrentWindow = editor.CurrentWindow.Next
	}
	editor.CurrentBuffer = editor.CurrentWindow.Buffer

	if editor.CurrentBuffer.WinCount > 1 {
		//w2b(Curwp) /* push win vars to buffer */
	}
}

// DeleteOtherWindows
func (wp *Window) DeleteOtherWindows() {
	if wp.Next == nil {
		wp.Editor.msg("Only 1 window")
		return
	}
	wp.FreeOtherWindows()
}

// FreeOtherWindows
func (winp *Window) FreeOtherWindows() {
	var editor = winp.Editor
	var wp *Window
	var next *Window
	wp = editor.RootWindow
	next = wp
	for next != nil {
		next = wp.Next /* get next before a call to free() makes wp undefined */
		if wp != winp {
			wp.DisassociateBuffer() /* this window no longer references its buffer */
		}
		wp = next
	}

	editor.RootWindow = winp
	editor.CurrentWindow = winp
	winp.OneWindow()
}

// AssociateBuffer
func (wp *Window) AssociateBuffer(bp *Buffer) {
	if bp != nil && wp != nil {
		wp.Buffer = bp
		bp.WinCount++
	}
}

// DisassociateBuffer
func (wp *Window) DisassociateBuffer() {
	// assert(wp != NULL);
	// assert(wp->Buffer != NULL);
	if wp != nil && wp.Buffer != nil {
		wp.Buffer.WinCount--
		wp.Buffer = nil
	}
}

func SyncBuffer(w *Window) { //sync w2b win to buff
	b := w.Buffer
	// w.w_bufp.int = w.w_int;
	b.SetPoint(w.Point)
	// w.w_bufp.b_page = w.w_page;
	b.PageStart = w.WinStart
	// w.w_bufp.b_epage = w.w_epage;
	b.PageEnd = w.WinEnd
	// w.w_bufp.b_row = w.w_row;
	b.PointRow = w.CurRow
	// w.w_bufp.b_col = w.w_col;
	b.PointCol = w.CurCol

	/* fixup Pointers in other windows of the same buffer, if size of edit text changed */
	// if (w.w_bufp.int > w.w_bufp.b_cint) {
	if b.Point() > b.OrigPoint {
		sizeDelta := b.TextSize - b.PrevSize
		// 	w.w_bufp.int += (w.w_bufp.b_size - w.w_bufp.b_psize);
		b.MoveGap(sizeDelta)
		// 	w.w_bufp.b_page += (w.w_bufp.b_size - w.w_bufp.b_psize);
		b.PageStart += sizeDelta
		// 	w.w_bufp.b_epage += (w.w_bufp.b_size - w.w_bufp.b_psize);
		b.PageEnd += sizeDelta
	}
}

func PushBuffer2Window(w *Window) { // b2w
	b := w.Buffer
	w.Point = b.Point()
	w.WinStart = b.PageStart
	w.WinEnd = b.PageEnd
	w.CurRow = b.PointRow
	w.CurCol = b.PointRow
	b.TextSize = b.BufferLen()
}
