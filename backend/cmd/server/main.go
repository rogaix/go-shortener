package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
	database "url-shortener/backend/internal/database"
	"url-shortener/backend/internal/models"
	url "url-shortener/backend/internal/url"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	err := database.DB.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci").AutoMigrate(&models.ShortenedURL{})
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}

	// Auto migrate the ShortenedURL model
	if err := database.DB.AutoMigrate(&models.ShortenedURL{}); err != nil {
		log.Fatalf("Failed to auto migrate ShortenedURL model: %v", err)
	}

	http.HandleFunc("/api/shorten", corsMiddleware(url.ShortenURL))
	http.HandleFunc("/", corsMiddleware(url.RedirectToURL))

	port := 8080
	var listener net.Listener
	for {
		listener, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err != nil {
			if strings.Contains(err.Error(), "address already in use") {
				port++
				continue
			}
			log.Fatalf("Failed to bind to port %d: %v", port, err)
		}
		break
	}

	fmt.Printf("Server is running on http://localhost:%d\n", port)

	if err := http.Serve(listener, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}