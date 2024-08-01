package main

import (
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
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

var (
	receipts = make(map[string]Receipt)
	mu       sync.Mutex
)

func main() {
	r := gin.Default()

	r.POST("/receipts/process", processReceipt)
	r.GET("/receipts/:id/points", getPoints)

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func processReceipt(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	var receipt Receipt

	if err := c.ShouldBindJSON(&receipt); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input format"})
		return
	}

	// Validate purchase date and time formats
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
	mu.Lock()
	defer mu.Unlock()

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

	// One point for every alphanumeric character in the retailer name
	points += countAlphanumeric(receipt.Retailer)

	// 50 points if the total is a round dollar amount with no cents
	total, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		log.Printf("Error parsing total amount: %v", err)
		return 0
	}
	if total == float64(int(total)) {
		points += 50
	}

	// 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	// 5 points for every two items on the receipt
	points += (len(receipt.Items) / 2) * 5

	// Points for item descriptions
	for _, item := range receipt.Items {
		trimmedDesc := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDesc)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				log.Printf("Error parsing item price: %v", err)
				continue
			}
			points += int(math.Ceil(price * 0.2))
		}
	}

	// 6 points if the day in the purchase date is odd
	date, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err != nil {
		log.Printf("Error parsing purchase date: %v", err)
		return points
	}
	if date.Day()%2 != 0 {
		points += 6
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err != nil {
		log.Printf("Error parsing purchase time: %v", err)
		return points
	}
	if purchaseTime.Hour() >= 14 && purchaseTime.Hour() < 16 {
		points += 10
	}

	return points
}

func countAlphanumeric(str string) int {
	re := regexp.MustCompile("[a-zA-Z0-9]")
	return len(re.FindAllString(str, -1))
}
