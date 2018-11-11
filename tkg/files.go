package tkg

// PosixFile foo
func (e *Editor) PosixFile(fname string) bool {
	// 	var fn []rune = []rune(fname)
	// 	if (fn[0] == '_') {
	// 		return false
	// }

	// 	for f := range fn {
	// 		if (!isalnum(f) && f != '.' && f != '_' && f != '-' && f != '/') {
	// 			return false
	// 		}
	// 	}
	return true
}

// Save foo
func (e *Editor) Save(fname string) bool {
	// 	//FILE *fp;
	// 	var length Point = 0

	// 	if (!PosixFile(fn)) {
	// 		msg("Not a portable POSIX file name.");
	// 		return false
	// 	}

	// 	fp = fopen(fn, "w");
	// 	if (fp == NULL) {
	// 		msg("Failed to open file \"%s\".", fn);
	// 		return false
	// 	}
	// 	(void) movegap(curbp, (point_t) 0);
	// 	length = (point_t) (curbp->b_ebuf - curbp->b_egap);
	// 	if (fwrite(curbp->b_egap, sizeof (char), (size_t) length, fp) != length) {
	// 		msg("Failed to write file \"%s\".", fn);
	// 		return false
	// 	}
	// 	if (fclose(fp) != 0) {
	// 		msg("Failed to close file \"%s\".", fn);
	// 		return false
	// 	}
	// 	curbp->b_flags &= ~B_MODIFIED;
	// 	msg("File \"%s\" %ld bytes saved.", fn, pos(curbp, curbp->b_ebuf));
	return true
}

// LoadFile foo
func (e *Editor) LoadFile(fname string) error {
	// 	/* reset the gap, make it the whole buffer */
	// 	curbp->b_gap = curbp->b_buf;
	// 	curbp->b_egap = curbp->b_ebuf;
	// 	top();
	return e.InsertFile(fname, false)
}

// InsertFile reads file into buffer at point
func (e *Editor) InsertFile(fname string, modflag bool) error {
	// {
	// 	// FILE *fp;
	// 	// size_t len;
	// 	// struct stat sb;

	// 	file, err := os.Open(fname)
	// if err != nil {
	//   fmt.Println(err)
	//   return false
	// }
	// defer file.Close()
	// 	// if (stat(fn, &sb) < 0) {
	// 	// 	msg("Failed to find file \"%s\".", fn);
	// 	// 	return false
	// 	// }
	// 	fileinfo, err := file.Stat()
	// 	if err != nil {
	// 	  fmt.Println(err)
	// 	  return false
	// 	}
	// 	// if (MAX_SIZE_T < sb.st_size) {
	// 	// 	msg("File \"%s\" is too big to load.", fn);
	// 	// 	return false
	// 	// }
	// 	filesize := fileinfo.Size()
	// 	buffer := make([]byte, filesize)

	// 	if (curbp->b_egap - curbp->b_gap < sb.st_size * sizeof (char_t) && !growgap(curbp, sb.st_size)) {
	// 		return false
	// }
	// 	// if ((fp = fopen(fn, "r")) == NULL) {
	// 	// 	msg("Failed to open file \"%s\".", fn);
	// 	// 	return false
	// 	// }
	// 	// curbp->b_point = movegap(curbp, curbp->b_point);
	// 	// curbp->b_gap += len = fread(curbp->b_gap, sizeof (char), (size_t) sb.st_size, fp);
	// 	bytesread, err := file.Read(buffer)
	// 	if err != nil {
	// 	  fmt.Println(err)
	//   return false
	// }

	// 	// if (fclose(fp) != 0) {
	// 	// 	msg("Failed to close file \"%s\".", fn);
	// 	// 	return false
	// 	// }

	// 	curbp->b_flags &= (modflag ? B_MODIFIED : ~B_MODIFIED);
	// 	msg("File \"%s\" %ld bytes read.", fn, len);
	return nil
}
