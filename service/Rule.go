package service

import (
	"time"

	"RuleEngineAST/dao"
	"RuleEngineAST/models"
)

type RuleInterface interface {
	FindRules() []models.Rule
	CreateRule(ruleStr string) models.Rule
}

type RuleManagerV1 struct {
}

func (ruleManager *RuleManagerV1) FindRules() []models.Rule {
	return dao.FindRule()
}

func (ruleManager *RuleManagerV1) CreateRule(ruleStr string) models.Rule {
	rule := models.Rule{
		Rule:      ruleStr,
		CreatedAt: time.Now(),
	}

	dao.CreateRule(rule)

	return rule
}
