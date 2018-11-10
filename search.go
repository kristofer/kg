package kg

/* search.c, Atto Emacs, Public Domain, Hugh Barney, 2016, Derived from: Anthony's Editor January 93 */

//#include "header.h"

// #define FWD_SEARCH 1
// #define REV_SEARCH 2
// const (
// 	FWD_SEARCH = 1
// 	REV_SEARCH = 2
// )

// func search() {
// 	cpos := 0
// 	c := 0
// 	var (
// 		o_point Point = Curbp.b_point
// 	found Point
// 	)

// 	searchtext = ""
// 	display_prompt_and_response("Search: ", searchtext);
// 	cpos = strlen(searchtext);

// 	for {
// 		c = getch();
// 		/* ignore control keys other than C-g, backspace, CR,  C-s, C-R, ESC */
// 		if (c < 32 && c != 07 && c != 0x08 && c != 0x13 && c != 0x12 && c != 0x1b) {}
// 			continue
// 		}
// 		switch(c) {
// 		case 0x1b: /* esc */
// 			searchtext[cpos] = '\0';
// 			flushinp(); /* discard any escape sequence without writing in buffer */
// 			return;

// 		case 0x07: /* ctrl-g */
// 			curbp->b_point = o_point;
// 			return;

// 		case 0x13: /* ctrl-s, do the search */
// 			found = search_forward(curbp, curbp->b_point, searchtext);
// 			display_search_result(found, FWD_SEARCH, "Search: ", searchtext);
// 			break;

// 		case 0x12: /* ctrl-r, do the search */
// 			found = search_backwards(curbp, curbp->b_point, searchtext);
// 			display_search_result(found, REV_SEARCH, "Search: ", searchtext);
// 			break;

// 		case 0x7f: /* del, erase */
// 		case 0x08: /* backspace */
// 			if (cpos == 0)
// 				continue;
// 			searchtext[--cpos] = '\0';
// 			display_prompt_and_response("Search: ", searchtext);
// 			break;

// 		default:
// 			if (cpos < STRBUF_M - 1) {
// 				searchtext[cpos++] = c;
// 				searchtext[cpos] = '\0';
// 				display_prompt_and_response("Search: ", searchtext);
// 			}
// 			break;
// 		}
// 	}
// }

// func display_search_result(found Point, dir int, prompt string, search string)
// {
// 	if (found != -1 ) {
// 		Curbp.b_point = found
// 		msg("%s%s",prompt, search)
// 		display(Curwp, true)
// 	} else {
// 		msg("Failing %s%s",prompt, search)
// 		dispmsg()
// 		Curbp.b_point = (dir == FWD_SEARCH ? 0 : pos(Curbp, Curbp.b_ebuf))
// 	}
// }

// func search_forward(bp *Buffer, start_p Point, stext string) Point
// {
// 	point_t end_p = pos(bp, bp->b_ebuf);
// 	point_t p,pp;
// 	char* s;

// 	if (0 == strlen(stext))
// 		return start_p;

// 	for (p=start_p; p < end_p; p++) {
// 		for (s=stext, pp=p; *s == *(ptr(bp, pp)) && *s !='\0' && pp < end_p; s++, pp++)
// 			;

// 		if (*s == '\0')
// 			return pp;
// 	}

// 	return -1;
// }

// func search_backwards(bp *Buffer, start_p Point, stext string) Point
// {
// 	var p, pp Point
// 	p = 0
// 	pp = 0
// 	s := ""

// 	if len(stext) == 0 {
// 		return start_p;
// 	}
// 	for p = start_p; p >= 0; p-- {
// 		for s=stext, pp=p; *s == *(ptr(bp, pp)) && *s != '\0' && pp > -1; s++, pp++ {}

// 		if (*s == '\0') {
// 			if (p > 0)
// 				p--;
// 			return p;
// 		}
// 	}
// 	return -1;
// }