package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"domainhub/internal/config"
	"domainhub/internal/database"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/term"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Config error:", err)
		return
	}

	db, err := database.ConnectDB(cfg)
	if err != nil {
		fmt.Println("Database connection error:", err)
		return
	}
	defer db.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter user email: ")
	email, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Failed to read email:", err)
		return
	}
	email = strings.TrimSpace(email)

	fmt.Print("Enter new password: ")
	passwordBytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		fmt.Println("Failed to read password:", err)
		return
	}

	password := string(passwordBytes)
	if len(password) < 8 {
		fmt.Println("Password must be at least 8 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword(
		passwordBytes,
		bcrypt.DefaultCost,
	)
	if err != nil {
		fmt.Println("Password hashing error:", err)
		return
	}

	result, err := db.Exec(
		`UPDATE users SET password_hash = $1 WHERE email = $2`,
		string(hash),
		email,
	)
	if err != nil {
		fmt.Println("Password update error:", err)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		fmt.Println("Could not verify update:", err)
		return
	}
	if rows == 0 {
		fmt.Println("No user found with that email")
		return
	}

	fmt.Println("Password reset successfully for:", email)
}
