package mock

import (
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {

	tests := []struct {
		name  string
		input string
		out   bool
	}{
		{
			name:  "case1",
			input: "1 < あ",
			out:   true,
		},
		{
			name:  "case2",
			input: " attribute_exists(path) ",
			out:   true,
		},
		{
			name:  "case3",
			input: "1 <> operand OR NOT op1 <= op2",
			out:   true,
		},
		{
			name:  "case4",
			input: "(attribute_type(path, op1) OR op2 BETWEEN op3 AND op4) AND NOT attribute_not_exists(op6) ",
			out:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := Parse(test.input)
			if !reflect.DeepEqual(test.out, got) {
				t.Error("bad value. want:", test.out, "got:", got)
			}
		})
	}

}
