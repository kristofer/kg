# kg Emacs

A small functional Emacs written in Go changed in time for the vault  

by Kristofer

kg Emacs is inspired by Atto Emacs (by Hugh Barney), MicroEmacs, Nano, Pico and his earlier project known as Perfect Emacs [1]. Also by JOE, uEmacs, MicroEmacs, mg, zile and of course, GNU Emacs. I learnt Emacs on a Sun 2/120 a long time ago in a galaxy far, far, away.


> A designer knows he has achieved perfection not when there is nothing left to add, but when there is nothing left to take away.
> -- <cite>Antoine de Saint-Exupery</cite>

> If you want to build a ship, don't drum up people to collect wood and don't assign them tasks and work, but rather teach them to long for the endless immensity of the sea.
> -- <cite>Antoine de Saint-Exupery</cite>

## Goals of kg Emacs

* Mine own, finally, in Go.
* uses an array of Rune to handle Unicode codepoints.
* no damn Overwrite mode. Too bad.
* pure Go implementation
* removal of the C-based hilite stuff
* add go routines for some operations.
* ...
* Be easy to understand without extensive study (to encourage further experimentation).

Using Atto as the lowest functional Emacs, Hugh had to consider the essential feature set that makes Emacs, 'Emacs'.  From his Readme: _I have defined this point as a basic Emacs command set and key bindings; the ability to edit multiple files (buffers), and switch between them; edit the buffers in mutliple windows, cut, copy and paste; forward and reverse searching, a replace function, basic syntax hilighting and UTF8 support. The proviso being that all this will fit in less than 2000 lines of C._


## Derivation

kg is based on the design of _Atto Emacs which is based on the public domain code of Anthony Howe's editor (commonly known as Anthony's Editor or AE, [2]).  Rather than representing a file as a linked list of lines, the AE Editor uses the concept of a Buffer-Gap [4,5,6].  A Buffer-Gap editor stores the file in a single piece of contiguous memory with some extra unused space known as the buffer gap.  On character insertion and deletion the gap is first moved to the current point.  A character deletion then extends the gap by moving the gap pointer back by 1 OR the gap is reduced by 1 when a character is inserted.  The Buffer-Gap technique is elegant and significantly reduces the amount of code required to load a file, modify it and redraw the display.  The proof of this is seen when you consider that Atto supports almost the same command set that Pico supports,  but Pico requires almost 17 times the amount of code._

## Comparisons with Other Emacs Implementations

    Editor         Binary   BinSize     KLOC  Files

    atto           atto       33002     1.9k     10
    pEmacs         pe         59465     5.7K     16
    Esatz-Emacs    ee         59050     5.7K     14
    GNOME          GNOME      55922     9.8k     13
    Zile           zile      257360    11.7k     48
    Mg             mg        585313    16.5K     50
    uEmacs/Pk      em        147546    17.5K     34
    Pico           pico      438534    24.0k     29
    Nano           nano      192008    24.8K     17
    jove           jove      248824    34.7k     94
    Qemacs         qe        379968    36.9k     59
    ue3.10         uemacs    171664    52.4K     16
    GNUEmacs       emacs   14632920   358.0k    186

    kg             kg       2719864     1.9k      8

* _Due to Go's single-binary, 2.7mb is pretty normal for a go binary size. It uses NO shared libraries, etc. The code itself though, is more like 720496, since that's (kg - HelloWorld) (2719864 - 1999368)._
* _While I did not aim for less than 2000 LOC, surprisingly, it is. And no, I did not count test LOC._ 

