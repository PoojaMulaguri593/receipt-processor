package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/google/uuid"
)

// Item represents a single item in a receipt.
type Item struct {
	Description string `json:"shortDescription"`
	Price       string `json:"price"`
}

// Receipt holds the details of a purchase receipt.
type Receipt struct {
	StoreName      string `json:"retailer"`
	DateOfPurchase string `json:"purchaseDate"`
	TimeOfPurchase string `json:"purchaseTime"`
	TotalAmount    string `json:"total"`
	PurchasedItems []Item `json:"items"`
}

// ReceiptResponse represents the response containing the receipt ID.
type ReceiptResponse struct {
	ReceiptID string `json:"id"`
}

// PointsResponse holds the calculated points for a receipt.
type PointsResponse struct {
	EarnedPoints int `json:"points"`
}

var (
	receiptStorage = make(map[string]Receipt)
	storageMutex   = &sync.Mutex{}
)

// processReceipt handles the processing and storage of receipts.
func processReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
		http.Error(w, "Invalid receipt format. Please verify input.", http.StatusBadRequest)
		return
	}

	// Validate receipt fields
	if receipt.StoreName == "" || receipt.DateOfPurchase == "" || receipt.TimeOfPurchase == "" || receipt.TotalAmount == "" || len(receipt.PurchasedItems) == 0 {
		http.Error(w, "Invalid receipt format. Please verify input.", http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()
	storageMutex.Lock()
	receiptStorage[receiptID] = receipt
	storageMutex.Unlock()

	json.NewEncoder(w).Encode(ReceiptResponse{ReceiptID: receiptID})
}

// getPoints retrieves the calculated points for a given receipt ID.
func getPoints(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	receiptID := strings.TrimPrefix(r.URL.Path, "/receipts/")
	receiptID = strings.TrimSuffix(receiptID, "/points")

	// Validate receipt ID format
	if !regexp.MustCompile(`^\S+$`).MatchString(receiptID) {
		http.Error(w, "Invalid receipt ID format", http.StatusBadRequest)
		return
	}

	storageMutex.Lock()
	receipt, exists := receiptStorage[receiptID]
	storageMutex.Unlock()
	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	points := computePoints(receipt)
	json.NewEncoder(w).Encode(PointsResponse{EarnedPoints: points})
}

// computePoints calculates the points earned based on the receipt details.
func computePoints(receipt Receipt) int {
	points := 0

	for _, char := range receipt.StoreName {
		if isAlphanumeric(char) {
			points++
		}
	}

	if strings.HasSuffix(receipt.TotalAmount, ".00") {
		points += 50
	}

	totalValue, _ := strconv.ParseFloat(receipt.TotalAmount, 64)
	if math.Mod(totalValue, 0.25) == 0 {
		points += 25
	}

	points += (len(receipt.PurchasedItems) / 2) * 5

	for _, item := range receipt.PurchasedItems {
		if len(strings.TrimSpace(item.Description))%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	dateParts := strings.Split(receipt.DateOfPurchase, "-")
	day, _ := strconv.Atoi(dateParts[2])
	if day%2 != 0 {
		points += 6
	}

	timeParts := strings.Split(receipt.TimeOfPurchase, ":")
	hour, _ := strconv.Atoi(timeParts[0])
	if hour >= 14 && hour < 16 {
		points += 10
	}

	return points
}

// isAlphanumeric checks if a character is alphanumeric.
func isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

// rootHandler displays a welcome message for the root URL.
func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Welcome to the Receipt Processor API! Use /receipts/process to submit a receipt and /receipts/{id}/points to get points."))
}

// main initializes the server and registers the endpoints.
func main() {
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/receipts/process", processReceipt)
	http.HandleFunc("/receipts/", getPoints)
	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
