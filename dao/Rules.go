package dao

import "RuleEngineAST/models"

func CreateRule(r models.Rule) {
	DB.Create(&r)
}

func FindRule() []models.Rule {
	var rule []models.Rule
	DB.Find(&rule)
	return rule
}
