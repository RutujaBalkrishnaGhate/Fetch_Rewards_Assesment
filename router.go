package main

import (
	"github.com/gin-gonic/gin"
)


func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/receipts/process", processReceipt)
	r.GET("/receipts/:id/points", getPoints)

	return r
}
