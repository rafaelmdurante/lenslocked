package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
)

func main() {
	// os.Args[0] is the file location
	switch os.Args[1] {
	case "hash":
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Printf("invalid command: %v\nuse 'hash' or 'compare'", os.Args[1])
	}
}

// sample: go run cmd/bcrypt/bcrypt.go hash "secret password"
func hash(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword(
		[]byte(password), bcrypt.DefaultCost)

	if err != nil {
		fmt.Printf("error hashing: %v\n", err)
		return
	}

	hash := string(hashedBytes)
	fmt.Println(hash)
}

// sample: go run cmd/bcrypt/bcrypt.go compare "secret password" '$2a$10$9gHK94mwD7KmxflYvnzjYuWRfRGx9/FtE7rB2.V9TaraRw4r40tGS'
func compare(password, hash string) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

	if err != nil {
		fmt.Printf("password is invalid: %v\n", password)
		return
	}

	fmt.Println("password is correct!")
}
