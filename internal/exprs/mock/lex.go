package mock

import (
	"strings"
	"unicode/utf8"
)

type ItemType int

// Syntax for filter and condition expressions
// https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Expressions.OperatorsAndFunctions.html
const (
	ItemError ItemType = iota
	ItemEOF
	ItemIdent
	ItemComparator
	ItemFunction
	ItemOperand
)

const (
	eof = -1
)

type Item struct {
	Type ItemType
	Pos  int
	Val  string
}

type lexer struct {
	input string
	start int
	pos   int
	width int
	items []Item
}

func lex(input string) *lexer {
	l := &lexer{
		input: input,
		items: []Item{},
	}
	return l
}

func (l *lexer) nextItem() (i Item) {
	i.Type, i.Pos, i.Val = l.next()
	l.items = append(l.items, i)
	return
}

func (l *lexer) next() (ItemType, int, string) {
	defer func() {
		l.pos += l.width
		l.width = 0
	}()

	var (
		counter  int
		val      string
		itemType ItemType
	)
	for {
		// eof
		if l.pos >= len(l.input) {
			return ItemEOF, l.pos, ""
		}

		// skip prefix blank
		if l.input[l.pos] == ' ' {
			l.pos++
			continue
		}

		// overflow
		if (l.pos + counter) >= len(l.input) {
			itemType = ItemOperand
			break
		}

		// take one utf8 string
		r, w := utf8.DecodeRuneInString(l.input[l.pos+counter:])
		l.width += w
		val += string(r)

		switch val {
		case "= ", "<> ", "< ", "<= ", "> ", ">= ":
			itemType = ItemComparator
		case "(", ")", "NOT ", "AND ", "OR ", "BETWEEN ", "IN ", ",":
			itemType = ItemIdent
		case "attribute_exists(", "attribute_not_exists(", "attribute_type(", "begins_with(", "contains(", "size(":
			itemType = ItemFunction
		}

		if itemType != ItemError {
			break
		}

		if r == ')' || r == ' ' || r == ',' {
			itemType = ItemOperand
			l.width--
			val = val[:len(val)-1]
			break
		}

		counter += w
	}

	return itemType, l.pos, strings.TrimSpace(val)
}
