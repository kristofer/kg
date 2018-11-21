package tkg

import (
	"log"
	"unicode"

	termbox "github.com/nsf/termbox-go"
)

//type editFunc ((*Editor)func())

func (e *Editor) quit() { e.Done = true }
func (e *Editor) quitquit() {
	e.EscapeFlag = false
	e.CtrlXFlag = false
	e.msg("Quit.\x07")
}
func (e *Editor) up() {
	e.CurrentBuffer.PointUp()
	//e.PointUp()
}
func (e *Editor) down() {
	e.CurrentBuffer.PointDown()

}
func (e *Editor) lnbegin() {
	e.CurrentBuffer.SetPoint(e.CurrentBuffer.LineForPoint(e.CurrentBuffer.Point()))
}
func (e *Editor) lnend() {
	e.CurrentBuffer.SetPoint(e.CurrentBuffer.LineEnd(e.CurrentBuffer.Point()))
}
func (e *Editor) version() { e.msg(VERSION) }
func (e *Editor) top() {
	e.CurrentBuffer.SetPoint(0)
}
func (e *Editor) bottom() {
	e.CurrentBuffer.SetPoint(e.CurrentBuffer.BufferLen())
	if e.CurrentBuffer.PageEnd < e.CurrentBuffer.BufferLen() {
		e.CurrentBuffer.Reframe = true
	}
}
func (e *Editor) block() {
	e.CurrentBuffer.Mark = e.CurrentBuffer.Point()
}
func (e *Editor) copy() {
	//copy_cut(false)
}
func (e *Editor) cut() {
	//copy_cut(true)
}
func (e *Editor) resize_terminal() {
	e.CurrentWindow.OneWindow()
}

func (e *Editor) quit_ask() {
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
	//var ch rune

	e.DisplayPromptAndResponse(prompt, "")
	e.MiniBufActive = true
	defer func() { e.MiniBufActive = false }()
	ev := <-e.EventChan
	log.Println("Mini ev", ev)
	// := e.HandleEvent(&ev)
	// ch = getch();
	ch := ev.Ch
	if ch == '\r' || ch == '\n' {
		return flag

	}
	return unicode.ToLower(ch) == 'y'
}

func (e *Editor) redraw() {
	log.Println("editor redraw")
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	e.CurrentWindow.Updated = true
	e.UpdateDisplay()
	termbox.Flush()
}

func (e *Editor) left() {
	e.CurrentBuffer.PointPrevious()
}

func (e *Editor) right() {
	e.CurrentBuffer.PointNext()
}

func (e *Editor) wleft() {

}

func (e *Editor) pgdown() {

}

func (e *Editor) pgup() {

}

func (e *Editor) wright() {

}

func (e *Editor) insert() {

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
	// int line;
	// point_t p;

	// if (getinput("Goto line: ", temp, STRBUF_S, F_CLEAR)) {
	// 	line = atoi(temp);
	// 	p = line_to_point(line);
	// 	if (p != -1) {
	// 		curbp->b_point = p;
	// 		msg("Line %d", line);
	// 	} else {
	// 		msg("Line %d, not found", line);
	// 	}
	// }
}

func (e *Editor) insertfile() {
	fname := e.GetFilename("Insert file: ")
	if fname != "" {
		res := e.InsertFile(fname, false)
		if res {
			//e.msg("Loaded file %s", fname)
		}
	}
}

func (e *Editor) readfile() {
	// buffer_t *bp;

	// temp[0] = '\0';
	// int result = getfilename("Find file: ", (char*)temp, NAME_MAX);
	// /* int result = getinput("Find file: ", (char*)temp, NAME_MAX, F_CLEAR); */

	// if (result) {
	// 	bp = find_buffer(temp, TRUE);
	// 	disassociate_b(curwp); /* we are leaving the old buffer for a new one */
	// 	curbp = bp;
	// 	associate_b2w(curbp, curwp);

	// 	/* load the file if not already loaded */
	// 	if (bp != NULL && bp->b_fname[0] == '\0') {
	// 		if (!load_file(temp)) {
	// 			msg("New file %s", temp);
	// 		}
	// 		strncpy(curbp->b_fname, temp, NAME_MAX);
	// 		curbp->b_fname[NAME_MAX] = '\0'; /* truncate if required */
	// 	}
	// }
}

func (e *Editor) savebuffer() {
	if e.CurrentBuffer.Filename != "" {
		e.Save(e.CurrentBuffer.Filename)
		return
	}
	// } else {
	// 	writefile();
	// }
	e.Refresh()
}

