package tkg

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

var LINES = 1
var COLS = 10
var MSGLINE = (LINES - 1)

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

var ()

type Keymapt struct {
	KeyDesc  string
	KeyBytes string
	Do       *func() // function to call for Keymap-ping
}

type Editor struct {
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
	Key_map    *Keymapt /* Command key mappings. */
	Keymap     []Keymapt
	Key_return *Keymapt /* Command key return */
	//
	Lines int
	Cols  int
}

// StartEditor is the old C main function
func (e *Editor) StartEditor(argv []string, argc int) {

	if argc > 1 {
		e.CurrentBuffer = FindBuffer(argv[1], true)
		e.InsertFile(argv[1], false)
		/* Save filename irregardless of load() success. */
		//strncpy(e.CurrentBuffer->b_fname, argv[1], NAME_MAX);
		e.CurrentBuffer.Filename = argv[1]
		//e.CurrentBuffer->b_fname[NAME_MAX] = '\0'; /* force truncation */
	} else {
		e.CurrentBuffer = FindBuffer("*scratch*", true)
		//strncpy(e.CurrentBuffer->b_bname, "*scratch*", STRBUF_S);
		e.CurrentBuffer.Buffername = "*scratch*"
	}
	e.CurrentWindow = NewWindow(e)
	e.RootWindow = e.CurrentWindow
	e.CurrentWindow.OneWindow()
	e.CurrentWindow.AssociateBuffer(e.CurrentBuffer)

	if !(e.CurrentBuffer.GrowGap(CHUNK)) {
		panic("%s: Failed to allocate required memory.\n")
	}
	//e.Key_map = e.Keymap

	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	e.Lines, e.Cols = termbox.Size()

	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	//draw_keyboard()
	e.Display()
	termbox.Flush()
	inputmode := 0
	ctrlxpressed := false
loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			//log.Println("[", e.Lines, e.Cols, ev, "event", ctrlxpressed)
			if ev.Key == termbox.KeyCtrlS && ctrlxpressed {
				termbox.Sync()
			}
			if ev.Key == termbox.KeyCtrlQ && ctrlxpressed {
				break loop
			}
			if ev.Key == termbox.KeyCtrlC && ctrlxpressed {
				chmap := []termbox.InputMode{
					termbox.InputEsc | termbox.InputMouse,
					termbox.InputAlt | termbox.InputMouse,
					termbox.InputEsc,
					termbox.InputAlt,
				}
				inputmode++
				if inputmode >= len(chmap) {
					inputmode = 0
				}
				termbox.SetInputMode(chmap[inputmode])
			}
			if ev.Key == termbox.KeyCtrlX {
				ctrlxpressed = true
			} else {
				ctrlxpressed = false
			}

			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			//draw_keyboard()
			e.msg("Key: %v", ev.Ch)
			e.Display()
			//dispatch_press(&ev)
			//pretty_print_press(&ev)
			termbox.Flush()
		case termbox.EventResize:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			e.Cols, e.Lines = termbox.Size()
			e.msg("Resize: h %d,w %d", e.Lines, e.Cols)
			//draw_keyboard()
			e.Display()
			//pretty_print_resize(&ev)
			termbox.Flush()
		case termbox.EventMouse:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			//log.Println("[", e.Lines, e.Cols, ev, "event", ctrlxpressed)
			e.msg("Mouse: %d,%d [%d %d]", ev.MouseX, ev.MouseY, e.Cols, e.Lines)
			//draw_keyboard()
			e.Display()
			//pretty_print_mouse(&ev)
			termbox.Flush()
		case termbox.EventError:
			panic(ev.Err)
		}
	}

	return
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
	e.Cols, e.Lines = termbox.Size()
	if e.Msgflag {
		e.drawstring(0, e.Lines-1, termbox.ColorWhite, termbox.ColorBlack, e.Msgline)
	}
	// clear to end of line?
}

func (e *Editor) Display() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	e.DisplayMsg()
	//termbox.Sync()
}
