package parse

import (
	"errors"
	"fmt"
	"unicode"
)

var ErrConfig = errors.New("config error")

var ErrParse = errors.New("error parsing")

var ErrEval = errors.New("eval error")

var ErrUnknownAST = errors.New("unknown AST node")

type AST interface {
	Parse(Parser) error
}

type Unparsed struct {
	Contents []string
}

func (u Unparsed) Parse(p Parser) error {
	return fmt.Errorf("%w: attempted to parse Unparsed node", ErrParse)
}

type Parser interface {
	Parse(tokens []string) (AST, error)
}

type Interpreter[T any] func(AST) (T, error)

func (i Interpreter[T]) WithFallback(b Interpreter[T]) Interpreter[T] {
	return func(ast AST) (T, error) {
		a, err := i(ast)
		if errors.Is(err, ErrUnknownAST) {
			return b(ast)
		}
		return a, err
	}
}

func Tokenize(str string, open, close rune, keywordMatcher *KeywordTrie) []string {
	runes := []rune(str)
	var substr []rune
	var result []string
	push := func() { // push substr onto result
		result = append(result, string(substr))
		substr = nil
	}

	for i := 0; i < len(runes); i++ {
		if runes[i] == open || runes[i] == close {
			if len(substr) > 0 {
				push()
			}
			result = append(result, string(runes[i]))
			continue
		}
		if unicode.IsSpace(runes[i]) {
			if len(substr) > 0 {
				push()
			}
			continue
		}
		matched := keywordMatcher.Match(runes[i:])
		if len(matched) > 0 {
			if len(substr) > 0 {
				push()
			}
			result = append(result, matched)
			i += len(matched) - 1
		} else {
			substr = append(substr, runes[i])
		}
	}
	if len(substr) > 0 {
		push()
	}
	return result
}
