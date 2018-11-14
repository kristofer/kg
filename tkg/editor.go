package tkg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	termbox "github.com/gdamore/tcell/termbox"
)

// var LINES = 1
// var COLS = 10
// var MSGLINE = (LINES - 1)

const (
	VERSION     = "kg 1.0, Public Domain, November 2018, Kristofer Younger,  No warranty."
	PROG_NAME   = "kg"
	B_MODIFIED  = 0x01 /* modified buffer */
	B_OVERWRITE = 0x02 /* overwite mode */
	//LINES            = 24
	NOMARK           = -1
	CHUNK            = 8096
	K_BUFFER_LENGTH  = 256
	TEMPBUF          = 512
	STRBUF_L         = 256
	STRBUF_M         = 64
	STRBUF_S         = 16
	MIN_GAP_EXPAND   = 512
	TEMPFILE         = "/tmp/kgXXXXXX"
	F_NONE           = 0
	F_CLEAR          = 1
	ID_DEFAULT       = 1
	ID_SYMBOL        = 2
	ID_MODELINE      = 3
	ID_DIGITS        = 4
	ID_LINE_COMMENT  = 5
	ID_BLOCK_COMMENT = 6
	ID_DOUBLE_STRING = 7
	ID_SINGLE_STRING = 8
)

type Keymapt struct {
	KeyDesc  string
	KeyBytes string
	Do       func(*Editor) // function to call for Keymap-ping
}

type Editor struct {
	EventChan     chan termbox.Event
	CurrentBuffer *Buffer /* current buffer */
	RootBuffer    *Buffer /* head of list of buffers */
	CurrentWindow *Window
	RootWindow    *Window
	// status vars
	// done int                /* Quit flag. */
	Done       bool   /* Quit flag. */
	Msgflag    bool   /* True if msgline should be displayed. */
	Nscrap     int    /* Length of scrap buffer. */
	Scrap      string /* Allocated scrap buffer. */
	Input      rune   // RUNE?????
	Msgline    string /* Message line input/output buffer. */
	Temp       string /* Temporary buffer. */
	Searchtext string
	Replace    string
	Keymap     []Keymapt
	KeysMapMap map[string]*Keymapt
	//
	Lines   int
	Cols    int
	FGColor termbox.Attribute
	BGColor termbox.Attribute
}

