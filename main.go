package main

import (
	"net/http"

	"RuleEngineAST/controller"
	"RuleEngineAST/dao"

	"github.com/gin-gonic/gin"
)

func main() {

	dao.ConnectDatabase("db")
	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//get all rules
	router.GET("/rules", controller.FindRules)

	//create a new rule
	router.POST("/rules", controller.CreateRule)

	//evaluate a rule with data
	router.POST("/rules/evaluate", controller.EvaluateRule)

	//merge rules
	router.POST("/rules/merge", controller.MergeRules)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
