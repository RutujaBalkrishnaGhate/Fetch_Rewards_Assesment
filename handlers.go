package main

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var receipts = make(map[string]Receipt)

func processReceipt(c *gin.Context) {
	var receipt Receipt

	if err := c.ShouldBindJSON(&receipt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	if _, err := time.Parse("2006-01-02", receipt.PurchaseDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase date format"})
		return
	}

	if _, err := time.Parse("15:04", receipt.PurchaseTime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase time format"})
		return
	}

	id := uuid.New().String()
	receipts[id] = receipt

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func getPoints(c *gin.Context) {
	id := c.Param("id")
	receipt, exists := receipts[id]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Receipt not found"})
		return
	}

	points := calculatePoints(receipt)
	c.JSON(http.StatusOK, gin.H{"points": points})
}

func calculatePoints(receipt Receipt) int {
	points := 0

	points += countAlphanumeric(receipt.Retailer)

	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.Items) / 2) * 5

	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
	if date.Day()%2 != 0 {
		points += 6
	}

	time, _ := time.Parse("15:04", receipt.PurchaseTime)
	if time.Hour() >= 14 && time.Hour() < 16 {
		points += 10
	}

	return points
}

func countAlphanumeric(str string) int {
	re := regexp.MustCompile("[a-zA-Z0-9]")
	return len(re.FindAllString(str, -1))
}
