// File: backend/cmd/tools/hashpwd.go
package main

import (
    "fmt"
    "os"
    "golang.org/x/crypto/bcrypt"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Please provide a password as an argument.")
        os.Exit(1)
    }
    
    password := os.Args[1]
    
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        fmt.Println("Error generating hash:", err)
        os.Exit(1)
    }
    
    // Print the hash. IMPORTANT: Bcrypt hashes can contain '$' which can be
    // misinterpreted by shells. Putting it in quotes is a good idea.
    fmt.Printf("Password: %s\n", password)
    fmt.Printf("Bcrypt Hash: '%s'\n", string(hash))
}