package main

import (
    "fmt"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    passwords := []string{
        "alice123",
        "bob123",
        "admin123",
    }

    for _, p := range passwords {
        hash, _ := bcrypt.GenerateFromPassword(
            []byte(p),
            bcrypt.DefaultCost,
        )

        fmt.Println(string(hash))
    }
}