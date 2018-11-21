package tkg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	termbox "github.com/nsf/termbox-go"
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
	CHUNK            = 16 //= 8096
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

type Editor struct {
	EventChan     chan termbox.Event
	MiniBufChan   chan termbox.Event
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
	Lines         int
	Cols          int
	FGColor       termbox.Attribute
	BGColor       termbox.Attribute
	EscapeFlag    bool
	CtrlXFlag     bool
	MiniBufActive bool
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
		/* Save filename regardless of load() success. */
		e.CurrentBuffer.Filename = argv[1]
	} else {
		e.msg("NO file to open, creating scratch buffer")
		e.CurrentBuffer = e.FindBuffer("*scratch*", true)
		e.CurrentBuffer.Buffername = "*scratch*"
		s := "Lorem ipsum dolor sit amet,\nconsectetur adipiscing elit,\nsed do eiusmod tempor incididunt ut\nlabore et dolore magna aliqua. \n"

		e.CurrentBuffer.SetText(s)
		e.CurrentBuffer.Insert(s)
		e.CurrentBuffer.Insert(" foo")
		e.CurrentBuffer.Insert(s)
		e.CurrentBuffer.Insert(" baz\n")
		e.CurrentBuffer.Insert(s)
		// e.CurrentBuffer.Insert(s)
		// e.CurrentBuffer.Insert(s)
		// e.CurrentBuffer.Insert(s)
		// e.CurrentBuffer.Insert(s)
		// e.CurrentBuffer.Insert(s)
		// e.CurrentBuffer.Insert(s)
		// e.CurrentBuffer.Insert(s)
		e.bottom()
		e.CurrentBuffer.Insert("\nbottom End Of File.\n")
		e.top()
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
	termbox.SetInputMode(termbox.InputAlt | termbox.InputEsc | termbox.InputMouse)
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
			//log.Printf("%#v\n", ev)
			log.Println(">>\n ", time.Now().Unix(), "\n>>")
			// if e.MiniBufActive {
			// 	log.Printf("MB %#v\n", ev)
			// 	e.MiniBufChan <- ev
			// 	if ev.Key == termbox.KeyCtrlQ {
			// 		return
			// 	}
			// } else {
			// log.Printf("Normal %#v\n", ev)
			ok := e.HandleEvent(&ev)
			if !ok {
				return
			}
			//e.ConsumeMoreEvents()
			e.UpdateDisplay()
			termbox.Flush()
			// }
		}
	}
	return

}

// HandleEvent
func (e *Editor) HandleEvent(ev *termbox.Event) bool {
	e.msg("")
	switch ev.Type {
	case termbox.EventKey:
		if ev.Ch != 0 && (e.CtrlXFlag || e.EscapeFlag) {
			_ = e.OnSysKey(ev)
			if e.Done {
				return false
			}
		} else if ev.Ch == 0 {
			_ = e.OnSysKey(ev)
			if e.Done {
				return false
			}
		} else {
			if ev.Mod&termbox.ModAlt != 0 && e.OnAltKey(ev) {
				break
			}
			e.CurrentWindow.OnKey(ev)
		}
		e.UpdateDisplay()
	case termbox.EventResize:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		e.Cols, e.Lines = termbox.Size()
		e.msg("Resize: h %d,w %d", e.Lines, e.Cols)
		e.CurrentWindow.WindowResize()
		e.UpdateDisplay()
	case termbox.EventMouse:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		e.msg("Mouse: %d,%d [%d %d] not impl.", ev.MouseX, ev.MouseY, e.Cols, e.Lines)
		// Need a hit test algo here
		//pt := e.CurrentWindow.PointForXY(ev.MouseX+1, ev.MouseY+1)
		//e.CurrentBuffer.SetPoint(pt)
		e.UpdateDisplay()
	case termbox.EventError:
		panic(ev.Err)
	}

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
	case termbox.KeyCtrlX:
		e.msg("C-X ")
		e.CtrlXFlag = true
		return true
	case termbox.KeyEsc:
		e.msg("Esc ")
		e.EscapeFlag = true
		return true
	case termbox.KeyCtrlQ:
		e.Done = true
		return true
	case termbox.KeySpace, termbox.KeyEnter, termbox.KeyCtrlJ, termbox.KeyTab:
		e.CurrentWindow.OnKey(ev)
		return true
	case termbox.KeyArrowDown, termbox.KeyArrowLeft, termbox.KeyArrowRight, termbox.KeyArrowUp:
		e.CtrlXFlag = false
		e.EscapeFlag = false
		return e.searchAndPerform(ev)
	default:
		return e.searchAndPerform(ev)
	}
	return false
}

