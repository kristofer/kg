package tkg

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
	//termbox "github.com/gdamore/tcell/termbox"
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

// OnKey handles the insertion of non-control/editor keys
func (wp *Window) OnKey(ev *termbox.Event) {
	switch ev.Key {
	case termbox.KeySpace:
		wp.Buffer.AddRune(' ')
	case termbox.KeyEnter, termbox.KeyCtrlJ:
		wp.Buffer.AddRune('\n')
	case termbox.KeyTab:
		wp.Buffer.AddRune('\t')
		//g.set_overlay_mode(init_extended_mode(g))
	// case termbox.KeyCtrlS:
	// 	g.set_overlay_mode(init_isearch_mode(g, false))
	// case termbox.KeyCtrlR:
	// 	g.set_overlay_mode(init_isearch_mode(g, true))
	default:
		// if ev.Mod&termbox.ModAlt != 0 && e.OnAltKey(ev) {
		// 	break
		// }
		//ch := ev.Ch
		//log.Printf("Win OnKey %#U Point is %d\n", ch, wp.Buffer.Point())
		wp.Buffer.AddRune(ev.Ch)
	}

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

// SyncBuffer
func SyncBuffer(w *Window) { //sync w2b win to buff
	b := w.Buffer
	b.SetPoint(w.Point)
	b.PageStart = w.WinStart
	b.PageEnd = w.WinEnd
	b.PointRow = w.CurRow
	b.PointCol = w.CurCol

	/* fixup Pointers in other windows of the same buffer, if size of edit text changed */
	// if b.Point() > b.OrigPoint {
	// 	sizeDelta := b.TextSize - b.PrevSize
	// 	b.MoveGap(sizeDelta)
	// 	b.PageStart += sizeDelta
	// 	b.PageEnd += sizeDelta
	// }
}

// PushBuffer2Window
func PushBuffer2Window(w *Window) { // b2w
	b := w.Buffer
	w.Point = b.Point()
	w.WinStart = b.PageStart
	w.WinEnd = b.PageEnd
	w.CurRow = b.PointRow
	w.CurCol = b.PointRow
	b.TextSize = b.BufferLen()
}
