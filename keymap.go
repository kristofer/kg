package kg

type actionFunc func(*Editor)

type keymapt struct {
	KeyDesc  string
	KeyBytes string
	Do       actionFunc // function to call for Keymap-ping
}

/* desc, keys, func */
//keymap_t keymap[] = {
var keymap = []keymapt{
	{"C-a beginning-of-line    ", "\x01", (*Editor).lnbegin},
	{"C-b backward-char        ", "\x02", (*Editor).left},
	{"ArrowLeft backward-char  ", "\uFFEB", (*Editor).left},
	{"C-d delete               ", "\x04", (*Editor).delete},
	{"C-e end-of-line          ", "\x05", (*Editor).lnend},
	{"C-f forward-char         ", "\u0006", (*Editor).right},
	{"ArrowRight forward-char  ", "\uFFEA", (*Editor).right},
	{"C-h backspace            ", "\x08", (*Editor).backsp},
	{"C-g quit-quit            ", "\x07", (*Editor).quitquit},
	{"C-k kill-to-eol          ", "\x0B", (*Editor).killtoeol},
	{"C-l refresh              ", "\x0C", (*Editor).redraw},
	{"C-n next-line            ", "\x0E", (*Editor).down},
	{"ArrowDown next-line      ", "\uFFEC", (*Editor).down},
	{"C-p previous-line        ", "\x10", (*Editor).up},
	{"ArrowUp next-line        ", "\uFFED", (*Editor).up},
	{"C-s search               ", "\x13", (*Editor).search},
	{"C-r search               ", "\x12", (*Editor).rsearch},
	{"C-v forward-page         ", "\x16", (*Editor).pgdown},
	{"C-w kill-region          ", "\x17", (*Editor).cut},
	{"C-y yank                 ", "\x19", (*Editor).paste},
	{"C-space set-mark         ", "\x00", (*Editor).iblock},
	{"C-x 1 delete-other-window", "\x18\x31", (*Editor).deleteOtherWindows},
	{"C-x 2 split-window       ", "\x18\x32", (*Editor).splitWindow},
	{"C-x o other-window       ", "\x18\x6F", (*Editor).nextWindow},
	{"C-x = cursor-position    ", "\x18\x3D", (*Editor).showpos},
	{"C-x i insert-file        ", "\x18\x69", (*Editor).insertfile},
	{"C-x k kill-buffer        ", "\x18\x6B", (*Editor).killBuffer},
	{"C-x C-n next-buffer      ", "\x18\x0E", (*Editor).nextBuffer},
	{"C-x n next-buffer        ", "\x18\x6E", (*Editor).nextBuffer},
	{"C-x C-f find-file        ", "\x18\x06", (*Editor).readfile},
	{"C-x C-s save-buffer      ", "\x18\x13", (*Editor).savebuffer},
	{"C-x C-w write-file       ", "\x18\x17", (*Editor).writefile}, /* write and prompt for name */
	{"C-x C-c exit             ", "\x18\x03", (*Editor).quitAsk},
	{"esc b back-word          ", "\x1B\x62", (*Editor).wleft},
	{"esc f forward-word       ", "\x1B\x66", (*Editor).wright},
	{"esc g gotoline           ", "\x1B\x67", (*Editor).gotoline},
	{"esc k kill-region        ", "\x1B\x6B", (*Editor).cut},
	{"esc r query-replace      ", "\x1B\x72", (*Editor).queryReplace},
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
	{"ins toggle-overwrite-mode", "\x1B\x5B\x32\x7E", (*Editor).toggleOverwriteMode}, /* Ins key */
	{"del forward-delete-char  ", "\x1B\x5B\x33\x7E", (*Editor).delete},              /* Del key */
	{"backspace delete-left    ", "\x7f", (*Editor).backsp},
	{"up previous-line         ", "\x1B\x5B\x41", (*Editor).up},
	{"down next-line           ", "\x1B\x5B\x42", (*Editor).down},
	{"left backward-character  ", "\x1B\x5B\x44", (*Editor).left},
	{"right forward-character  ", "\x1B\x5B\x43", (*Editor).right},
	{"home beginning-of-line   ", "\x1B\x5B\x48", (*Editor).lnbegin},
	{"end end-of-line          ", "\x1B\x5B\x46", (*Editor).lnend},
	{"pgup backward-page       ", "\x1B\x5B\x35\x7E", (*Editor).pgup},   /* PgUp key */
	{"pgdn forward-page        ", "\x1B\x5B\x36\x7E", (*Editor).pgdown}, /* PgDn key */
	{"resize resize-terminal   ", "\x9A", (*Editor).resizeTerminal},
	{"K_ERROR                  ", "", nil},
}

// func get_key() {

// }

// // getinput get input from the Msgline.
// func getinput(prompt string, buf string, nbuf int, flag bool) rune {

// 	return '|'
// }
