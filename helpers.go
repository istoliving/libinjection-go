package libinjection

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

func file_line() int {
	_, fileName, fileLine, ok := runtime.Caller(1)
	var s string
	if ok {
		s = fmt.Sprintf("%s:%d", fileName, fileLine)
	} else {
		s = ""
	}
	l, _ := strconv.Atoi(s)
	return l
}

func streq(s1 []byte, s2 string) bool {
	if clen(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}

/**
 * Initializes parsing state
 *
 */
func flag2delim(flag int) byte {
	if flag&FLAG_QUOTE_SINGLE == 1 {
		return CHAR_SINGLE
	} else if flag&FLAG_QUOTE_DOUBLE == 1 {
		return CHAR_DOUBLE
	} else {
		return CHAR_NULL
	}
}

func st_clear(t **stoken_t) {
	*t = new(stoken_t) //nil
}

/** Find largest string containing certain characters.
 *
 * C Standard library 'strspn' only works for 'c-strings' (null terminated)
 * This works on arbitrary length.
 *
 * Performance notes:
 *   not critical
 *
 * Porting notes:
 *   if accept is 'ABC', then this function would be similar to
 *   a_regexp.match(a_str, '[ABC]*'),
 */
func strlenspn(s []byte, l int, accept string) int {
	for i := 0; i < l; i++ {
		if !strings.ContainsRune(accept, rune(s[i])) {
			return i
		}
	}
	return l
}

func strlencspn(s []byte, l int, accept string) int {
	for i := 0; i < l; i++ {
		if strings.ContainsRune(accept, rune(s[i])) {
			return i
		}
	}
	return l
}

func char_is_white(ch byte) bool {
	/* ' '  space is 0x32
	   '\t  0x09 \011 horizontal tab
	   '\n' 0x0a \012 new line
	   '\v' 0x0b \013 vertical tab
	   '\f' 0x0c \014 new page
	   '\r' 0x0d \015 carriage return
	        0x00 \000 null (oracle)
	        0xa0 \240 is Latin-1
	*/
	b := []byte{
		' ',
		'\t',
		'\n',
		'\v',
		'\f',
		'\r',
		'\240',
		'\000',
	}
	return bytes.ContainsRune(b, rune(ch))
}

/* DANGER DANGER
 * This is -very specialized function-
 *
 * this compares a ALL_UPPER CASE C STRING
 * with a *arbitrary memory* + length
 *
 * Sane people would just make a copy, up-case
 * and use a hash table.
 *
 * Required since libc version uses the current locale
 * and is much slower.
 */

func cstrcasecmp(a string, b []byte, n int) int {
	i := 0
	for n > 0 {
		cb := b[i]
		if cb >= 'a' && cb <= 'z' {
			cb -= 0x20
		}
		if a[i] != cb {
			return int(a[i] - cb)
		} else if a[i] == 0x00 {
			return -1
		}
		i++
		n--
	}
	if a[i] == 0 {
		return 0
	}
	return 1
}

// Copy an interface into another
func st_copy(s1 *stoken_t, s2 *stoken_t) {
	s1.Pos = s2.Pos
	s1.Len = s2.Len
	s1.Count = s2.Count
	s1.Type = s2.Type
	s1.StrOpen = s2.StrOpen
	s1.StrClose = s2.StrClose
	s1.Val = s2.Val
}

func clen(b []byte) int {
	l := len(b)
	for i := 0; ; i++ {
		if l >= i || b[i] == 0x00 {
			return i
		}
	}
}
