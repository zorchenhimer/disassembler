package config

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	EOF_Item lexItem = lexItem{typ: lex_EOF, val:""}
	EOF_Rune rune = '\u0004' // EOT ascii control character
)

type stateFn func(l *lexer) stateFn

type lexer struct {
	input []byte
	//start int
	//pos   int
	//width int
	//line  int
	//col   int

	current  runeinfo
	previous runeinfo

	lineStart int
	lastLineStart int
	lineDelta int

	items chan lexItem
}

type runeinfo struct {
	start int
	pos int
	width int
	line int
	col int

	lineStart int
	lineDelta int
}

func newLexer(raw []byte) (*lexer, chan lexItem) {
	l := &lexer{
		input: raw,
		items: make(chan lexItem),
		//line:  1,

		current: runeinfo{
			line: 1,
		},
		previous: runeinfo{
			line: 1,
		},
	}

	return l, l.items
}

func (l *lexer) Run() {
	for state := lexStateStart; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) emit(t lexItemType) {
	itm := lexItem{
		typ:  t,
		val:  string(l.input[l.current.start:l.current.pos]),
		line: l.current.line,
		col:  l.current.start -  l.current.lineStart + 1,
		pos:  l.current.pos,
		start: l.current.start,
		lineStart: l.current.lineStart,
	}

	l.items <- itm
	l.current.start = l.current.pos
	l.current.line += l.current.lineDelta
	l.current.lineDelta = 0
}

func (l *lexer) emitError(message string) {
	itm := lexItem{
		typ:  lex_Error,
		val:  message,
		line: l.current.line,
		col:  l.current.start - l.current.lineStart + 1,
		pos:  l.current.start,
	}
	l.items <- itm
	l.current.start = l.current.pos
	l.current.line += l.current.lineDelta
	l.current.lineDelta = 0
}

func (l *lexer) next() rune {
	l.previous = l.current

	if l.current.pos >= len(l.input) {
		l.current.width = 0
		return EOF_Rune
	}

	r, size := utf8.DecodeRune(l.input[l.current.pos:])
	if size == 0 || r == utf8.RuneError {
		panic(fmt.Sprintf("Invalid unicode at %d", l.current.pos))
	}

	//fmt.Printf("%q\n", r)

	l.current.width = size
	l.current.pos += l.current.width

	if r == '\n' {
		l.current.lineDelta++
		l.current.lineStart = l.current.pos
	}

	return r
}

func (l *lexer) peek() rune {
	if l.current.pos >= len(l.input) {
		l.current.width = 0
		return EOF_Rune
	}

	r, size := utf8.DecodeRune(l.input[l.current.pos:])
	if size == 0 || r == utf8.RuneError {
		panic(fmt.Sprintf("Invalid unicode at %d", l.current.pos))
	}

	return r
}

func (l *lexer) backup() {
	l.current = l.previous
	//l.pos -= l.width
	//if l.lineDelta > 0 {
	//	l.lineDelta--
	//}
	//l.lineStart = l.lastLineStart
}

func (l *lexer) ignore() {
	l.current.start = l.current.pos
	l.current.line += l.current.lineDelta
	l.current.lineDelta = 0
}

func (l *lexer) discard() {
	l.next()
	l.current.start = l.current.pos
	l.current.line += l.current.lineDelta
	l.current.lineDelta = 0
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}

	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for {
		r := l.next()
		if r == EOF_Rune {
			break
		}

		if strings.IndexRune(valid, r) < 0 {
			break
		}
	}
	l.backup()
}

type lexItemType int

const (
	lex_EOF lexItemType = iota

	lex_OpenBracket
	lex_CloseBracket
	lex_OpenSquare
	lex_CloseSquare

	lex_Ident
	lex_String
	lex_Semicolon
	lex_Number

	lex_Comment
	lex_Space
	lex_Error
)

func (lit lexItemType) String() string {
	switch lit {
	case lex_EOF:
		return "lex_EOF"
	case lex_OpenBracket:
		return "lex_OpenBracket"
	case lex_CloseBracket:
		return "lex_CloseBracket"
	case lex_OpenSquare:
		return "lex_OpenSquare"
	case lex_CloseSquare:
		return "lex_CloseSquare"
	case lex_Ident:
		return "lex_Ident"
	case lex_String:
		return "lex_String"
	case lex_Semicolon:
		return "lex_Semicolon"
	case lex_Number:
		return "lex_Number"
	case lex_Space:
		return "lex_Space"
	case lex_Error:
		return "lex_Error"
	case lex_Comment:
		return "lex_Comment"
	}

	return "lex_UNKNOWN"
}

type lexItem struct {
	typ lexItemType
	val string
	pos int

	line int
	col  int
	start int
	lineStart int
}

func (itm lexItem) GoString() string {
	return fmt.Sprintf("{%s [%d:%d {%d:%d} <%d>] %q}", itm.typ, itm.line, itm.col, itm.start, itm.pos, itm.lineStart, itm.val)
}

func (itm lexItem) String() string {
	return fmt.Sprintf("[%d:%d] %s %q", itm.line, itm.col, itm.typ, itm.val)
}

func lexComment(l *lexer) stateFn {
	l.lineStart -= 1
	for {
		r := l.next()
		if r == EOF_Rune {
			l.emit(lex_Comment)
			return nil
		}

		if r == '\n' {
			l.backup()
			l.emit(lex_Comment)
			break
		}
	}

	return lexStateStart
}

func lexIdent(l *lexer) stateFn {
	l.acceptRun("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	l.emit(lex_Ident)
	return lexStateStart
}

func lexNumber(l *lexer) stateFn {
	l.accept("-+")
	l.accept("x$%")
	l.acceptRun("0123456789ABCDEFabcdef_") // accept underscores for visual separation
	l.emit(lex_Number)
	return lexStateStart
}

func lexString(l *lexer) stateFn {
	for {
		r := l.next()
		if r == '\\' && l.peek() == '"' {
			l.next()
			continue
		}

		if r == '"' {
			break
		}

		if r == EOF_Rune {
			l.emitError("EOF in string")
			return nil
		}
	}

	l.backup()
	l.emit(lex_String)
	l.discard()
	return lexStateStart
}

func lexStateStart(l *lexer) stateFn {
	for {
		r := l.next()
		if r == EOF_Rune {
			break
		}

		if unicode.IsSpace(r) {
			l.ignore()
			continue
		}

		switch r {
		case '[':
			l.emit(lex_OpenSquare)
			continue
		case ']':
			l.emit(lex_CloseSquare)
			continue
		case '{':
			l.emit(lex_OpenBracket)
			continue
		case '}':
			l.emit(lex_CloseBracket)
			continue

		case ';':
			l.emit(lex_Semicolon)
			continue

		case '/':
			r = l.next()
			if r != '/' {
				l.emitError("Invalid slash")
				return nil
			}
			return lexComment

		case '"':
			l.ignore()
			return lexString
		}

		if r == '$' || r == '%' || r == '+' || r == '-' || ('0' <= r && r <= '9') {
			l.backup()
			return lexNumber
		}

		if ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') {
			l.backup()
			return lexIdent
		}
	}

	if l.current.pos > l.current.start {
		l.emit(lex_Space)
	}

	l.emit(lex_EOF)
	return nil
}
