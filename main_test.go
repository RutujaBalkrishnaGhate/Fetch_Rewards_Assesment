package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestProcessReceipt_ValidInput tests the processReceipt function with valid input
func TestProcessReceipt_ValidInput(t *testing.T) {
	router := setupRouter() 

	receipt := Receipt{
		Retailer:     "Test Store",
		PurchaseDate: "2024-08-01",
		PurchaseTime: "14:30",
		Items: []Item{
			{ShortDescription: "Item1", Price: "10.00"},
			{ShortDescription: "Item2", Price: "15.50"},
		},
		Total: "25.50",
	}
	jsonData, _ := json.Marshal(receipt)

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "id")
}

// TestProcessReceipt_InvalidDate tests the processReceipt function with an invalid date
func TestProcessReceipt_InvalidDate(t *testing.T) {
	router := setupRouter() 

	receipt := Receipt{
		Retailer:     "Test Store",
		PurchaseDate: "InvalidDate",
		PurchaseTime: "14:30",
		Items:        []Item{},
		Total:        "25.50",
	}
	jsonData, _ := json.Marshal(receipt)

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid purchase date format")
}

// TestProcessReceipt_InvalidTime tests the processReceipt function with an invalid time
func TestProcessReceipt_InvalidTime(t *testing.T) {
	router := setupRouter() 

	receipt := Receipt{
		Retailer:     "Test Store",
		PurchaseDate: "2024-08-01",
		PurchaseTime: "InvalidTime",
		Items:        []Item{},
		Total:        "25.50",
	}
	jsonData, _ := json.Marshal(receipt)

	req, _ := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusBadRequest, resp.Code)
	assert.Contains(t, resp.Body.String(), "Invalid purchase time format")
}

// TestGetPoints_ValidReceipt tests the getPoints function with a valid receipt
func TestGetPoints_ValidReceipt(t *testing.T) {
	router := setupRouter() 

	receipt := Receipt{
		Retailer:     "Test Store",
		PurchaseDate: "2024-08-01",
		PurchaseTime: "14:30",
		Items: []Item{
			{ShortDescription: "Item1", Price: "10.00"},
		},
		Total: "10.00",
	}
	id := "test-id" 
	receipts[id] = receipt

	req, _ := http.NewRequest("GET", "/receipts/"+id+"/points", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Contains(t, resp.Body.String(), "points")
}

// TestGetPoints_ReceiptNotFound tests the getPoints function with an invalid receipt ID
func TestGetPoints_ReceiptNotFound(t *testing.T) {
	router := setupRouter() 
	req, _ := http.NewRequest("GET", "/receipts/invalid-id/points", nil)
	resp := httptest.NewRecorder()
	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
	assert.Contains(t, resp.Body.String(), "Receipt not found")
}
