package mock

import (
	"fmt"
)

func Parse(input string) bool {
	lexer := lex(input)
	return expression(lexer)
}

func expression(lex *lexer) bool {
	val := factor(lex)

	for {
		item := lex.nextItem()

		switch item.Type {
		case ItemEOF:
			// end
			return val

		case ItemIdent:
			switch item.Val {
			case "AND":
				// condition AND condition
				val2 := factor(lex)
				val = val && val2
			case "OR":
				// condition OR condition
				val2 := factor(lex)
				val = val || val2
			default:
				return val
			}

		default:
			handleError(fmt.Errorf("dynamo: expression error at %v", item))
		}
	}
}

func factor(lex *lexer) bool {
	itemA := lex.nextItem()

	switch itemA.Type {
	case ItemOperand:
		// operand comparator operand, operand BETWEEN operand AND operand,
		itemB := lex.nextItem()

		if itemB.Type == ItemComparator {
			// operand comparator operand
			itemC := lex.nextItem()
			if itemC.Type != ItemOperand {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemC))
			}
			fmt.Printf("%s %s %s\n", itemA.Val, itemB.Val, itemC.Val)
			return true

		} else if itemB.Val == "BETWEEN" {
			// operand BETWEEN operand AND operand,
			itemC := lex.nextItem()
			if itemC.Type != ItemOperand {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemC))
			}
			itemD := lex.nextItem()
			if itemD.Val != "AND" {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemD))
			}
			itemE := lex.nextItem()
			if itemE.Type != ItemOperand {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemE))
			}
			fmt.Printf("%s %s %s %s %s\n", itemA.Val, itemB.Val, itemC.Val, itemD.Val, itemE.Val)
			return true

		} else if itemB.Val == "IN" {
			// operand IN ( operand (',' operand (, ...) ))
			// TODO:
		} else {
			handleError(fmt.Errorf("dynamo: invalid factor at %v", itemB))
		}

	case ItemFunction:
		switch itemA.Val {
		case "attribute_exists(":
			// attribute_exists (path)
			itemB := lex.nextItem()
			if itemB.Type != ItemOperand {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemB))
			}
			itemC := lex.nextItem()
			if itemC.Val != ")" {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemC))
			}
			fmt.Printf("%s %s %s\n", itemA.Val, itemB.Val, itemC.Val)
			return true

		case "attribute_not_exists(":
			// attribute_not_exists (path)
			itemB := lex.nextItem()
			if itemB.Type != ItemOperand {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemB))
			}
			itemC := lex.nextItem()
			if itemC.Val != ")" {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemC))
			}
			fmt.Printf("%s %s %s\n", itemA.Val, itemB.Val, itemC.Val)
			return true

		case "attribute_type(":
			// attribute_type (path, type)
			itemB := lex.nextItem()
			if itemB.Type != ItemOperand {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemB))
			}
			itemC := lex.nextItem()
			if itemC.Val != "," {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemC))
			}
			itemD := lex.nextItem()
			if itemD.Type != ItemOperand {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemD))
			}
			itemE := lex.nextItem()
			if itemE.Val != ")" {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemE))
			}
			fmt.Printf("%s %s %s %s %s\n", itemA.Val, itemB.Val, itemC.Val, itemD.Val, itemE.Val)
			return true

		case "begins_with(":
			// begins_with (path, substr)
			// TODO:
		case "contains(":
			// contains (path, operand)
			// TODO:
		case "size(":
			// size (path)
			// TODO:
		default:
			handleError(fmt.Errorf("dynamo: invalid factor at %v", itemA))
		}

	case ItemIdent:

		switch itemA.Val {
		case "NOT":
			// NOT condition
			val := factor(lex)
			return !val

		case "(":
			// ( condition )
			val := expression(lex)
			itemB := lex.items[len(lex.items)-1]
			if itemB.Val != ")" {
				handleError(fmt.Errorf("dynamo: invalid factor at %v", itemB))
			}
			fmt.Printf("%s %v %s\n", itemA.Val, val, itemB.Val)
			return val
		}

	default:
		handleError(fmt.Errorf("dynamo: invalid factor at %v", itemA))
	}

	return false
}

func handleError(err error) {
	panic(err)
}