func (e *Editor) searchAndPerform(ev *termbox.Event) bool {
	lookfor := ""
	rch := ev.Ch
	if ev.Ch == 0 {
		rch = rune(ev.Key)
	}
	lookfor = fmt.Sprintf("%c", rch)
	if e.CtrlXFlag {
		log.Printf("lookfor %c\n", rch)
		lookfor = fmt.Sprintf("\x18%c", rch)
	}
	if e.EscapeFlag {
		lookfor = fmt.Sprintf("\x1B%c", rch)
	}
	for i, j := range e.Keymap {
		if strings.Compare(lookfor, j.KeyBytes) == 0 {
			//log.Println("SearchAndPerform FOUND ", lookfor, e.Keymap[i])
			jj := e.Keymap[i].Do
			e.performAction(jj)
			e.CtrlXFlag = false
			e.EscapeFlag = false
			return true
		}
	}
	return false
}
func (e *Editor) performAction(fn actionFunc) {
	if fn != nil {
		fn(e)
	}
}

// OnAltKey on Alt key pressed
func (e *Editor) OnAltKey(ev *termbox.Event) bool {
	e.msg("AltKey\n")
	// switch ev.Ch {
	// case 'g':
	// 	//g.set_overlay_mode(init_line_edit_mode(g, g.goto_line_lemp()))
	// 	return true
	// 	e.msg("Alt G")
	// 	e.Msgflag = true
	// 	return true
	// case '/':
	// 	//g.set_overlay_mode(init_autocomplete_mode(g))
	// 	return true
	// 	e.msg("Alt /")
	// 	e.Msgflag = true
	// 	return true
	// case 'q':
	// 	//g.set_overlay_mode(init_fill_region_mode(g))
	// 	e.msg("Alt Q")
	// 	e.Msgflag = true
	// 	return true
	// }
	return false
}

