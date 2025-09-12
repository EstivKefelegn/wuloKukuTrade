package main

import (
	"chickenTrade/API/internal/api/router"
	"chickenTrade/API/internal/repository/sqlconnect"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	db, err := sqlconnect.ConnectDB("trade_chicken")
	if err != nil {
		fmt.Println("Error : ----")
		return
	}
	_ = db

	cert := os.Getenv("CERT_FILE")
	key := os.Getenv("KEY_FILE")

	port :=  os.Getenv("API_PORT")
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}

	router := router.MainRouter()

	fmt.Println("Server is going to start")
	server := &http.Server{
		Addr: port,
		Handler: router,
		TLSConfig: tlsConfig,
	}

	fmt.Println("Server running on port 3000")
	err = server.ListenAndServeTLS(cert, key)
	if err != nil {
		log.Println("Tls server failed to connect ", err)
	}
}
