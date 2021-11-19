package kg

import (
	"io/ioutil"
	"unicode"
)

// Refresh editor display(!)
func (e *Editor) Refresh() {}

// PosixFile foo
func (e *Editor) PosixFile(fname string) bool {
	fn := []rune(fname)
	if fn[0] == '_' {
		return false
	}
	for f := range fn {
		if (unicode.IsLetter(rune(f))) && f != '.' && f != '_' && f != '-' && f != '/' {
			return false
		}
	}
	return true
}

// Save foo
func (e *Editor) Save(fname string) bool {
	if e.PosixFile(fname) != true {
		e.msg("Not a portable POSIX file name.")
		return false
	}
	d1 := []byte(e.CurrentBuffer.getText())
	rch := d1[len(d1)-1]
	if rch != '\n' {
		prompt := "Last character is not newline. Add one?"
		if e.yesno(true, prompt) {
			d1 = append(d1, '\n')
		}
	}
	err := ioutil.WriteFile(fname, d1, 0644)
	if err != nil {
		e.msg("Failed to save file \"%s\".", fname)
		return false
	}
	e.CurrentBuffer.modified = false
	e.msg("File \"%s\" %d bytes saved.", fname, len(d1))
	return true
}

// LoadFile foo
// func (e *Editor) LoadFile(fname string) bool {
// 	return false
// }

// InsertFile reads file into buffer at point
func (e *Editor) InsertFile(fname string, modflag bool) bool {
	bp := e.CurrentBuffer
	dat, err := ioutil.ReadFile(fname)
	if err != nil {
		e.msg("Failed to read and insert file \"%s\".", fname)
	}
	if !modflag { // just do a load into buffer with no modification
		bp.setText(string(dat))
		bp.modified = false
	} else { // insert into buffer and mark as modified.
		bp.Insert(string(dat))
		bp.modified = true
	}
	e.msg("File \"%s\" %d bytes read.", fname, len(dat))
	return true
}
