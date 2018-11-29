package kg

const (
	fwdsearch = 1
	revsearch = 2
)

func (e *Editor) search() {
	e.Searchtext = e.getInput("Search: ")
	found := e.CurrentBuffer.searchForward(e.CurrentBuffer.Point, e.Searchtext)
	e.displaySearchResult(found, fwdsearch, "Search: ", e.Searchtext)
}
func (e *Editor) rsearch() {
	e.Searchtext = e.getInput("R-Search: ")
	found := e.CurrentBuffer.searchBackwards(e.CurrentBuffer.Point, e.Searchtext)
	e.displaySearchResult(found, revsearch, "R-Search: ", e.Searchtext)
}

func (e *Editor) displaySearchResult(found int, dir int, prompt string, search string) {
	if found != -1 {
		e.CurrentBuffer.SetPoint(found)
		e.Display(e.CurrentWindow, true)
	} else {
		e.msg("Failing %s%s", prompt, search)
		e.displayMsg()
	}
}

func (bp *Buffer) searchForward(startp int, stext string) int {
	endpt := bp.TextSize - 1
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
			return pp
		}
	}
	return -1
}

func (bp *Buffer) searchBackwards(startp int, stext string) int {
	endpt := bp.TextSize - 1
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
			return pp
		}
	}
	return -1
}
