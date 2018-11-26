package kg

import "fmt"

/*search for a string and replace it with another string */
func (e *Editor) queryReplace() {
	// point_t o_point = curbp->b_point;
	// point_t l_point = -1;
	// point_t found;
	// char question[STRBUF_L];
	// int slen, rlen;   /* length of search and replace strings */
	// int numsub = 0;   /* number of substitutions */
	// int ask = TRUE;
	// int c;

	// searchtext[0] = '\0';
	// replace[0] = '\0';

	e.Searchtext = e.getInput("Query replace: ")
	if len(e.Searchtext) < 1 {
		return
	}

	e.Replace = e.getInput("With: ")
	slen := len(e.Searchtext)
	rlen := len(e.Replace)
	bp := e.CurrentBuffer
	opoint := bp.Point()
	lpoint := -1
	ask := true
	/* build query replace question string */
	question := fmt.Sprintf("Replace '%s' with '%s' ? ", e.Searchtext, e.Replace)

	/* scan through the file, from point */
	numsub := 0
	for {
		found := bp.searchForward(bp.Point(), e.Searchtext)

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
		//curbp->b_point -= slen
		for k := 0; k < slen; k++ {
			bp.PointPrevious()
		}

		if ask == true {
			e.msg(question)

			for {
				e.Display(e.CurrentWindow, true)
				//c = getch();
				resp := e.getInput(question)
				c := resp[1]
				switch c {
				case 'y': /* yes, substitute */
					break

				case 'n': /* no, find next */
					bp.SetPoint(found) /* set to end of search string */
					continue

				case '!': /* yes/stop asking, do the lot */
					ask = false
					break

				case 0x1B: /* esc */
					//flushinp() /* discard any escape sequence without writing in buffer */
				case 'q': /* controlled exit */
					return

				default: /* help me */
					e.msg("(y)es, (n)o, (!)do the rest, (q)uit")
					continue
				}
			}
		}

		if rlen > slen {
			// 	movegap(curbp, found);
			// 	/*check enough space in gap left */
			// 	if (rlen - slen < curbp->b_egap - curbp->b_gap)
			// 		growgap(curbp, rlen - slen);
			// 	/* shrink gap right by r - s */
			// 	curbp->b_gap = curbp->b_gap + (rlen - slen);
			// } else if (slen > rlen) {
			// 	movegap(curbp, found);
			// 	/* stretch gap left by s - r, no need to worry about space */
			// 	curbp->b_gap = curbp->b_gap - (slen - rlen);
			// } else {
			// 	/* if rlen = slen, we just overwrite the chars, no need to move gap */
		}

		/* now just overwrite the chars at point in the buffer */
		// l_point = curbp->b_point;
		// memcpy(ptr(curbp, curbp->b_point), replace, rlen * sizeof (char_t));
		// curbp->b_flags |= B_MODIFIED;
		// curbp->b_point = found - (slen - rlen); /* end of replcement */
		numsub++
	}

	e.msg("%d substitutions", numsub)
}
