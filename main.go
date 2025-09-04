package main

import (
	"fmt"
	"log"
	"passwordStorage/database"
	"passwordStorage/handlers"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found or failed to load")
	}

	handlers.InitOAuthConfigs()

	db, err := database.DbInit()
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}

	defer db.Close()

	router := handlers.RoutesHandler(db)

	err = router.Run(":2137")
	fmt.Println("Listning on port :2137...")
	if err != nil {
		log.Fatalf("Error while running server: %v", err)
	}

}
