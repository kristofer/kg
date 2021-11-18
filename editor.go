package kg

import (
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

const (
	version        = "kg 1.1, Public Domain, November 2021, Kristofer Younger,  No warranty."
	nomark         = -1
	gapchunk       = 16 //= 8096
	idDefault      = 1
	idSymbol       = 2
	idModeline     = 3
	idDigits       = 4
	idLineComment  = 5
	idBlockComment = 6
	idDoubleString = 7
	idSingleString = 8
)

// Editor struct
type Editor struct {
	EventChan     chan termbox.Event
	CurrentBuffer *Buffer /* current buffer */
	RootBuffer    *Buffer /* head of list of buffers */
	CurrentWindow *Window
	RootWindow    *Window
	// status vars
	Done          bool   /* Quit flag. */
	Msgflag       bool   /* True if msgline should be displayed. */
	PasteBuffer   string /* Allocated scrap buffer. */
	Msgline       string /* Message line input/output buffer. */
	Searchtext    string
	Replace       string
	Keymap        []keymapt
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
	checkErr(err)
	defer termbox.Close()
	e.Cols, e.Lines = termbox.Size()

	if argc > 1 {
		for k := 1; k < argc; k++ {
			e.CurrentBuffer = e.FindBuffer(argv[k], true)
			e.InsertFile(argv[k], false)
			/* Save filename regardless of load() success. */
			e.CurrentBuffer.Filename = argv[k]
		}
	} else {
		e.msg("NO file to open, creating scratch buffer")
		e.CurrentBuffer = e.FindBuffer("*scratch*", true)
		e.CurrentBuffer.Buffername = "*scratch*"
		e.top()
	}
	e.CurrentWindow = NewWindow(e)
	e.RootWindow = e.CurrentWindow
	e.CurrentWindow.OneWindow()
	e.CurrentWindow.AssociateBuffer(e.CurrentBuffer)

	if !(e.CurrentBuffer.GrowGap(gapchunk)) {
		panic("%s: Failed to allocate required memory.\n")
	}
	e.Keymap = keymap
	termbox.SetInputMode(termbox.InputAlt | termbox.InputEsc | termbox.InputMouse)
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	e.updateDisplay()
	termbox.Flush()

	e.EventChan = make(chan termbox.Event, 20)
	go func() {
		for {
			e.EventChan <- termbox.PollEvent()
		}
	}()

	// Instead of using for {
	// 	select {
	// 	case ev := <-e.EventChan:

	for ev := range e.EventChan {
		ok := e.handleEvent(&ev)
		if !ok {
			return
		}
		e.updateDisplay()
		termbox.Flush()
	}
}

// handleEvent
func (e *Editor) handleEvent(ev *termbox.Event) bool {
	e.msg("")
	switch ev.Type {
	case termbox.EventKey:
		if (ev.Mod & termbox.ModAlt) != 0 {
			e.msg("FOUND ALT...")
			switch ev.Ch {
			case 'j':
				e.msg("FOUND ALT J")
			// and others..
			default:
			}
		}
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
			// if ev.Mod&termbox.ModAlt != 0 && e.OnAltKey(ev) {
			// 	break
			// }
			if (ev.Mod & termbox.ModAlt) != 0 {
				switch ev.Ch {
				case 'j':
					e.msg("FOUND ALT J")
				// and others..
				default:
				}
			}
			e.CurrentWindow.OnKey(ev)
		}
		e.updateDisplay()
	case termbox.EventResize:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		e.Cols, e.Lines = termbox.Size()
		e.msg("Resize: h %d,w %d", e.Lines, e.Cols)
		e.CurrentWindow.WindowResize()
		e.updateDisplay()
	case termbox.EventMouse:
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		e.msg("Mouse: r %d, c %d ", ev.MouseY, ev.MouseX)
		e.SetPointForMouse(ev.MouseX, ev.MouseY)
		e.updateDisplay()
	case termbox.EventError:
		panic(ev.Err)
	}

	return true
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
}

