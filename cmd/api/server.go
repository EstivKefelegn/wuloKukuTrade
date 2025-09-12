package main

import (
	"chickenTrade/API/internal/api/router"
	"crypto/tls"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//go:embed .env
var envfile embed.FS

func loadEnvFromEmbededFile() {
	content, err := envfile.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	tempFile, err := os.CreateTemp("", ".env")
	if err != nil {
		log.Fatalf("Error creating temp .env file %v:", err)
	}

	defer os.Remove(tempFile.Name())

	_, err = tempFile.Write(content)
	if err != nil {
		log.Fatalf("Error closing tempfile %v:", err)
	}
	err = tempFile.Close()
	if err != nil {
		log.Fatalf("Error closing tempfile %v:", err)
	}

	err = godotenv.Load(tempFile.Name())
	if err != nil {
		log.Fatalf("Error loading .env file %v:", err)
	}
}

func main() {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Fatalf("Error loading .env file %v:", err)
	// }

	loadEnvFromEmbededFile()

	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	port := os.Getenv("API_PORT")
	fmt.Println("PORT:", port)
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	router := router.MainRouter()

	fmt.Println("Server is going to start")
	server := &http.Server{
		Addr:      port,
		Handler:   router,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server running on port 3000")
	err := server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Println("Tls server failed to connect ", err)
	}
}
