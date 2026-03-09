package main

import (
	"log"

	"github.com/CedricThomas/console/internal/config"
)

func main() {
	_, err := config.New()
	if err != nil {
		log.Fatalf("Cannot initialize configuration: %v", err)
	}

	log.Println("PC-Agent server started successfully")
}
