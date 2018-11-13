package tkg

/* desc, keys, func */
//keymap_t keymap[] = {
var keymap = []Keymapt{
	{"C-a beginning-of-line    ", "\x01", (*Editor).lnbegin},
	{"C-b backward-char        ", "\x02", (*Editor).left},
	{"C-d delete               ", "\x04", (*Editor).delete},
	{"C-e end-of-line          ", "\x05", (*Editor).lnend},
	{"C-f foward-char          ", "\x06", (*Editor).right},
	{"C-h backspace            ", "\x08", (*Editor).backsp},
	{"C-k kill-to-eol          ", "\x0B", (*Editor).killtoeol},
	{"C-l refresh              ", "\x0C", (*Editor).redraw},
	{"C-n next-line            ", "\x0E", (*Editor).down},
	{"C-p previous-line        ", "\x10", (*Editor).up},
	// {"C-s search               ", "\x13", (*Editor).search},
	// {"C-r search               ", "\x12", (*Editor).search},
	{"C-v forward-page         ", "\x16", (*Editor).pgdown},
	{"C-w kill-region          ", "\x17", (*Editor).cut},
	{"C-y yank                 ", "\x19", (*Editor).paste},
	{"C-space set-mark         ", "\x00", (*Editor).iblock},
	{"C-x 1 delete-other-window", "\x18\x31", (*Editor).DeleteOtherWindows},
	{"C-x 2 split-window       ", "\x18\x32", (*Editor).SplitWindow},
	{"C-x o other-window       ", "\x18\x6F", (*Editor).NextWindow},
	{"C-x = cursor-position    ", "\x18\x3D", (*Editor).showpos},
	//	{"C-x i insert-file        ", "\x18\x69", (*Editor).InsertFile},
	{"C-x k kill-buffer        ", "\x18\x6B", (*Editor).killbuffer},
	{"C-x C-n next-buffer      ", "\x18\x0E", (*Editor).NextBuffer},
	{"C-x n next-buffer        ", "\x18\x6E", (*Editor).NextBuffer},
	{"C-x C-f find-file        ", "\x18\x06", (*Editor).readfile},
	{"C-x C-s save-buffer      ", "\x18\x13", (*Editor).savebuffer},
	{"C-x C-w write-file       ", "\x18\x17", (*Editor).writefile}, /* write and prompt for name */
	{"C-x C-c exit             ", "\x18\x03", (*Editor).quit_ask},
	{"esc b back-word          ", "\x1B\x62", (*Editor).wleft},
	{"esc f forward-word       ", "\x1B\x66", (*Editor).wright},
	{"esc g gotoline           ", "\x1B\x67", (*Editor).gotoline},
	{"esc k kill-region        ", "\x1B\x6B", (*Editor).cut},
	//	{"esc r query-replace      ", "\x1B\x72", (*Editor).query_replace},
	{"esc v backward-page      ", "\x1B\x76", (*Editor).pgup},
	{"esc w copy-region        ", "\x1B\x77", (*Editor).copy},
	{"esc @ set-mark           ", "\x1B\x40", (*Editor).iblock}, /* esc-@ */
	{"esc < beg-of-buf         ", "\x1B\x3C", (*Editor).top},
	{"esc > end-of-buf         ", "\x1B\x3E", (*Editor).bottom},
	{"esc home, beg-of-buf     ", "\x1B\x1B\x4F\x48", (*Editor).top},
	{"esc end, end-of-buf      ", "\x1B\x1B\x4F\x46", (*Editor).bottom},
	{"esc up, beg-of-buf       ", "\x1B\x1B\x5B\x41", (*Editor).top},
	{"esc down, end-of-buf     ", "\x1B\x1B\x5B\x42", (*Editor).bottom},
	{"esc esc show-version     ", "\x1B\x1B", (*Editor).version},
	{"ins toggle-overwrite-mode", "\x1B\x5B\x32\x7E", (*Editor).toggle_overwrite_mode}, /* Ins key */
	{"del forward-delete-char  ", "\x1B\x5B\x33\x7E", (*Editor).delete},                /* Del key */
	{"backspace delete-left    ", "\x7f", (*Editor).backsp},
	{"up previous-line         ", "\x1B\x5B\x41", (*Editor).up},
	{"down next-line           ", "\x1B\x5B\x42", (*Editor).down},
	{"left backward-character  ", "\x1B\x5B\x44", (*Editor).left},
	{"right forward-character  ", "\x1B\x5B\x43", (*Editor).right},
	{"home beginning-of-line   ", "\x1B\x4F\x48", (*Editor).lnbegin},
	{"end end-of-line          ", "\x1B\x4F\x46", (*Editor).lnend},
	{"home beginning-of-line   ", "\x1B\x5B\x48", (*Editor).lnbegin},
	{"end end-of-line          ", "\x1B\x5B\x46", (*Editor).lnend},
	{"pgup backward-page       ", "\x1B\x5B\x35\x7E", (*Editor).pgup},   /* PgUp key */
	{"pgdn forward-page        ", "\x1B\x5B\x36\x7E", (*Editor).pgdown}, /* PgDn key */
	{"resize resize-terminal   ", "\x9A", (*Editor).resize_terminal},
	{"K_ERROR                  ", "", nil},
}

