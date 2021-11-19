package kg

import (
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

func (e *Editor) quit() { e.Done = true }
func (e *Editor) quitquit() {
	e.EscapeFlag = false
	e.CtrlXFlag = false
	e.msg("Quit.\x07")
}
func (e *Editor) up() {
	e.CurrentBuffer.PointUp()
}
func (e *Editor) down() {
	e.CurrentBuffer.PointDown()
}
func (e *Editor) lnbegin() {
	e.CurrentBuffer.SetPoint(e.CurrentBuffer.LineStart(e.CurrentBuffer.Point))
}
func (e *Editor) lnend() {
	e.CurrentBuffer.SetPoint(e.CurrentBuffer.LineEnd(e.CurrentBuffer.Point))
}
func (e *Editor) version() { e.msg(version) }
func (e *Editor) top() {
	e.CurrentBuffer.SetPoint(0)
}
func (e *Editor) bottom() {
	e.CurrentBuffer.SetPoint(e.CurrentBuffer.TextSize - 1)
	e.CurrentBuffer.Reframe = true
	e.CurrentBuffer.PageEnd = e.CurrentBuffer.TextSize - 1
}
func (e *Editor) block() {
	e.CurrentBuffer.Mark = e.CurrentBuffer.Point
}
func (e *Editor) copy() {
	e.copyCut(false)
}
func (e *Editor) cut() {
	e.copyCut(true)
}
func (e *Editor) resizeTerminal() {
	e.CurrentWindow.OneWindow()
}

func (e *Editor) quitAsk() {
	if e.ModifiedBuffers() == true {
		prompt := "Modified buffers exist; really exit (y/n) ?"
		if !e.yesno(false, prompt) {
			return
		}
	}
	e.quit()
}

/* flag = default answer, FALSE=n, TRUE=y */
func (e *Editor) yesno(flag bool, prompt string) bool {
	e.displayPromptAndResponse(prompt, "")
	e.MiniBufActive = true
	defer func() { e.MiniBufActive = false }()
	ev := <-e.EventChan
	ch := ev.Ch
	if ch == '\r' || ch == '\n' {
		return flag

	}
	return unicode.ToLower(ch) == 'y'
}

func (e *Editor) redraw() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	e.CurrentWindow.Updated = true
	e.CurrentBuffer.Reframe = true
	k := 0
	for wp := e.RootWindow; wp != nil; wp = wp.Next {
		wp.Updated = true
		k++
	}
	e.msg("editor redraw win(%d)", k)
	e.updateDisplay()
}

func (e *Editor) left() {
	e.CurrentBuffer.PointPrevious()
}

func (e *Editor) right() {
	e.CurrentBuffer.PointNext()
}

func (e *Editor) wleft() {

}
func (e *Editor) wright() {

}

func (e *Editor) pgdown() {
	pt := e.CurrentBuffer.Point
	l1 := e.CurrentBuffer.LineForPoint(e.CurrentBuffer.PageEnd)
	l2 := l1 + e.CurrentWindow.Rows - 2
	npt := e.CurrentBuffer.PointForLine(l2)
	log.Printf("start %d last line %d next %d new pt %d\n", pt, l1, l2, npt)
	e.CurrentBuffer.SetPoint(npt)
	e.CurrentWindow.Updated = true
}

func (e *Editor) pgup() {
	l1 := e.CurrentBuffer.LineForPoint(e.CurrentBuffer.PageStart)
	l2 := l1 - e.CurrentWindow.Rows - 2
	npt := e.CurrentBuffer.PointForLine(l2)
	log.Printf("last line %d next %d new pt %d\n", l1, l2, npt)
	e.CurrentBuffer.SetPoint(npt)
}

func (e *Editor) backsp() {
	e.CurrentBuffer.Backspace()
	e.CurrentBuffer.MarkModified()
}

func (e *Editor) delete() {
	e.CurrentBuffer.Delete()
	e.CurrentBuffer.MarkModified()
}

