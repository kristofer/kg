package tkg

import "log"

const (
	fwdsearch = 1
	revsearch = 2
)

func (e *Editor) search() {
	searchtext := e.getInput("Search: ")
	found := e.CurrentBuffer.searchForward(e.CurrentBuffer.Point(), searchtext)
	e.displaySearchResult(found, fwdsearch, "Search: ", searchtext)
}
func (e *Editor) rsearch() {
	searchtext := e.getInput("R-Search: ")
	found := e.CurrentBuffer.searchBackwards(e.CurrentBuffer.Point(), searchtext)
	e.displaySearchResult(found, revsearch, "R-Search: ", searchtext)
}

func (e *Editor) displaySearchResult(found int, dir int, prompt string, search string) {
	if found != -1 {
		e.CurrentBuffer.SetPoint(found)
		e.msg("%s%s", prompt, search)
		e.Display(e.CurrentWindow, true)
	} else {
		e.msg("Failing %s%s", prompt, search)
		e.displayMsg()
		// if dir == fwdsearch {
		// 	e.bottom()
		// } else {
		// 	e.top()
		// }
	}
}

func (bp *Buffer) searchForward(startp int, stext string) int {
	endpt := bp.BufferLen() - 1
	if len(stext) == 0 {
		return -1
	}
	for p := startp; p < endpt; p++ {
		s := []rune(stext)
		ss := 0
		pp := 0
		for pp = p; pp < endpt; pp++ {
			rch, _ := bp.RuneAt(pp)
			if ss < len(s) && s[ss] == rch {
				ss++
			} else {
				break
			}
		}
		if ss == len(s) {
			log.Printf("Found %s at pt %d\n", stext, pp)
			return pp
		}
	}
	return -1
}

func (bp *Buffer) searchBackwards(startp int, stext string) int {
	endpt := bp.BufferLen() - 1
	if len(stext) == 0 {
		return startp
	}
	for p := startp; p >= 0; p-- {
		s := []rune(stext)
		ss := 0
		pp := 0
		for pp = p; pp < endpt; pp++ {
			rch, _ := bp.RuneAt(pp)
			if ss < len(s) && s[ss] == rch {
				ss++
			} else {
				break
			}
		}
		if ss == len(s) {
			log.Printf("Found %s at pt %d\n", stext, pp)
			return pp
		}
	}
	return -1
}