func (e *Editor) searchAndPerform(ev *termbox.Event) bool {
	rch := ev.Ch
	if ev.Ch == 0 {
		rch = rune(ev.Key)
	}
	lookfor := fmt.Sprintf("%c", rch)
	if e.CtrlXFlag {
		lookfor = fmt.Sprintf("\x18%c", rch)
	}
	if e.EscapeFlag {
		lookfor = fmt.Sprintf("\x1B%c", rch)
	}
	for i, j := range e.Keymap {
		if strings.Compare(lookfor, j.KeyBytes) == 0 {
			//// log.Println("SearchAndPerform FOUND ", lookfor, e.Keymap[i])
			do := e.Keymap[i].Do
			if do != nil {
				do(e) // execute function for key
			}
			e.CtrlXFlag = false
			e.EscapeFlag = false
			return true
		}
	}
	return false
}

// OnAltKey on Alt key pressed
func (e *Editor) OnAltKey(ev *termbox.Event) bool {
	e.msg("AltKey")
	return false
}

func (e *Editor) msg(fm string, args ...interface{}) {
	e.Msgline = fmt.Sprintf(fm, args...)
	e.Msgflag = true
}

func (e *Editor) drawString(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (e *Editor) displayMsg() {
	if e.Msgflag {
		e.drawString(0, e.Lines-1, e.FGColor, termbox.ColorDefault, e.Msgline)
	}
	e.blankFrom(e.Lines-1, len(e.Msgline))
}

// Display draws the window, minding the buffer pagestart/pageend
func (e *Editor) Display(wp *Window, shouldDrawCursor bool) {
	bp := wp.Buffer
	pt := bp.Point
	// /* find start of screen, handle scroll up off page or top of file  */
	if pt < bp.PageStart {
		bp.PageStart = bp.SegStart(bp.LineStart(pt), pt, e.Cols)
	}

	if bp.Reframe || (pt > bp.PageEnd && pt != bp.PageEnd && !(pt >= bp.TextSize)) {
		bp.Reframe = false
		i := 0
		/* Find end of screen plus one. */
		bp.PageStart = bp.DownDown(pt, e.Cols)
		/* if we scroll to EOF we show 1 blank line at bottom of screen */
		if bp.PageEnd <= bp.PageStart {
			bp.PageStart = bp.PageEnd
			i = wp.Rows - 1 // 1
		} else {
			i = wp.Rows - 0
		}
		/* Scan backwards the required number of lines. */
		for i > 0 {
			bp.PageStart = bp.UpUp(bp.PageStart, e.Cols)
			i--
		}
	}

	l1 := bp.LineForPoint(bp.PageStart)
	l2 := l1 + wp.Rows
	l2end := bp.LineEnd(bp.PointForLine(l2))
	bp.PageEnd = l2end
	r, c := wp.TopPt, 0
	for k := bp.PageStart; k <= bp.PageEnd; k++ {
		if pt == k {
			bp.PointCol = c
			bp.PointRow = r
		}
		rch, err := bp.RuneAt(k)
		if err != nil {
			e.msg("Error on RuneAt", err)
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
			e.blankFrom(r, c)
			c -= e.Cols
			if c < 0 {
				c = 0
			}
			r++
		}
	}
	for k := r; k < wp.TopPt+wp.Rows+1; k++ {
		e.blankFrom(k, 0)
	}

	buffer2Window(wp)
	e.ModeLine(wp)
	if wp == e.CurrentWindow && shouldDrawCursor {
		e.displayMsg()
		e.setTermCursor(wp.Col, wp.Row) //bp.PointCol, bp.PointRow)
	}
	termbox.Flush() //refresh();
	wp.Updated = false
}

func (e *Editor) blankFrom(r, c int) { // blank line to end of term
	for k := c; k < (e.Cols - 1); k++ {
		termbox.SetCell(k, r, ' ', e.FGColor, termbox.ColorDefault)
	}
}
func (e *Editor) setTermCursor(c, r int) {
	wp := e.CurrentWindow
	wp.Col, wp.Row = c, r
	termbox.SetCursor(c, r)
	//// log.Printf("c %d r %d\n", c, r)
}

func (e *Editor) updateDisplay() {
	bp := e.CurrentWindow.Buffer
	bp.OrigPoint = bp.Point /* OrigPoint only ever set here */
	/* only one window */
	if e.RootWindow.Next == nil {
		e.Display(e.CurrentWindow, true)
		termbox.Flush()
		bp.PrevSize = bp.TextSize
		return
	}
	/* this is key, we must call our win first to get accurate page and epage etc */
	e.Display(e.CurrentWindow, false)
	/* never CurrentWin,  but same buffer in different window or update flag set*/
	for wp := e.RootWindow; wp != nil; wp = wp.Next {
		if wp != e.CurrentWindow && (wp.Buffer == bp || wp.Updated) {
			window2Buffer(wp)
			e.Display(wp, false)
		}
	}
	/* now display our window and buffer */
	window2Buffer(e.CurrentWindow)
	e.displayMsg()
	e.setTermCursor(e.CurrentWindow.Col, e.CurrentWindow.Row)
	bp.PrevSize = bp.TextSize /* now safe to save previous size for next time */
}

// SetPointForMouse xxx
func (e *Editor) SetPointForMouse(mc, mr int) {
	c, r := e.setWindowForMouse(mc, mr)
	bp := e.CurrentBuffer
	sl := bp.LineForPoint(bp.PageStart) // sl is startline of buffer frame
	ml := sl + r
	mlpt := bp.PointForLine(ml)
	mll := bp.LineLenAtPoint(mlpt) // how wide is line?
	nc := c + 1
	if mll < c {
		nc = mll
	}
	npt := bp.PointForXY(nc, ml)
	bp.SetPoint(npt)
}

func (e *Editor) setWindowForMouse(mc, mr int) (c, r int) {
	log.Printf("setWindowForMouse col %d row %d ", mc, mr)

	wp := e.RootWindow
	// if mr is modeline or modeline+1, reduce to last wp.Rows
	if mr > wp.Rows {
		mr = wp.Rows
	}
	for wp != nil {
		if (mr <= wp.Rows+wp.TopPt) && (mr >= wp.TopPt) {
			log.Printf("set win rows %d top %d\n", wp.Rows, wp.TopPt)
			e.setWindow(wp)
			r = mr - wp.TopPt
			// if mr == wp.Rows+wp.TopPt {
			// 	r--
			// }
			c = mc
			return
		}
		wp = wp.Next
	}
	return 0, e.Lines - 1
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
	temp := fmt.Sprintf("%c%c%c kg: %c%c %s L%d wp(%d,%d)", lch, och, mch, lch, lch,
		e.GetBufferName(wp.Buffer),
		wp.Buffer.PointRow, wp.Row, wp.Col)
	x := 0
	y := wp.TopPt + wp.Rows + 1
	for _, c := range temp {
		termbox.SetCell(x, y, c, termbox.ColorBlack, e.BGColor)
		x++
	}

	for i := len(temp); i <= e.Cols; i++ {
		termbox.SetCell(i, y, lch, termbox.ColorBlack, e.BGColor) // e.FGColor
	}
}

func (e *Editor) displayPromptAndResponse(prompt string, response string) {
	e.drawString(0, e.Lines-1, e.FGColor, termbox.ColorDefault, prompt)
	if response != "" {
		e.drawString(len(prompt), e.Lines-1, e.FGColor, termbox.ColorDefault, response)
	}
	e.blankFrom(e.Lines-1, len(prompt)+len(response))
	termbox.SetCursor(len(prompt)+len(response), e.Lines-1)
	termbox.Flush()
}

func (e *Editor) getInput(prompt string) string {
	fname := ""
	var ev termbox.Event
	e.displayPromptAndResponse(prompt, "")
	e.MiniBufActive = true
loop:
	for {
		ev = <-e.EventChan
		if ev.Ch != 0 {
			ch := ev.Ch
			fname = fname + string(ch)
		}
		if ev.Ch == 0 {
			switch ev.Key {
			case termbox.KeyTab:
				fname = fname + string('\t')
			case termbox.KeySpace:
				fname = fname + string(' ')
			case termbox.KeyEnter, termbox.KeyCtrlR:
				break loop
			case termbox.KeyBackspace2, termbox.KeyBackspace:
				if len(fname) > 0 {
					fname = fname[:len(fname)-1]
				} else {
					fname = ""
				}
			case termbox.KeyCtrlG:
				return ""
			default:

			}
		}
		e.displayPromptAndResponse(prompt, fname)
	}
	e.MiniBufActive = false
	return fname
}

// DeleteBuffer unlink from the list of buffers, free associated memory,
// assumes buffer has been saved if modified
func (e *Editor) deleteBuffer(bp *Buffer) bool {
	var sb *Buffer

	/* we must have switched to a different buffer first */
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
		bp = nil
	} else {
		return false
	}
	return true
}