func (e *Editor) gotoline() {
	fname := e.getInput("Goto Line: ")
	ln, err := strconv.Atoi(fname)
	if err != nil {
		e.msg("Invalid Line.")
	}
	e.CurrentBuffer.gotoLine(ln)
}

func (e *Editor) insertfile() {
	fname := e.getInput("Insert file: ")
	if fname != "" {
		res := e.InsertFile(fname, true)
		if res {
			e.msg("Loaded file %s", fname)
		}
	}
}

func (e *Editor) readfile() {
	fname := e.getInput("Find file: ")
	if fname == "" {
		e.msg("Nope")
		return
	}
	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		e.msg("Failed to find file \"%s\".", fname)
		return
	}
	bp := e.FindBuffer(fname, true)
	e.CurrentWindow.DisassociateBuffer()
	e.CurrentBuffer = bp
	e.CurrentWindow.AssociateBuffer(bp)
	e.CurrentBuffer.Filename = fname
	bp.setText(string(dat))
}

func (e *Editor) savebuffer() {
	if e.CurrentBuffer.Filename != "" {
		e.Save(e.CurrentBuffer.Filename)
	} else {
		e.writefile()
	}
	e.Refresh()
}

func (e *Editor) writefile() {
	fname := e.getInput("Write file: ")
	if e.Save(fname) == true {
		e.CurrentBuffer.Filename = fname
	}
}

func (e *Editor) killBuffer() {
	killbp := e.CurrentBuffer
	bcount := e.CountBuffers()
	// do nothing if only buffer left is the scratch buffer
	if bcount == 1 && strings.Compare(e.GetBufferName(e.CurrentBuffer), "*scratch*") == 0 {
		return
	}
	if e.CurrentBuffer.modified == true {
		q := "Discard changes (y/n) ?"
		if !e.yesno(false, q) {
			return
		}
	}
	if bcount == 1 {
		bp := e.FindBuffer("*scratch*", true)
		bp.Filename = "*scratch*"
	}
	e.nextBuffer()
	if killbp != e.CurrentBuffer {
		e.deleteBuffer(killbp)
	}
}

func (e *Editor) iblock() {
	e.block()
	e.msg("Mark set")
}

func (e *Editor) toggleOverwriteMode() {
	e.msg("NEVER!! no overwite mode, you philistine.")
}

func (e *Editor) killtoeol() {
	bp := e.CurrentBuffer
	pt := e.CurrentBuffer.Point
	for i := 0; i < bp.LineLenAtPoint(pt)-bp.ColumnForPoint(pt); i++ {
		bp.Delete()
	}
}

func (e *Editor) copyCut(cut bool) {
	bp := e.CurrentBuffer
	pt := bp.Point
	if bp.Mark == nomark || pt == bp.Mark {
		return
	}
	extent := 0
	start := 0
	if pt < bp.Mark {
		extent = bp.Mark - pt
		start = pt
	} else { // bp.Point > bp.Mark
		extent = pt - bp.Mark
		start = bp.Mark
	}
	scrap := make([]rune, extent)
	l := start
	for k := 0; k < extent; k++ {
		rch, err := bp.RuneAt(l)
		if err != nil {
			e.msg("Copy/Cut failed. %s", err)
		}
		scrap[k] = rch
		l++
	}
	e.PasteBuffer = string(scrap)
	if cut == true {
		bp.Remove(start, extent)
		e.msg("%d characters cut.", extent)
	} else {
		e.msg("%d bytes copied.", extent)
	}
	bp.Mark = nomark
}

func (e *Editor) paste() {
	if len(e.PasteBuffer) <= 0 {
		e.msg("PasteBuffer is empty.  Nothing to paste.")
	} else {
		e.CurrentBuffer.Insert(e.PasteBuffer)
	}
}

func (e *Editor) showpos() {
	x, y := e.CurrentBuffer.XYForPoint(e.CurrentBuffer.Point)
	cl, ll := e.CurrentBuffer.GetLineStats()
	e.msg("(%d,%d) CurrLine %d LastLine %d", x, y, cl, ll)
}
