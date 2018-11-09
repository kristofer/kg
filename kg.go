// Package kg is a very small emacs editor written in Go.
// ported from https://github.com/hughbarney/atto
package kg

/* header.h, Atto Emacs, Public Domain, Hugh Barney, 2016, Derived from: Anthony's Editor January 93 */
//  _XOPEN_SOURCE
// #include <locale.h>
// #include <stdlib.h>
// #include <stdarg.h>
// #include <assert.h>
// #include <curses.h>
// #include <stdio.h>
// #include <sys/types.h>
// #include <ctype.h>
// #include <limits.h>
// #include <string.h>
// #include <unistd.h>
// #include <wchar.h>

//int mkstemp(char *);

/*
 * Some compilers define size_t as a unsigned 16 bit number while
 * Point and off_t might be defined as a signed 32 bit number.
 * malloc(), realloc(), fread(), and fwrite() take size_t parameters,
 * which means there will be some size limits because size_t is too
 * small of a type.
 */

//MAX_SIZE_T      ((unsigned long) (size_t) ~0)

/*
 * routines to still port..
 */

//  DONE void fatal(char *);
//  DONE void msg(char *, ...);
//  void display(window_t *, int);
//  void dispmsg(void);
//  void modeline(window_t *);
//  int utf8_size(char_t);
//  int prev_utf8_char_size(void);
//  void display_utf8(buffer_t *, char_t, int);
//  Point lnstart(buffer_t *, point_t);
//  Point lncolumn(buffer_t *, point_t, int);
//  Point segstart(buffer_t *, point_t, point_t);
//  Point segnext(buffer_t *, point_t, point_t);
//  Point upup(buffer_t *, point_t);
//  Point dndn(buffer_t *, point_t);
//  char_t *get_key(keymap_t *, keymap_t **);
//  int getinput(char *, char *, int, int);
//  int getfilename(char *, char *, int);
//  void display_prompt_and_response(char *, char *);
//  int growgap(buffer_t *, point_t);
//  Point movegap(buffer_t *, point_t);
//  Point pos(buffer_t *, char_t *);
//  char_t *ptr(buffer_t *, point_t);
//  int posix_file(char *);
//  int save(char *);
//  int load_file(char *);
//  int insert_file(char *, int);
//  void backsp(void);
//  void block(void);
//  void iblock(void);
//  void bottom(void);
//  void cut(void);
//  void copy(void);
//  void copy_cut(int);
//  void delete(void);
//  void toggle_overwrite_mode(void);
//  void down(void);
//  void insert(void);
//  void left(void);
//  void lnbegin(void);
//  void lnend(void);
//  void paste(void);
//  void pgdown(void);
//  void pgup(void);
//  void quit(void);
//  int yesno(int);
//  void quit_ask(void);
//  void redraw(void);
//  void readfile(void);
//  void insertfile(void);
//  void right(void);
//  void top(void);
//  void up(void);
//  void version(void);
//  void wleft(void);
//  void wright(void);
//  void writefile(void);
//  void savebuffer(void);
//  void showpos(void);
//  void killtoeol(void);
//  void gotoline(void);
//  void search(void);
//  void query_replace(void);
//  Point line_to_point(int);
//  Point search_forward(buffer_t *, point_t, char *);
//  Point search_backwards(buffer_t *, point_t, char *);
//  void update_search_prompt(char *, char *);
//  void display_search_result(point_t, int, char *, char *);
//  buffer_t* find_buffer(char *, int);
//  void buffer_init(buffer_t *);
//  int delete_buffer(buffer_t *);
//  void next_buffer(void);
//  int count_buffers(void);
//  int modified_buffers(void);
//  void killbuffer(void);
//  char* get_buffer_name(buffer_t *);
//  void get_line_stats(int *, int *);
//  void query_replace(void);
//  DONE window_t *new_window();
//  DONE void one_window(window_t *);
//  DONE void split_window();
//  DONE void next_window();
//  DONE void delete_other_windows();
//  DONE void free_other_windows();
//  DONE void update_display();
//  void w2b(window_t *);
//  void b2w(window_t *);
//  DONE void associate_b2w(buffer_t *, window_t *);
//  DONE void disassociate_b(window_t *);
//  void set_parse_state(buffer_t *, point_t);
//  void set_parse_state2(buffer_t *, point_t);
//  int parse_text(buffer_t *, point_t);
//  void resize_terminal();