// NextBuffer returns next buffer after current
func (e *Editor) nextBuffer() {
	if e.CurrentBuffer != nil && e.RootBuffer != nil {
		e.CurrentWindow.DisassociateBuffer()
		if e.CurrentBuffer.Next != nil {
			e.CurrentBuffer = e.CurrentBuffer.Next

		} else {
			e.CurrentBuffer = e.RootBuffer
		}
		e.CurrentWindow.AssociateBuffer(e.CurrentBuffer)
		e.CurrentBuffer.Reframe = true
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
		if bp.modified {
			return true
		}
	}
	return false
}

// FindBuffer Find a buffer by filename or create if requested
func (e *Editor) FindBuffer(fname string, cflag bool) *Buffer {
	bp := e.RootBuffer
	for bp != nil {
		if strings.Compare(fname, bp.Filename) == 0 || strings.Compare(fname, bp.Buffername) == 0 {
			return bp
		}
		bp = bp.Next
	}
	if cflag {
		bp = NewBuffer()
		/* find the place in the list to insert this buffer */
		if e.RootBuffer == nil {
			e.RootBuffer = bp
		} else if strings.Compare(e.RootBuffer.Filename, fname) > 0 {
			/* insert at the begining */
			bp.Next = e.RootBuffer
			e.RootBuffer = bp
		} else {
			sb := e.RootBuffer
			for sb.Next != nil {
				if strings.Compare(sb.Next.Filename, fname) > 0 {
					break
				}
				sb = sb.Next
			}
			/* and insert it */
			bp.Next = sb.Next
			sb.Next = bp
		}
	}
	return bp
}

