package controller

import (
	"fmt"
	"net/http"

	"RuleEngineAST/service"
	"github.com/gin-gonic/gin"
)

var ruleManager service.RuleInterface = &service.RuleManagerV1{}

var ruleEngine = &RuleEngine{}

func FindRules(c *gin.Context) {
	var rules = ruleManager.FindRules()
	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

func CreateRule(c *gin.Context) {

	type request struct {
		Rule string `json:"rule"`
	}

	req := &request{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	_, err = ruleEngine.parseTree(req.Rule)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid rule. err : %s", err.Error()))
		return
	}

	rule := ruleManager.CreateRule(req.Rule)

	c.JSON(http.StatusOK, rule)
}

func MergeRules(c *gin.Context) {

	type request struct {
		FirstRule  string `json:"first_rule"`
		SecondRule string `json:"second_rule"`
		Strategy   string `json:"merge_strategy"`
	}

	req := &request{}
	err := c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if req.Strategy != "AND" && req.Strategy != "OR" {
		c.JSON(http.StatusBadRequest, "invalid strategy")
		return
	}

	// combine rule has onl 2 strategy for now  i.e AND & OR
	mergedRule := ruleEngine.combineRule(req.FirstRule, req.SecondRule, req.Strategy)
	_, err = ruleEngine.parseTree(mergedRule)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("cannot merge rules. err : %s", err.Error()))
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"merged_rule": mergedRule,
	})
}

func EvaluateRule(c *gin.Context) {

	type payloadStruct struct {
		Rule string            `json:"rule"`
		Data map[string]string `json:"data"`
	}

	payload := &payloadStruct{}

	err := c.BindJSON(payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, fmt.Sprintf("invalid rule. err : %s", err.Error()))
		return
	}

	ast, err := ruleEngine.parseTree(payload.Rule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	evalNode := ruleEngine.evaluateRule(ast, payload.Data)

	c.JSON(http.StatusOK, map[string]bool{
		"rule_match": evalNode.MatchValue,
	})
}
