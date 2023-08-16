package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	switch os.Args[1] {
	case "hash":
		// hash password
		hash(os.Args[2])
	case "compare":
		compare(os.Args[2], os.Args[3])
	default:
		fmt.Printf("Invalid command: %v\n", os.Args[1])
	}
}

func hash(password string) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("error hashing: %v\n", err)
		return
	}
	fmt.Println(string(hashedBytes))
}

func compare(password, hash string) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		fmt.Printf("Invalid Password: %v\n", password)
		return
	}
	fmt.Println("Password Match.")
}
