package main

import (
	"fmt"
	"log"
	"passwordStorage/database"
	"passwordStorage/handlers"
)

func main() {
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
