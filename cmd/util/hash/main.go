package main

import (
	"fmt"
	"nds-go-starter/internal/core/auth"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/util/hash/main.go <password>")
		os.Exit(1)
	}

	password := os.Args[1]
	hash, err := auth.HashPassword(password)
	if err != nil {
		fmt.Printf("Error hashing password: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Password:", password)
	fmt.Println("Hash:    ", hash)
}
