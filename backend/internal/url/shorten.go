package url

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	database "url-shortener/backend/internal/database"
	"url-shortener/backend/internal/models"
	qrcode "url-shortener/backend/internal/qrcode"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
	QRCode   string `json:"qr_code"`
}

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length  = 4
)

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func generateShortID() string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ShortenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.URL == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	fmt.Printf("Received request to shorten URL: %s\n", req.URL)

	// Check if the URL already exists in the database
	var existingURL models.ShortenedURL
	result := database.DB.Where("long_url = ?", req.URL).First(&existingURL)
	if result.Error == nil {
		// URL already exists, return the existing short URL and QR code
		qrCodeBase64 := base64.StdEncoding.EncodeToString(existingURL.QRCode)
		fmt.Printf("QR Code: %s\n", qrCodeBase64)
		resp := ShortenResponse{
			ShortURL: existingURL.ShortURL,
			QRCode:   qrCodeBase64,
		}

		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(resp)
		if err != nil {
			fmt.Printf("Error encoding response: %v\n", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		fmt.Printf("Returning existing short URL: %s with QR code\n", existingURL.ShortURL)
		return
	}

	// Generate a new short ID
	var shortID string
	for {
		shortID = generateShortID()
		var existingShortURL models.ShortenedURL
		result := database.DB.Where("id = ?", shortID).First(&existingShortURL)
		if result.Error != nil {
			// Short ID doesn't exist and can be used
			break
		}
	}

	shortURL := fmt.Sprintf("http://localhost:3000/%s", shortID)

	// Generate QR code
	qrCodeFileName := fmt.Sprintf("%s.png", shortID)
	qrCodePath, err := qrcode.GenerateQRCode(shortURL, qrCodeFileName)
	if err != nil {
		fmt.Printf("Error generating QR code: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Read QR code file
	qrCodeData, err := os.ReadFile(qrCodePath)
	if err != nil {
		fmt.Printf("Error reading QR code file: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	qrCodeBase64 := base64.StdEncoding.EncodeToString(qrCodeData)
	newURL := models.ShortenedURL{
		ID:       shortID,
		LongURL:  req.URL,
		ShortURL: shortURL,
		QRCode:   qrCodeData,
	}

	result = database.DB.Create(&newURL)
	if result.Error != nil {
		fmt.Printf("Error saving to database: %v\n", result.Error)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = os.Remove(qrCodePath)
	if err != nil {
		return
	}

	resp := ShortenResponse{
		ShortURL: shortURL,
		QRCode:   qrCodeBase64,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		fmt.Printf("Error encoding response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Printf("Created and returned new short URL: %s with QR code\n", shortURL)
}

func RedirectToURL(w http.ResponseWriter, r *http.Request) {
	shortID := r.URL.Path[1:]
	log.Printf("Received redirect request for shortID: %s\n", shortID)

	if shortID == "" {
		log.Println("Error: Empty shortID")
		http.Error(w, "Invalid short URL", http.StatusBadRequest)
		return
	}

	var shortenedURL models.ShortenedURL
	result := database.DB.First(&shortenedURL, "id = ?", shortID)
	if result.Error != nil {
		log.Printf("Error finding URL in database: %v\n", result.Error)
		http.Error(w, "URL not found", http.StatusNotFound)
		return
	}

	shortenedURL.ClickCount++
	result = database.DB.Save(&shortenedURL)
	if result.Error != nil {
		log.Printf("Error updating click count: %v\n", result.Error)
	}

	log.Printf("Found long URL: %s, Click count: %d\n", shortenedURL.LongURL, shortenedURL.ClickCount)

	qrCodeBase64 := base64.StdEncoding.EncodeToString(shortenedURL.QRCode)
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{
		"long_url": shortenedURL.LongURL,
		"qr_code":  qrCodeBase64,
	})

	if err != nil {
		log.Printf("Error encoding JSON: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Println("Successfully sent JSON response")
}