func (e *Editor) splitWindow() {
	if e.CurrentWindow.Rows < 3 {
		e.msg("Cannot split a %d line window", e.CurrentWindow.Rows)
		return
	}

	nwp := NewWindow(e)
	nwp.AssociateBuffer(e.CurrentWindow.Buffer)
	buffer2Window(nwp)

	ntru := (e.CurrentWindow.Rows - 1) / 2    /* Upper size */
	ntrl := (e.CurrentWindow.Rows - 1) - ntru /* Lower size */

	/* Old is upper window */
	e.CurrentWindow.Rows = ntru
	nwp.TopPt = e.CurrentWindow.TopPt + ntru + 2
	nwp.Rows = ntrl - 1

	/* insert it in the list */
	wp2 := e.CurrentWindow.Next
	e.CurrentWindow.Next = nwp
	nwp.Next = wp2
	/* mark the lot for update */
	e.redraw()
}

// NextWindow
func (e *Editor) nextWindow() {
	e.CurrentWindow.Updated = true /* make sure modeline gets updated */
	//Curwp = (Curwp.Next == nil ? Wheadp : Curwp.Next)
	if e.CurrentWindow.Next == nil {
		e.CurrentWindow = e.RootWindow
	} else {
		e.CurrentWindow = e.CurrentWindow.Next
	}
	e.CurrentBuffer = e.CurrentWindow.Buffer

	if e.CurrentBuffer.WinCount > 1 {
		/* push win vars to buffer */
		window2Buffer(e.CurrentWindow)
	}
}

func (e *Editor) setWindow(wp *Window) {
	e.CurrentWindow.Updated = true /* make sure modeline gets updated */
	e.CurrentWindow = wp
	e.CurrentWindow.Updated = true /* make sure modeline gets updated */
	e.CurrentBuffer = e.CurrentWindow.Buffer
	// if e.CurrentBuffer.WinCount > 1 {
	// 	/* push win vars to buffer */
	// 	window2Buffer(e.CurrentWindow)
	// }
	e.updateDisplay()
}

// DeleteOtherWindows
func (e *Editor) deleteOtherWindows() {
	wp := e.RootWindow
	if wp.Next == nil {
		wp.Editor.msg("Only 1 window")
		return
	}
	e.freeOtherWindows()
}

// FreeOtherWindows
func (e *Editor) freeOtherWindows() {
	wp := e.RootWindow
	winp := e.CurrentWindow
	next := wp
	for next != nil {
		next = wp.Next /* get next before a call to free() makes wp undefined */
		if wp != winp {
			wp.DisassociateBuffer() /* this window no longer references its buffer */
		}
		wp = next
	}
	e.RootWindow = winp
	e.CurrentWindow = winp
	winp.OneWindow()
}
