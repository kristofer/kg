package kg

import (
	"fmt"
)

/* main.c, Atto Emacs, Public Domain, Hugh Barney, 2016, Derived from: Anthony's Editor January 93 */

const (
	VERSION          = "kg 1.0, Public Domain, November 2018, Kristofer Younger,  No warranty."
	PROG_NAME        = "kg"
	B_MODIFIED       = 0x01 /* modified buffer */
	B_OVERWRITE      = 0x02 /* overwite mode */
	MSGLINE          = (LINES - 1)
	NOMARK           = -1
	CHUNK            = 8096
	K_BUFFER_LENGTH  = 256
	TEMPBUF          = 512
	STRBUF_L         = 256
	STRBUF_M         = 64
	STRBUF_S         = 16
	MIN_GAP_EXPAND   = 512
	TEMPFILE         = "/tmp/feXXXXXX"
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

//#include "header.h"

// int done;
// point_t nscrap;
// char_t *scrap;

// int msgflag;
// char_t *input;
// char msgline[TEMPBUF];
// char temp[TEMPBUF];
// char searchtext[STRBUF_M];
// char replace[STRBUF_M];

// keymap_t *key_return;
// keymap_t *key_map;
// buffer_t *curbp;			/* current buffer */
// buffer_t *bheadp;			/* head of list of buffers */
// window_t *curwp;
// window_t *wheadp;

var (
	// done int                /* Quit flag. */
	Done       bool   /* Quit flag. */
	Msgflag    bool   /* True if msgline should be displayed. */
	Nscrap     Point  /* Length of scrap buffer. */
	Scrap      string /* Allocated scrap buffer. */
	Input      ch     // RUNE?????
	Msgline    string /* Message line input/output buffer. */
	Temp       string /* Temporary buffer. */
	Searchtext string
	Replace    string
	Key_map    *Keymapt /* Command key mappings. */
	Keymap     []Keymapt
	Key_return *Keymapt /* Command key return */
	//
	Curbp  *Buffer /* current buffer */
	Bheadp *Buffer /* head of list of buffers */
	Curwp  *Window
	Wheadp *Window
)

// typedef unsigned char char_t;
//type Char character

// typedef long point_t
// typedef struct keymap_t {
// 	char *key_desc;                 /* name of bound function */
// 	char *key_bytes;		/* the string of bytes when this key is pressed */
// 	void (*func)(void);
// } keymap_t;
type Keymapt struct {
	KeyDesc  string
	KeyBytes string
	Do       *func() // function to call for Keymap-ping
}

// StartEditor is the old C main function
func StartEditor(argv []string, argc int) {

	// Still need for Termbox...

	// setlocale(LC_ALL, "") ; /* required for 3,4 byte UTF8 chars */
	if initscr() == nil {
		panic("%s: Failed to initialize the screen.\n")
	}
	//raw();
	//noecho();
	//idlok(stdscr, TRUE);

	// start_color();
	// init_pair(ID_DEFAULT, COLOR_CYAN, COLOR_BLACK);          /* alpha */
	// init_pair(ID_SYMBOL, COLOR_WHITE, COLOR_BLACK);          /* non alpha, non digit */
	// init_pair(ID_MODELINE, COLOR_BLACK, COLOR_WHITE);        /* modeline */
	// init_pair(ID_DIGITS, COLOR_YELLOW, COLOR_BLACK);         /* digits */
	// init_pair(ID_BLOCK_COMMENT, COLOR_GREEN, COLOR_BLACK);   /* block comments */
	// init_pair(ID_LINE_COMMENT, COLOR_GREEN, COLOR_BLACK);    /* line comments */
	// init_pair(ID_SINGLE_STRING, COLOR_YELLOW, COLOR_BLACK);  /* single quoted strings */
	// init_pair(ID_DOUBLE_STRING, COLOR_YELLOW, COLOR_BLACK);  /* double quoted strings */

	if argc > 1 {
		Curbp = find_buffer(argv[1], TRUE)
		insert_file(argv[1], FALSE)
		/* Save filename irregardless of load() success. */
		//strncpy(curbp->b_fname, argv[1], NAME_MAX);
		Curbp.b_fname = argv[1]
		//curbp->b_fname[NAME_MAX] = '\0'; /* force truncation */
	} else {
		Curbp = find_buffer("*scratch*", TRUE)
		//strncpy(curbp->b_bname, "*scratch*", STRBUF_S);
		Curbp.b_bname = "*scratch*"
	}
	Curwp = new_window()
	Wheadp = Curwp
	one_window(Curwp)
	associate_b2w(Curbp, Curwp)

	if !growgap(Curbp, CHUNK) {
		panic("%s: Failed to allocate required memory.\n")
	}
	Key_map = Keymap

	for Done != true {
		update_display()
		Input = get_key(Key_map, &key_return)

		if key_return != nil {
			key_return.Do()
		} else {
			/* allow TAB and NEWLINE, otherwise any Control Char is 'Not bound' */
			if unicode.isControl(Input) || Input == '\n' || Input == '\t' {
				insert()
			} else {
				flushinp() /* discard without writing in buffer */
				msg("Not bound")
			}
		}
	}

	//if (scrap != NULL) free(scrap);
	Scrap = ""
	move(LINES-1, 0)
	//refresh();
	//noraw();
	//endwin();
	return 0
}

func msg(args ...string) {
	Msgline = fmt.Sprintf("%#v", args)
	Msgflag = true
	return
}