func (e *Editor) writefile() {
	// if (getinput("Write file: ", temp, NAME_MAX, F_NONE))
	fname := e.GetFilename("Write file: ")
	if e.Save(fname) == true {
		e.CurrentBuffer.Filename = fname
	}
}

func (e *Editor) killbuffer() {
	// buffer_t *kill_bp = curbp;
	// buffer_t *bp;
	// int bcount = count_buffers();

	// /* do nothing if only buffer left is the scratch buffer */
	// if (bcount == 1 && 0 == strcmp(get_buffer_name(curbp), "*scratch*"))
	// 	return;

	// if (curbp->b_flags & B_MODIFIED) {
	// 	mvaddstr(MSGLINE, 0, "Discard changes (y/n) ?");
	// 	clrtoeol();
	// 	if (!yesno(FALSE))
	// 		return;
	// }

	// if (bcount == 1) {
	// 	/* create a scratch buffer */
	// 	bp = find_buffer("*scratch*", TRUE);
	// 	strcpy(bp->b_bname, "*scratch*");
	// }

	// next_buffer();
	// assert(kill_bp != curbp);
	// delete_buffer(kill_bp);
}

func (e *Editor) iblock() {
	e.block()
	e.msg("Mark set")
}

func (e *Editor) toggle_overwrite_mode() {
	// NEVER!!
}

func (e *Editor) killtoeol() {
	//     if (curbp->b_point == pos(curbp, curbp->b_ebuf))
	// 	return; /* do nothing if at end of file */
	// if (*(ptr(curbp, curbp->b_point)) == 0xa) {
	// 	delete(); /* delete CR if at start of empty line */
	// } else {
	// 	curbp->b_mark = curbp->b_point;
	// 	lnend();
	// 	if (curbp->b_mark != curbp->b_point) copy_cut(TRUE);
	// }
}

func (e *Editor) copy_cut(cut int) {
	// char_t *p;
	// /* if no mark or point == marker, nothing doing */
	// if (curbp->b_mark == NOMARK || curbp->b_point == curbp->b_mark)
	// 	return;
	// if (scrap != NULL) {
	// 	free(scrap);
	// 	scrap = NULL;
	// }
	// if (curbp->b_point < curbp->b_mark) {
	// 	/* point above marker: move gap under point, region = marker - point */
	// 	(void) movegap(curbp, curbp->b_point);
	// 	p = ptr(curbp, curbp->b_point);
	// 	nscrap = curbp->b_mark - curbp->b_point;
	// } else {
	// 	/* if point below marker: move gap under marker, region = point - marker */
	// 	(void) movegap(curbp, curbp->b_mark);
	// 	p = ptr(curbp, curbp->b_mark);
	// 	nscrap = curbp->b_point - curbp->b_mark;
	// }
	// if ((scrap = (char_t*) malloc(nscrap)) == NULL) {
	// 	msg("No more memory available.");
	// } else {
	// 	(void) memcpy(scrap, p, nscrap * sizeof (char_t));
	// 	if (cut) {
	// 		curbp->b_egap += nscrap; /* if cut expand gap down */
	// 		curbp->b_point = pos(curbp, curbp->b_egap); /* set point to after region */
	// 		curbp->b_flags |= B_MODIFIED;
	// 		msg("%ld bytes cut.", nscrap);
	// 	} else {
	// 		msg("%ld bytes copied.", nscrap);
	// 	}
	// 	curbp->b_mark = NOMARK;  /* unmark */
	// }
}

func (e *Editor) paste() {
	// if(curbp->b_flags & B_OVERWRITE)
	// 	return;
	// if (nscrap <= 0) {
	// 	msg("Scrap is empty.  Nothing to paste.");
	// } else if (nscrap < curbp->b_egap - curbp->b_gap || growgap(curbp, nscrap)) {
	// 	curbp->b_point = movegap(curbp, curbp->b_point);
	// 	memcpy(curbp->b_gap, scrap, nscrap * sizeof (char_t));
	// 	curbp->b_gap += nscrap;
	// 	curbp->b_point = pos(curbp, curbp->b_egap);
	// 	curbp->b_flags |= B_MODIFIED;
	// }
}

func (e *Editor) showpos() {
	x, y := e.CurrentBuffer.XYForPoint(e.CurrentBuffer.Point())
	cl, ll := e.CurrentBuffer.GetLineStats()
	e.msg("(%d,%d) CurrLine %d LastLine %d", x, y, cl, ll)
}