// func get_key(keymap_t *keys, keymap_t **key_return) {
// type Keymapt struct {
//	KeyDesc  string
//	KeyBytes string
//	Do       func(*Editor) // function to call for Keymap-ping
// }

func get_key() {
	// keymap_t *k;
	// int submatch;
	// static char_t buffer[K_BUFFER_LENGTH];
	// static char_t *record = buffer;

	// *key_return = NULL;

	// /* if recorded bytes remain, return next recorded byte. */
	// if (*record != '\0') {
	// 	*key_return = NULL;
	// 	return record++;
	// }
	// /* reset record buffer. */
	// record = buffer;

	// do {
	// 	assert(K_BUFFER_LENGTH > record - buffer);
	// 	/* read and record one byte. */
	// 	*record++ = (unsigned)getch();
	// 	*record = '\0';

	// 	/* if recorded bytes match any multi-byte sequence... */
	// 	for (k = keys, submatch = 0; k->key_bytes != NULL; ++k) {
	// 		char_t *p, *q;

	// 		for (p = buffer, q = (char_t *)k->key_bytes; *p == *q; ++p, ++q) {
	// 		        /* an exact match */
	// 			if (*q == '\0' && *p == '\0') {
	//     				record = buffer;
	// 				*record = '\0';
	// 				*key_return = k;
	// 				return record; /* empty string */
	// 			}
	// 		}
	// 		/* record bytes match part of a command sequence */
	// 		if (*p == '\0' && *q != '\0') {
	// 			submatch = 1;
	// 		}
	// 	}
	// } while (submatch);
	// /* nothing matched, return recorded bytes. */
	// record = buffer;
	// return (record++);
}

// getinput get input from the Msgline.
func getinput(prompt string, buf string, nbuf int, flag bool) rune {
	// int cpos = 0;
	// int c;
	// int start_col = strlen(prompt);

	// mvaddstr(MSGLINE, 0, prompt);
	// clrtoeol();

	// if (flag == F_CLEAR) buf[0] = '\0';

	// /* if we have a default value print it and go to end of it */
	// if (buf[0] != '\0') {
	// 	addstr(buf);
	// 	cpos = strlen(buf);
	// }

	// for (;;) {
	// 	refresh();
	// 	c = getch();
	// 	/* ignore control keys other than backspace, cr, lf */
	// 	if (c < 32 && c != 0x07 && c != 0x08 && c != 0x0a && c != 0x0d)
	// 		continue;

	// 	switch(c) {
	// 	case 0x0a: /* cr, lf */
	// 	case 0x0d:
	// 		buf[cpos] = '\0';
	// 		return (cpos > 0 ? TRUE : FALSE);

	// 	case 0x07: /* ctrl-g */
	// 		return FALSE;

	// 	case 0x7f: /* del, erase */
	// 	case 0x08: /* backspace */
	// 		if (cpos == 0)
	// 			continue;

	// 		move(MSGLINE, start_col + cpos - 1);
	// 		addch(' ');
	// 		move(MSGLINE, start_col + cpos - 1);
	// 		buf[--cpos] = '\0';
	// 		break;

	// 	default:
	// 		if (cpos < nbuf -1) {
	// 			addch(c);
	// 			buf[cpos++] = c;
	// 			buf[cpos] ='\0';
	// 		}
	// 		break;
	// 	}
	// }
	return '|'
}
