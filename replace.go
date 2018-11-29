package kg

import "fmt"

/*search for a string and replace it with another string */
func (e *Editor) queryReplace() {
	e.Searchtext = e.getInput("Query replace: ")
	if len(e.Searchtext) < 1 {
		return
	}
	e.Replace = e.getInput("With: ")
	slen := len(e.Searchtext)
	bp := e.CurrentBuffer
	opoint := bp.Point
	lpoint := -1
	ask := true
	/* build query replace question string */
	question := fmt.Sprintf("Replace '%s' with '%s' ? ", e.Searchtext, e.Replace)
	/* scan through the file, from point */
	numsub := 0
outer:
	for {
		found := bp.searchForward(bp.Point, e.Searchtext)
		/* if not found set the point to the last point of replacement, or where we started */
		if found == -1 {
			if lpoint == -1 {
				bp.SetPoint(opoint)
			} else {
				bp.SetPoint(lpoint)
			}
			break
		}
		bp.SetPoint(found)
		/* search_forward places point at end of search, move to start of search */
		for k := 0; k < slen; k++ {
			bp.PointPrevious()
		}
		e.Display(e.CurrentWindow, true)

		if ask == true {
			answer := e.getInput(question)

		inner:
			for {
				e.Display(e.CurrentWindow, true)
				resp := []rune(answer)
				c := ' '
				if len(resp) > 0 {
					c = resp[0]
				}
				switch c {
				case 'y': /* yes, substitute */
					break inner
				case 'n': /* no, find next */
					bp.SetPoint(found) /* set to end of search string */
				case '!': /* yes/stop asking, do the lot */
					ask = false
					break inner
				//case 0x1B: /* esc */
				//flushinp() /* discard any escape sequence without writing in buffer */
				case 'q': /* controlled exit */
					break outer
				default: /* help me */
					answer = e.getInput("(y)es, (n)o, (!)do the rest, (q)uit: ")
					//continue inner
				}
			}
		}
		for k := 0; k < slen; k++ { // delete found search text
			bp.Delete()
		}
		bp.Insert(e.Replace) // qed
		lpoint = bp.Point
		numsub++
	}
	e.msg("%d substitutions", numsub)
}
