package mock

import (
	"reflect"
	"testing"
)

func TestLexer(t *testing.T) {

	tests := []struct {
		name  string
		input string
		out   []Item
	}{
		{
			name:  "case1",
			input: "1 < あ",
			out: []Item{
				{ItemOperand, 0, "1"},
				{ItemComparator, 2, "<"},
				{ItemOperand, 4, "あ"},
			},
		},
		{
			name:  "case2",
			input: " attribute_exists(path) ",
			out: []Item{
				{ItemFunction, 1, "attribute_exists("},
				{ItemOperand, 18, "path"},
				{ItemIdent, 22, ")"},
			},
		},
		{
			name:  "case3",
			input: "1 <> operand AND  NOT op",
			out: []Item{
				{ItemOperand, 0, "1"},
				{ItemComparator, 2, "<>"},
				{ItemOperand, 5, "operand"},
				{ItemIdent, 13, "AND"},
				{ItemIdent, 18, "NOT"},
				{ItemOperand, 22, "op"},
			},
		},
		{
			name:  "case4",
			input: "(contains(path, op1) OR op2 BETWEEN op3 AND op4) NOT op5 IN (op6,op7) ",
			out: []Item{
				{ItemIdent, 0, "("},
				{ItemFunction, 1, "contains("},
				{ItemOperand, 10, "path"},
				{ItemIdent, 14, ","},
				{ItemOperand, 16, "op1"},
				{ItemIdent, 19, ")"},
				{ItemIdent, 21, "OR"},
				{ItemOperand, 24, "op2"},
				{ItemIdent, 28, "BETWEEN"},
				{ItemOperand, 36, "op3"},
				{ItemIdent, 40, "AND"},
				{ItemOperand, 44, "op4"},
				{ItemIdent, 47, ")"},
				{ItemIdent, 49, "NOT"},
				{ItemOperand, 53, "op5"},
				{ItemIdent, 57, "IN"},
				{ItemIdent, 60, "("},
				{ItemOperand, 61, "op6"},
				{ItemIdent, 64, ","},
				{ItemOperand, 65, "op7"},
				{ItemIdent, 68, ")"},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				lexer = lex(test.input)
				got   = make([]Item, 0, len(test.out))
			)
			var counter int
			for {
				item := lexer.nextItem()
				if item.Type == ItemEOF {
					break
				}
				got = append(got, item)
				counter++
				if counter > 100 {
					break
				}
			}
			if !reflect.DeepEqual(test.out, got) {
				t.Error("bad value. want:", test.out, "got:", got)
			}
		})
	}

}
