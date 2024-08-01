package main

import (
	"github.com/gin-gonic/gin"
)

// setupRouter initializes and returns a new Gin engine with defined routes.
// This function sets up the HTTP routes for processing receipts and retrieving points.
func setupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/receipts/process", processReceipt)
	r.GET("/receipts/:id/points", getPoints)

	return r
}
