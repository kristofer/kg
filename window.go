package kg

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
	//termbox "github.com/gdamore/tcell/termbox"
)

var winCount = 0

// Window main type
type Window struct {
	Editor   *Editor
	Next     *Window /* w_next Next window */
	Buffer   *Buffer /* w_bufp Buffer displayed in window */
	Point    int     // w_point
	Mark     int     // w_mark
	WinStart int     // w_page
	WinEnd   int     // w_epage
	TopPt    int     /* w_top Origin 0 top row of window  on screen */
	Rows     int     /* w_rows no. of rows of text in window */
	Row      int     /* w_row cursor row */
	Col      int     /* w_col cursor col */
	Updated  bool    // int w_update
	Name     string  // w_name[STRBUF_S];
}

// NewWindow xxx
func NewWindow(e *Editor) *Window {
	wp := &Window{}
	wp.Editor = e
	wp.Next = nil
	wp.Buffer = nil
	wp.Point = 0
	wp.Mark = nomark
	wp.TopPt = 0
	wp.Rows = 0
	wp.Updated = false
	winCount++
	wp.Name = fmt.Sprintf("W%d", winCount)
	return wp
}

// OneWindow xxx
func (wp *Window) OneWindow() {
	wp.TopPt = 0
	wp.Rows = wp.Editor.Lines - 3
	// log.Printf("OneWindow rows %d line %d\n", wp.Rows, wp.Editor.Lines)
	wp.Next = nil
}

// WindowResize xxx
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
	default:
		if ev.Mod&termbox.ModAlt != 0 && wp.Editor.OnAltKey(ev) {
			// log.Println("Alt!", ev.Key, ev.Ch)
			break
		}
		wp.Buffer.AddRune(ev.Ch)
	}
}

// AssociateBuffer xxx
func (wp *Window) AssociateBuffer(bp *Buffer) {
	if bp != nil && wp != nil {
		wp.Buffer = bp
		bp.WinCount++
	}
}

// DisassociateBuffer xxx
func (wp *Window) DisassociateBuffer() {
	if wp != nil && wp.Buffer != nil {
		wp.Buffer.WinCount--
		wp.Buffer = nil
	}
}

// SyncBuffer xxx
func window2Buffer(w *Window) {
	b := w.Buffer
	b.SetPoint(w.Point)
	b.PageStart = w.WinStart
	b.PageEnd = w.WinEnd
	b.PointRow = w.Row
	b.PointCol = w.Col
	// this should be figured out.
	/* fixup Pointers in other windows of the same buffer, if size of edit text changed */
	if b.Point > b.OrigPoint {
		sizeDelta := b.TextSize - b.PrevSize
		b.MoveGap(sizeDelta)
		b.SetPoint(b.Point + sizeDelta)
		b.PageStart += sizeDelta
		b.PageEnd += sizeDelta
	}
}

// PushBuffer2Window xxx
func buffer2Window(w *Window) {
	b := w.Buffer
	w.Point = b.Point
	w.WinStart = b.PageStart
	w.WinEnd = b.PageEnd
	w.Row = b.PointRow
	w.Col = b.PointCol
	//b.TextSize = b.TextSize
}
