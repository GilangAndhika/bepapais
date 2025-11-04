package main
import (
	"fmt"
	"log"
	"golang.org/x/crypto/bcrypt"
)
func main() {
	// GANTI "password-rahasia-anda" DENGAN PASSWORD YANG ANDA INGINKAN
	pass := "admin123" 
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Password:", pass)
	fmt.Println("Hash:", string(hash))
}