// StartEditor is the old C main function
func (e *Editor) StartEditor(argv []string, argc int) {
	// log setup....
	f, err := os.OpenFile("logfile", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	f.Truncate(0)
	log.Println("Start of Log...")
	//
	e.FGColor = termbox.ColorDefault
	e.BGColor = termbox.ColorWhite
	err = termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	e.Cols, e.Lines = termbox.Size()

	if argc > 1 {
		e.CurrentBuffer = e.FindBuffer(argv[1], true)
		e.InsertFile(argv[1], false)
		/* Save filename irregardless of load() success. */
		e.CurrentBuffer.Filename = argv[1]
	} else {
		e.msg("NO file to open, creating scratch buffer")
		e.CurrentBuffer = e.FindBuffer("*scratch*", true)
		e.CurrentBuffer.Buffername = "*scratch*"
		s := "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut\nlabore et dolore magna aliqua. "

		e.CurrentBuffer.SetText(s)
	}
	e.CurrentWindow = NewWindow(e)
	e.RootWindow = e.CurrentWindow
	e.CurrentWindow.OneWindow()
	e.CurrentWindow.AssociateBuffer(e.CurrentBuffer)

	if !(e.CurrentBuffer.GrowGap(CHUNK)) {
		panic("%s: Failed to allocate required memory.\n")
	}
	//e.Key_map = e.Keymap
	e.Keymap = keymap
	e.KeysMapMap = make(map[string]*Keymapt)
	for _, k := range e.Keymap {
		e.KeysMapMap[k.KeyBytes] = &k
	}
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	//termbox.SetInputMode(termbox.InputAlt)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	e.UpdateDisplay()
	termbox.Flush()

	e.EventChan = make(chan termbox.Event, 20)
	go func() {
		for {
			e.EventChan <- termbox.PollEvent()
		}
	}()
	for {
		select {
		case ev := <-e.EventChan:
			ok := e.HandleEvent(&ev)
			if !ok {
				return
			}
			//e.consume_more_events()
			e.ConsumeMoreEvents()
			e.UpdateDisplay()
			termbox.Flush()
		}
	}
	return

}

// HandleEvent
func (e *Editor) HandleEvent(ev *termbox.Event) bool {
	switch ev.Type {
	case termbox.EventKey:
		e.OnSysKey(ev)
		if e.Done {
			return false
		}
		e.OnKey(ev)
	case termbox.EventResize:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		// g.resize()
		// if g.overlay != nil {
		// 	g.overlay.on_resize(ev)
		// }
		e.Cols, e.Lines = termbox.Size()
		e.msg("Resize: h %d,w %d", e.Lines, e.Cols)
		e.CurrentWindow.WindowResize()
		e.UpdateDisplay()
	case termbox.EventMouse:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		//log.Println("[", e.Lines, e.Cols, ev, "event", ctrlxpressed)
		e.msg("Mouse: %d,%d [%d %d]", ev.MouseX, ev.MouseY, e.Cols, e.Lines)
		// TODO: need to set the Point to mouse click location.
		e.UpdateDisplay()
		termbox.Flush()
	case termbox.EventError:
		panic(ev.Err)
	}

	// just dump the current view location from the view to the buffer
	// after each event, it's cheap and does what it needs to be done
	//v := g.active.leaf
	//v.buf.loc = v.view_location
	return true
}

// ConsumeMoreEvents handles
func (e *Editor) ConsumeMoreEvents() bool {
	for {
		select {
		case ev := <-e.EventChan:
			ok := e.HandleEvent(&ev)
			if !ok {
				return false
			}
		default:
			return true
		}
	}
	panic("unreachable")
}

// OnSysKey on Ctrl key pressed
func (e *Editor) OnSysKey(ev *termbox.Event) bool {
	switch ev.Key {
	// case termbox.KeyCtrlG:
	// 	//v := g.active.leaf
	// 	//v.ac = nil
	// 	//g.set_overlay_mode(nil)
	// 	//g.set_status("Quit")
	// 	e.msg("Quit")
	// 	e.Msgflag = true
	// case termbox.KeyCtrlZ:
	// 	//suspend(e)
	case termbox.KeyCtrlQ:
		e.Done = true
	}
	lookfor := fmt.Sprintf("%#2x %U", ev.Key, ev.Key)
	log.Println("OnSysKey is ", ev.Key, lookfor, e.KeysMapMap[lookfor])
	log.Println("OnSysKey is ", ev.Key, lookfor, e.KeysMapMap["\x04"])
	if k, ok := e.KeysMapMap[lookfor]; ok == true {
		// perform
		log.Printf("K %#v (%t)", k, ok)
		e.PerformAction(k.Do)
		return true
	}
	return false
}

func (e *Editor) PerformAction(fn func(*Editor)) {
	if fn != nil {
		fn(e)
	}
}

// OnAltKey on Alt key pressed
func (e *Editor) OnAltKey(ev *termbox.Event) bool {
	switch ev.Ch {
	case 'g':
		//g.set_overlay_mode(init_line_edit_mode(g, g.goto_line_lemp()))
		return true
		e.msg("Alt G")
		e.Msgflag = true
		return true
	case '/':
		//g.set_overlay_mode(init_autocomplete_mode(g))
		return true
		e.msg("Alt /")
		e.Msgflag = true
		return true
	case 'q':
		//g.set_overlay_mode(init_fill_region_mode(g))
		e.msg("Alt Q")
		e.Msgflag = true
		return true
	}
	return false
}

// OnKey some key
func (e *Editor) OnKey(ev *termbox.Event) {
	switch ev.Key {
	case termbox.KeyCtrlX:
		//g.set_overlay_mode(init_extended_mode(g))
	// case termbox.KeyCtrlS:
	// 	g.set_overlay_mode(init_isearch_mode(g, false))
	// case termbox.KeyCtrlR:
	// 	g.set_overlay_mode(init_isearch_mode(g, true))
	default:
		if ev.Mod&termbox.ModAlt != 0 && e.OnAltKey(ev) {
			break
		}
		e.CurrentWindow.OnKey(ev)
	}
}

func (e *Editor) msg(fm string, args ...interface{}) {
	e.Msgline = fmt.Sprintf(fm, args...)
	e.Msgflag = true
	return
}
func (e *Editor) drawstring(x, y int, fg, bg termbox.Attribute, msg string) {
	//log.Println(msg)
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (e *Editor) DisplayMsg() {
	//e.Cols, e.Lines = termbox.Size()
	if e.Msgflag {
		e.drawstring(0, e.Lines-1, e.FGColor, termbox.ColorDefault, e.Msgline)
	}
}

func (e *Editor) Display(wp *Window, flag bool) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	r, c := 0, 0
	idx := 0
	var rch rune
	bp := wp.Buffer

	// token_type := ID_DEFAULT

	// /* find start of screen, handle scroll up off page or top of file  */
	// /* Point is always within b_page and b_epage */
	if bp.Point() < bp.PageStart {
		bp.PageStart = e.SegStart(e.LineStart(bp.Point()), bp.Point())
	}

	// /* reframe when scrolled off bottom */
	if bp.Reframe == true || (bp.PageEnd <= bp.Point() && e.CurrentBuffer.Point() != e.CurrentBuffer.PageEnd) {
		bp.Reframe = true
		i := 0
		/* Find end of screen plus one. */
		bp.PageStart = e.DownDown(bp.Point())
		/* if we scoll to EOF we show 1 blank line at bottom of screen */
		if bp.BufferLen() <= bp.PageStart {
			bp.PageStart = bp.BufferLen()
			i = wp.Rows - 1
		} else {
			i = wp.Rows - 0
		}
		/* Scan backwards the required number of lines. */
		for 0 < i {
			bp.PageStart = e.UpUp(bp.PageStart)
			i--
		}

	}

	// move(wp.TopPt, 0); /* start from top of window */
	r = wp.TopPt
	c = 0
	// // bp.b_epage = bp.b_page;
	bp.PageEnd = bp.PageStart
	// // set_parse_state(bp, bp.b_epage); /* are we in a multline comment ? */

	// // /* paint screen from top of page until we hit maxline */
	for {
		// 	/* reached Point - store the Point position */
		if bp.Point() == bp.PageEnd {
			bp.PointRow = r
			bp.PointCol = c
		}
		// 	p = ptr(bp, bp.b_epage);
		// 	nch = 1;
		if r > wp.TopPt+wp.Rows || idx >= bp.BufferLen() { /* maxline */
			break
		}
		rch = bp.RuneAt(idx)
		//log.Println(rch, c, r)
		if rch != '\r' {
			if unicode.IsPrint(rch) || rch == '\t' || rch == '\n' {
				if rch == '\t' {
					c += 4 //? 8-(j&7) : 1;
				}
				termbox.SetCell(c, r, rch, e.FGColor, termbox.ColorDefault)
				c++

			} else {
				// const char *ctrl = unctrl(*p);
				// j += (int) strlen(ctrl);
				// addstr(ctrl); '\u2318'
				termbox.SetCell(c, r, '\uFFFD', e.FGColor, termbox.ColorDefault)
				c++
			}
		}

		if rch == '\n' || e.Cols <= c {
			c -= e.Cols
			if c < 0 {
				c = 0
			}
			r++
		}
		idx++
		// 	bp.b_epage = bp.b_epage + nch;
	}

	// /* replacement for clrtobot() to bottom of window */
	// for (k=i; k < wp.w_top + wp.w_rows; k++) {
	// 	move(k, j) /* clear from very last char not start of line */
	// 	clrtoeol()
	// 	j = 0 /* thereafter start of line */
	// }

	PushBuffer2Window(wp)
	e.ModeLine(wp)
	e.DisplayMsg()
	bp.PointCol, bp.PointRow = bp.XYForPoint(bp.Point())
	termbox.SetCursor(bp.PointCol-1, bp.PointRow-1)
	//termbox.Sync()
	wp.Updated = false
}

func (e *Editor) UpdateDisplay() {
	bp := e.CurrentWindow.Buffer
	bp.OrigPoint = bp.Point() /* OrigPoint only ever set here */

	/* only one window */
	if e.RootWindow.Next == nil {
		e.Display(e.CurrentWindow, true)
		//refresh()
		bp.PrevSize = bp.TextSize
		return
	}

	e.Display(e.CurrentWindow, false) /* this is key, we must call our win first to get accurate page and epage etc */

	/* never CurrentWin,  but same buffer in different window or update flag set*/
	for wp := e.RootWindow; wp != nil; wp = wp.Next {
		if wp != e.CurrentWindow && (wp.Buffer == bp || wp.Updated) {
			SyncBuffer(wp)
			e.Display(wp, false)
		}
	}

	/* now display our window and buffer */
	SyncBuffer(e.CurrentWindow)
	//e.DisplayMsg()
	e.CurrentWindow.CurCol, e.CurrentWindow.CurRow = e.CurrentBuffer.XYForPoint(bp.Point())
	termbox.SetCursor(e.CurrentWindow.CurRow-1, e.CurrentWindow.CurCol-1) /* set cursor for CurrentWin */
	//refresh()
	bp.PrevSize = bp.TextSize /* now safe to save previous size for next time */
}

func (e *Editor) ModeLine(wp *Window) {
	//i := 0
	var lch, mch, och rune
	e.Cols, e.Lines = termbox.Size()

	//standout();
	//move(wp.TopPt+wp.Rows, 0)
	// lch = (wp == CurrentWin ? '=' : '-')
	if wp == e.CurrentWindow {
		lch = '='
	} else {
		lch = '-'
	}
	// mch = ((wp.Buffer.Flags & B_MODIFIED) ? '*' : lch);
	mch = lch
	if wp.Buffer.modified {
		mch = '*'
	}
	// och = ((wp.Buffer.Flags & B_OVERWRITE) ? 'O' : lch);
	och = lch
	// if wp.Buffer.Flags&B_OVERWRITE != 0 {
	// 	och = 'O'
	// }

	temp := fmt.Sprintf("%c%c%c kg: %c%c %s (h %d, w%d) 0 y %d", lch, och, mch, lch, lch,
		e.GetBufferName(wp.Buffer), e.Lines, e.Cols, wp.TopPt+wp.Rows)
	//fmt.Println(temp)
	//e.drawstring(0, e.Lines-1, termbox.ColorWhite, termbox.ColorBlack, temp)
	x := 0
	y := wp.TopPt + wp.Rows
	//e.msg("win x %d y %d ", x, y)
	for _, c := range temp {
		termbox.SetCell(x, y, c, e.FGColor, e.BGColor)
		x++
	} //addstr(temp) // term

	for i := len(temp); i <= e.Cols; i++ {
		termbox.SetCell(i, y, lch, e.FGColor, e.BGColor)
	}
	//standend();
}

func (e *Editor) DisplayPromptAndResponse(prompt string, response string) {
	e.drawstring(0, e.Lines-1, e.FGColor, termbox.ColorDefault, prompt)
	/* if we have a value print it and go to end of it */
	if response != "" {
		e.drawstring(len(prompt), e.Lines-1, e.FGColor, termbox.ColorDefault, response)
	}

}

/* Reverse scan for start of logical line containing offset */
func (e *Editor) LineStart(off int) int {
	return e.CurrentBuffer.LineStart(off)
}

/* Forward scan for start of logical line segment (corresponds to screen line)  containing 'finish' */
func (e *Editor) SegStart(start int, finish int) int {
	bp := e.CurrentBuffer
	var p rune
	c := 0
	scan := start

	for scan < finish {
		//p = ptr(bp, scan);
		p = bp.RuneAt(scan)
		if p == '\n' {
			c = 0
			start = scan + 1
		} else {
			if e.Cols <= c {
				c = 0
				start = scan
			}
		}
		scan++
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if p == '\t' {
			c += 4 //8 - (c % 7)
		} else {
			c++
		}
	}
	// (c < COLS ? start : finish);
	if c < e.Cols {
		return start
	}
	return finish
}

/* Forward scan for start of logical line segment following 'finish' */
func (e *Editor) SegNext(start, finish int) int {
	// char_t *p;
	// int c = 0;
	bp := e.CurrentBuffer
	var p rune
	//var pptr int
	c := 0

	scan := e.SegStart(start, finish)
	for {
		//p = ptr(bp, scan);
		//p, pptr = bp.GetCurrentRune()
		p = bp.RuneAt(scan)
		//if (bp.b_ebuf <= p || COLS <= c)
		if e.Cols <= c {
			break
		}
		//scan += utf8_size(*ptr(bp,scan));
		scan++
		if p == '\n' {
			break
		}
		//c += *p == '\t' ? 8 - (c & 7) : 1;
		if p == '\t' {
			c += 4 //8 - (c % 7)
		} else {
			c++
		}
	}
	//(p < bp.b_ebuf ? scan : );
	if scan < bp.BufferLen() {
		return scan
	}
	return bp.BufferLen()
}

/* Move up one screen line */
func (e *Editor) UpUp(off int) int {
	curr := e.LineStart(off)
	seg := e.SegStart(curr, off)
	if curr < seg {
		off = e.SegStart(curr, seg-1)
	} else {
		off = e.SegStart(e.LineStart(curr-1), curr-1)
	}
	return off
}

/* Move down one screen line */
func (e *Editor) DownDown(off int) int {
	return (e.SegNext(e.LineStart(off), off))
}

// OffsetForColumn ln column - Return the offset of a column on the specified line
func (e *Editor) OffsetForColumn(offset int, column int) int {
	var p rune
	c := 0
	p = e.CurrentBuffer.RuneAt(offset)
	for offset < e.CurrentBuffer.PageEnd && p != '\n' && c < column {
		if p == '\t' {
			c += 4 //8 - (c % 7)
		} else {
			c++
		}
		offset++
		p = e.CurrentBuffer.RuneAt(offset)
	}
	return offset
}

// DeleteBuffer unlink from the list of buffers, free associated memory,
// assumes buffer has been saved if modified
func (e *Editor) DeleteBuffer(bp *Buffer) bool {
	//editor := bp.CurrentWindow.Editor
	var sb *Buffer

	/* we must have switched to a different buffer first */
	//assert(bp != CurrentBuffer)
	if bp != e.CurrentBuffer {
		/* if buffer is the head buffer */
		if bp == e.RootBuffer {
			e.RootBuffer = bp.Next
		} else {
			/* find place where the bp buffer is next */
			for sb = e.RootBuffer; sb.Next != bp && sb.Next != nil; sb = sb.Next {
			}
			if sb.Next == bp || sb.Next == nil {
				sb.Next = bp.Next
			}
		}

		/* now we can delete */
		//free(bp.BufferStart);
		//bp.BufferStart = nil
		//free(bp);
		bp = nil
	} else {
		return false
	}
	return true
}

// NextBuffer returns next buffer after current
func (e *Editor) NextBuffer() {
	if e.CurrentBuffer != nil && e.RootBuffer != nil {
		e.CurrentWindow.DisassociateBuffer()
		if e.CurrentBuffer.Next != nil {
			e.CurrentBuffer = e.CurrentBuffer.Next

		} else {
			e.CurrentBuffer = e.RootBuffer
		}
		e.CurrentWindow.AssociateBuffer(e.CurrentBuffer)
	}
}

// GetBufferName returns buffer name
func (e *Editor) GetBufferName(bp *Buffer) string {
	if bp.Filename != "" {
		return bp.Filename
	}
	return bp.Buffername
}

// CountBuffers how many buffers in list
func (e *Editor) CountBuffers() int {
	var bp *Buffer
	i := 0

	for bp = e.RootBuffer; bp != nil; bp = bp.Next {
		i++
	}
	return i
}

// ModifiedBuffers true is any buffers modified
func (e *Editor) ModifiedBuffers() bool {
	var bp *Buffer

	for bp = e.RootBuffer; bp != nil; bp = bp.Next {
		if bp.modified == true {
			return true
		}
	}
	return false
}

/* Buffer lists manipulation */
/* Find a buffer by filename or create if requested */
func (e *Editor) FindBuffer(fname string, cflag bool) *Buffer {
	var bp *Buffer
	var sb *Buffer

	bp = e.RootBuffer
	for bp != nil {
		if strings.Compare(fname, bp.Filename) == 0 || strings.Compare(fname, bp.Buffername) == 0 {
			return bp
		}
		bp = bp.Next
	}

	if cflag != false {
		// if ((bp = (buffer_t *) malloc (sizeof (buffer_t))) == nil)
		// 	return (0);
		bp = NewBuffer()

		//BufferInit(bp)
		//assert(bp != nil);

		/* find the place in the list to insert this buffer */
		if e.RootBuffer == nil {
			e.RootBuffer = bp
		} else if strings.Compare(e.RootBuffer.Filename, fname) > 0 {
			/* insert at the begining */
			bp.Next = e.RootBuffer
			e.RootBuffer = bp
		} else {
			for sb = e.RootBuffer; sb.Next != nil; sb = sb.Next {
				if strings.Compare(sb.Next.Filename, fname) > 0 {
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

func (e *Editor) SplitWindow() {
	var editor = e
	var wp *Window
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
func (e *Editor) NextWindow() {
	var editor = e
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
func (e *Editor) DeleteOtherWindows() {
	wp := e.RootWindow
	if wp.Next == nil {
		wp.Editor.msg("Only 1 window")
		return
	}
	e.FreeOtherWindows()
}

// FreeOtherWindows
func (e *Editor) FreeOtherWindows() {
	var editor = e
	var winp *Window
	var wp *Window
	var next *Window
	wp = e.RootWindow
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