// OnKey some key
func (e *Editor) OnKey(ev *termbox.Event) {
	// switch ev.Key {
	// case termbox.KeyCtrlX:
	// 	//g.set_overlay_mode(init_extended_mode(g))
	// // case termbox.KeyCtrlS:
	// // 	g.set_overlay_mode(init_isearch_mode(g, false))
	// // case termbox.KeyCtrlR:
	// // 	g.set_overlay_mode(init_isearch_mode(g, true))
	// default:

	// }
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

// Display draws the window, minding the buffer pagestart/pageend
func (e *Editor) Display(wp *Window, flag bool) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	bp := wp.Buffer
	// /* find start of screen, handle scroll up off page or top of file  */
	// /* Point is always within b_page and b_epage */
	if bp.Point() < bp.PageStart {
		bp.PageStart = bp.SegStart(bp.LineStart(bp.Point()), bp.Point(), e.Cols)
	}

	if bp.Reframe == true || (bp.Point() > bp.PageEnd && bp.Point() != bp.PageEnd) {
		bp.Reframe = false
		i := 0
		/* Find end of screen plus one. */
		bp.PageStart = bp.DownDown(bp.Point(), e.Cols)
		//log.Printf("P1 PageStart %d Point %d, bp.PageEnd %d", bp.PageStart, bp.Point(), bp.PageEnd)
		/* if we scoll to EOF we show 1 blank line at bottom of screen */
		if bp.PageEnd <= bp.PageStart {
			bp.PageStart = bp.PageEnd
			i = wp.Rows - 1 // 1
		} else {
			i = wp.Rows - 0
		}
		/* Scan backwards the required number of lines. */
		//log.Printf("Before BWscan i %d PageStart %d Point %d, bp.PageEnd %d", i, bp.PageStart, bp.Point(), bp.PageEnd)
		for i > 0 {
			bp.PageStart = bp.UpUp(bp.PageStart, e.Cols)
			i--
			//log.Printf("P3 i %d PageStart %d Point %d, bp.PageEnd %d", i, bp.PageStart, bp.Point(), bp.PageEnd)
		}
	}

	l1 := bp.LineForPoint(bp.PageStart)
	l2 := l1 + wp.Rows
	l2end := bp.LineEnd(bp.PointForLine(l2))
	bp.PageEnd = l2end
	log.Printf("P0 lines %d %d PageStart %d Point %d, bp.PageEnd %d BufL %d", l1, l2, bp.PageStart, bp.Point(), bp.PageEnd, bp.BufferLen())
	//toPrint := bp.GetTextForLines(l1, l2+1)
	//runeArray := []rune(toPrint)
	r, c := 0, 0
	for k := bp.PageStart; k <= bp.PageEnd; k++ {
		/* reached point - store the cursor position */
		if bp.Point() == k {
			bp.PointCol = c
			bp.PointRow = r
		}
		//rch := runeArray[k]
		rch, err := bp.RuneAt(k)
		if err != nil {
			log.Println("Error on RuneAt", err)
		}
		if rch != '\r' {
			if unicode.IsPrint(rch) || rch == '\t' || rch == '\n' {
				if rch == '\t' {
					c += 3 //? 8-(j&7) : 1;
				}
				termbox.SetCell(c, r, rch, e.FGColor, termbox.ColorDefault)
				c++
			} else {
				termbox.SetCell(c, r, rch, e.FGColor, termbox.ColorDefault)
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
	}

	// /* replacement for clrtobot() to bottom of window */
	// for k := idx; k < wp.TopPt+wp.Rows; k++ {
	// 	//move(k, j) /* clear from very last char not start of line */
	// 	//clrtoeol()
	// 	for cc := 0; cc < e.Cols; cc++ {
	// 		termbox.SetCell(cc, k, ' ', e.FGColor, termbox.ColorDefault)
	// 	}
	// 	//sc = 0 /* thereafter start of line */
	// }

	PushBuffer2Window(wp)
	e.ModeLine(wp)
	if wp == e.CurrentWindow && flag {
		e.DisplayMsg()
		//move(bp->b_row, bp->b_col); /* set cursor */
		e.SetTermCursor(bp.PointCol, bp.PointRow)
		termbox.Flush() //refresh();
	}
	wp.Updated = false
	// b2w(wp); /* save buffer stuff on window */
	// modeline(wp);
	// if (wp == curwp && flag) {
	// 	dispmsg();
	// 	move(bp->b_row, bp->b_col); /* set cursor */
	// 	refresh();
	// }
	// wp->w_update = FALSE;
}

func (e *Editor) UpdateDisplay() {
	bp := e.CurrentWindow.Buffer
	bp.OrigPoint = bp.Point() /* OrigPoint only ever set here */

	/* only one window */
	if e.RootWindow.Next == nil {
		e.Display(e.CurrentWindow, true)
		bp.PrevSize = bp.TextSize
		return
	}

	e.Display(e.CurrentWindow, false)
	/* this is key, we must call our win first to get accurate page and epage etc */

	/* never CurrentWin,  but same buffer in different window or update flag set*/
	for wp := e.RootWindow; wp != nil; wp = wp.Next {
		if wp != e.CurrentWindow && (wp.Buffer == bp || wp.Updated) {
			SyncBuffer(wp)
			e.Display(wp, false)
		}
	}

	/* now display our window and buffer */
	SyncBuffer(e.CurrentWindow)
	e.DisplayMsg()
	bp.PrevSize = bp.TextSize /* now safe to save previous size for next time */
}

// SetTermCursor -
func (e *Editor) SetTermCursor(c, r int) {
	wp := e.CurrentWindow
	//log.Println("wp t,p", wp.TopPt, wp.Rows)
	if r > wp.Rows {
		r = wp.Rows + 1
	}
	pt := wp.Buffer.Point()
	wp.Buffer.logBufferEOB(pt)
	wp.CurCol, wp.CurRow = c, r
	termbox.SetCursor(c, r)
	if wp.Buffer.EndOfBuffer(pt) == true {
		eol := wp.Buffer.ColumnForPoint(wp.Buffer.LineEnd(pt)) - 1
		wp.CurCol, wp.CurRow = eol, r-1
		termbox.SetCursor(eol, r-1)
	}
	/* set cursor for CurrentWin */
}

// ModeLine draw modeline for window
func (e *Editor) ModeLine(wp *Window) {
	var lch, mch, och rune
	e.Cols, e.Lines = termbox.Size()

	if wp == e.CurrentWindow {
		lch = '='
	} else {
		lch = '-'
	}
	mch = lch
	if wp.Buffer.modified {
		mch = '*'
	}
	och = lch
	temp := fmt.Sprintf("%c%c%c kg: %c%c %s wp(%d,%d) (h %d, w%d) rows %d", lch, och, mch, lch, lch,
		e.GetBufferName(wp.Buffer),
		wp.CurCol, wp.CurRow,
		//c, r,
		e.Lines, e.Cols, wp.TopPt+wp.Rows)
	x := 0
	y := wp.TopPt + wp.Rows + 1
	for _, c := range temp {
		termbox.SetCell(x, y, c, termbox.ColorBlack, e.BGColor)
		//mch = c
		x++
	}

	for i := len(temp); i <= e.Cols; i++ {
		termbox.SetCell(i, y, lch, termbox.ColorBlack, e.BGColor) // e.FGColor
	}
}

func (e *Editor) DisplayPromptAndResponse(prompt string, response string) {
	e.drawstring(0, e.Lines-1, e.FGColor, termbox.ColorDefault, prompt)
	/* if we have a value print it and go to end of it */
	if response != "" {
		e.drawstring(len(prompt), e.Lines-1, e.FGColor, termbox.ColorDefault, response)
	}
	termbox.SetCursor(len(prompt)+len(response)+1, e.Lines-1)
	termbox.Flush()
}

func (e *Editor) GetFilename(prompt string) string {
	fname := ""
	var ev termbox.Event
	e.DisplayPromptAndResponse(prompt, "")
	e.MiniBufActive = true
	for {
		ev = <-e.MiniBufChan
		ch := ev.Ch
		if ch == '\r' || ch == '\n' {
			break
		}
		fname = fname + string(ch)
		e.DisplayPromptAndResponse(prompt, fname)
	}
	e.MiniBufActive = false
	return fname
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

// attempt for new PointUp
func (e *Editor) PointUp() {
	// urbp->b_point = lncolumn(curbp, upup(curbp, curbp->b_point),curbp->b_col);
	bp := e.CurrentBuffer
	pt := bp.Point()
	c1 := bp.ColumnForPoint(pt)
	npt := bp.UpUp(pt, c1)
	bp.SetPointAndCursor(npt)
}
