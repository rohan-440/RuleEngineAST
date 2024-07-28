package controller

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEvaluateRule(t *testing.T) {

	re := NewRuleEngine()

	testCases := []struct {
		desc          string
		ruleString    string
		dataMap       map[string]string
		expectedMatch bool
	}{
		{
			desc:       "rule is match for AND operator",
			ruleString: "age > 30 AND department == 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "31",
				"department": "ENGINEERING",
			},
			expectedMatch: true,
		},
		{
			desc:       "rule is match for OR operator",
			ruleString: "age > 30 OR department == 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "31",
				"department": "SALES",
			},
			expectedMatch: true,
		},
		{
			desc:       "rule is not match for non equal check",
			ruleString: "age > 30 AND department != 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "31",
				"department": "ENGINEERING",
			},
			expectedMatch: false,
		},
		{
			desc:       "rule is not match for == check",
			ruleString: "age > 30 AND department == 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "31",
				"department": "SALES",
			},
			expectedMatch: false,
		},
		{
			desc:       "rule is not match for > check",
			ruleString: "age > 30 AND department == 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "30",
				"department": "ENGINEERING",
			},
			expectedMatch: false,
		},
		{
			desc:       "rule is not match for >= check",
			ruleString: "age > 30 AND department == 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "29",
				"department": "ENGINEERING",
			},
			expectedMatch: false,
		},
		{
			desc:       "rule is not match for < check",
			ruleString: "age < 30 AND department == 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "31",
				"department": "ENGINEERING",
			},
			expectedMatch: false,
		},
		{
			desc:       "rule is not match for <= check",
			ruleString: "age <= 30 AND department == 'ENGINEERING'",
			dataMap: map[string]string{
				"age":        "31",
				"department": "ENGINEERING",
			},
			expectedMatch: false,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {

			//build tree
			ast, err := re.parseTree(tt.ruleString)
			assert.Nil(t, err)

			//match the rule
			evalNode := re.evaluateRule(ast, tt.dataMap)
			assert.Equal(t, evalNode.MatchValue, tt.expectedMatch)
		})
	}
}

func TestParseTree(t *testing.T) {

	re := NewRuleEngine()

	testCases := []struct {
		desc          string
		ruleString    string
		expectedError error
	}{
		{
			desc:       "successfully parse tree",
			ruleString: "age > 30 AND department == 'ENGINEERING'",
		},
		{
			desc:          "invalid rule string",
			ruleString:    "age > AND department == 'ENGINEERING'",
			expectedError: errors.New("error parsing comparison: error parsing: unexpected end of expression\n"),
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			_, err := re.parseTree(tt.ruleString)
			assert.Equal(t, err, tt.expectedError)
		})
	}
}

func TestCombineRule(t *testing.T) {

	re := NewRuleEngine()

	testCases := []struct {
		desc              string
		firstRule         string
		SecondRule        string
		Strategy          string
		ExpectedMergeRule string
	}{
		{
			desc:              "successfully merge rules",
			firstRule:         "age > 30",
			SecondRule:        "department == 'ENGINEERING'",
			Strategy:          "AND",
			ExpectedMergeRule: "(age > 30) AND (department == 'ENGINEERING')",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			mergeRule := re.combineRule(tt.firstRule, tt.SecondRule, tt.Strategy)
			assert.Equal(t, tt.ExpectedMergeRule, mergeRule)
		})
	}
}
