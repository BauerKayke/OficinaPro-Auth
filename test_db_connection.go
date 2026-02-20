package main

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Test 1: Raw DSN from environment
	fmt.Println("=== TEST 1: DSN from ENV ===")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSL_MODE")

	fmt.Printf("Host: %s\n", host)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("User: %s\n", user)
	fmt.Printf("DB: %s\n", database)
	fmt.Printf("SSL: %s\n", sslmode)
	fmt.Printf("Password length: %d\n", len(password))
	fmt.Printf("Password (first 5 chars): %s...\n", password[:5])

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, database, sslmode,
	)

	fmt.Printf("\nDSN (censored): host=%s port=%s user=%s password=*** dbname=%s sslmode=%s\n",
		host, port, user, database, sslmode)

	fmt.Println("\n=== Attempting connection ===")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf("Failed to ping: %v", err)
	}

	fmt.Println("✅ Connection successful!")

	// Test query
	var result string
	if err := db.Raw("SELECT current_database()").Scan(&result).Error; err != nil {
		log.Fatalf("Failed to query: %v", err)
	}

	fmt.Printf("✅ Current database: %s\n", result)
	sqlDB.Close()
}
