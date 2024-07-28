package controller

import (
	"fmt"
	"strconv"
	"strings"

	"RuleEngineAST/ast/parse"
	bools "RuleEngineAST/ast/parse/bool"
	"RuleEngineAST/ast/parse/comp"
)

type RuleEngine struct {
}

func NewRuleEngine() *RuleEngine {
	return &RuleEngine{}
}

type EvaluateNode struct {
	Key        string
	MatchValue bool
}

func (re *RuleEngine) parseTree(ruleString string) (parse.AST, error) {
	bParser, _ := bools.NewParser()
	cParser, _ := comp.NewParser()

	ast, err := bParser.ParseStr(ruleString)
	if err != nil {
		return nil, err
	}

	// parse comparisons
	err = ast.Parse(cParser)
	if err != nil {
		return nil, fmt.Errorf("error parsing comparison: %v\n", err)
	}

	return ast, nil
}

// TODO: optimise combine rule method using AST if time available
func (re *RuleEngine) combineRule(r1, r2, strategy string) string {
	return fmt.Sprintf("(%s) %s (%s)", r1, strategy, r2)
}

func (re *RuleEngine) evaluateRule(ast parse.AST, dataMap map[string]string) *EvaluateNode {
	switch ast := (ast).(type) {
	case *bools.BinExpr:
		leftEval := re.evaluateRule(ast.LHS, dataMap)
		rightEval := re.evaluateRule(ast.RHS, dataMap)

		switch ast.Op.String() {
		case "AND":
			return &EvaluateNode{MatchValue: leftEval.MatchValue && rightEval.MatchValue}
		case "OR":
			return &EvaluateNode{MatchValue: leftEval.MatchValue || rightEval.MatchValue}
		default:
			return &EvaluateNode{MatchValue: false}
		}

	case *comp.EqualExpr:
		leftEval := re.evaluateRule(ast.LHS, dataMap)
		rightEval := re.evaluateRule(ast.RHS, dataMap)

		val, ok := dataMap[leftEval.Key]
		if !ok {
			return &EvaluateNode{MatchValue: false}
		}

		switch ast.Op {
		case comp.OpEqual:
			return &EvaluateNode{MatchValue: rightEval.Key == val}
		case comp.OpNotEqual:
			return &EvaluateNode{MatchValue: rightEval.Key != val}
		}

	case *comp.OrdinalExpr:
		leftEval := re.evaluateRule(ast.LHS, dataMap)
		rightEval := re.evaluateRule(ast.RHS, dataMap)

		val, ok := dataMap[leftEval.Key]
		if !ok {
			return &EvaluateNode{MatchValue: false}
		}

		ruleData, err := strconv.ParseFloat(rightEval.Key, 32)
		if err != nil {
			return &EvaluateNode{MatchValue: false}
		}

		dataVal, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return &EvaluateNode{MatchValue: false}
		}

		switch ast.Op {

		case comp.OpGreater:
			return &EvaluateNode{MatchValue: dataVal > ruleData}

		case comp.OpGreaterOrEqual:
			return &EvaluateNode{MatchValue: dataVal >= ruleData}

		case comp.OpLess:
			return &EvaluateNode{MatchValue: dataVal < ruleData}

		case comp.OpLessOrEqual:
			return &EvaluateNode{MatchValue: dataVal <= ruleData}

		}

	case parse.Unparsed:
		key := strings.Join(ast.Contents, " ")
		key = strings.ReplaceAll(key, "'", "")
		return &EvaluateNode{Key: key, MatchValue: true}
	}

	return nil
}
