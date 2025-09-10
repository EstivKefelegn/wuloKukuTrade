package main

import (
	"chickenTrade/API/internal/repository/sqlconnect"
	"fmt"
)

func main() {
	db, err := sqlconnect.ConnectDB("trade_chicken")
	if err != nil {
		fmt.Println("Error : ----")
		return
	}

	_ = db

}