## kg Key Bindings

    C-A   begining-of-line
    C-B   backward-character
    C-D   delete-char
    C-E   end-of-line
    C-F   forward Character
    C-G	  Abort (at prompts)
    C-H   backspace
    C-I   handle-tab
    C-J   newline
    C-K   kill-to-eol
    C-L   refresh display
    C-M   Carrage Return
    C-N   next line
    C-P   previous line
    C-R   search-backwards
    C-S	  search-forwards
    C-U   Undo
    C-V   Page Down
    C-W   Kill Region (Cut)
    C-X   CTRL-X command prefix
    C-Y   Yank (Paste)

    M-<   Start of file
    M->   End of file
    M-v   Page Up
    M-f   Forward Word
    M-b   Backwards Word
    M-g   goto-line
    M-r   Search and Replace
    M-w   copy-region

    C-<spacebar> Set mark at current position.

    ^X^C  Exit. Any unsaved files will require confirmation.
    ^X^F  Find file; read into a new buffer created from filename.
    ^X^S  Save current buffer to disk, using the buffer's filename as the name of
    ^X^W  Write current buffer to disk. Type in a new filename at the prompt to
    ^Xi   Insert file at point
    ^X=   Show Character at position
    ^X^N  next-buffer
    ^Xn   next-buffer
    ^Xk   kill-buffer
    ^X1   delete-other-windows
    ^X2   split-window
    ^Xo   other-window

    Home  Beginning-of-line
    End   End-of-line
    Del   Delete character under cursor
    Ins   Toggle Overwrite Mode
    Left  Move left
    Right Move point right
    Up    Move to the previous line
    Down  Move to the next line
    Backspace delete caharacter on the left
    Ctrl+Up      beginning of file
    Ctrl+Down    end of file
    Ctrk+Left    Page Down
    Ctrl+Right   Page Up

### Copying and moving

    C-<spacebar> Set mark at current position
    ^W   Delete region
    ^Y   Yank back kill buffer at cursor
    M-w  Copy Region

A region is defined as the area between this mark and the current cursor position. The kill buffer is the text which has been most recently deleted or copied.

Generally, the procedure for copying or moving text is:
1. Mark out region using M-<spacebar> at the beginning and move the cursor to the end.
2. Delete it (with ^W) or copy it (with M-W) into the kill buffer.
3. Move the cursor to the desired location and yank it back (with ^Y).

### Searching

    C-S or C-R enters the search prompt, where you type the search string
    BACKSPACE - will reduce the search string, any other character will extend it
    C-S at the search prompt will search forward, will wrap at end of the buffer
    C-R at the search prompt will search backwards, will wrap at start of the buffer
    ESC will escape from the search prompt and return to the point of the match
    C-G abort the search and return to point before the search started

## Building on Linux/MacOS

When building on Linux/MacOS you will need to install Go v1.11 or greater.
You will also need the `github.com/nsf/termbox-go`

    $ go get github.com/kristofer/kg
    $ go get github.com/nsf/termbox-go

cd to the kg source directory...

    $ cd cmd
    $ go build -o kg main.go

and then move into your binary PATH

## Future Enhancements

Maybe a piece-table or piece-chain implementation? Maybe a different key mapping set to make it more mac-like.

## Multiple Windows or Not?

Kg supports multiple windows.

## Known Issues

* Goto-line will fail to go to the very last line.  This is a special case that could easily be fixed.

## Copying

  Kg code is released to the public domain.
  kryounger AT gmail.com - 2018

## Acknowledgements

    Hugh Barney for Atto (and his other work).
    Ed Davies for bringing Athony's Editor to my attention
    Anthony Howe for his original codebase
    Matt Fielding (Magnetic Realms) for providing fixes for multi-byte / wide characters, delete, backspace and cursor position
    The Infinnovation team for bug fixes to complete.c
    James Gosling for telling me long ago that to build an Emacs just something you did.

## References

    [1] Perfect Emacs - https://github.com/hughbarney/pEmacs
    [2] Anthony's Editor - https://github.com/hughbarney/Anthony-s-Editor
    [3] MG - https://github.com/rzalamena/mg
    [4] Jonathan Payne, Buffer-Gap: http://ned.rubyforge.org/doc/buffer-gap.txt
    [5] Anthony Howe,  http://ned.rubyforge.org/doc/editor-101.txt
    [6] Hugh Barney, https://github.com/hughbarney/atto

