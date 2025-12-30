package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	passwords := map[string]string{
		"Admin@123":  "admin@oficinapro.com",
		"Client@123": "client@oficinapro.com",
	}

	for pass, user := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		fmt.Printf("-- %s (%s)\n", user, pass)
		fmt.Printf("'%s',\n\n", string(hash))
	}
}
