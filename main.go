package main

import (
	"ygsx-backend/database"
	"ygsx-backend/router"
)

func main() {
	database.InitDB()
	r := router.SetupRouter()
	r.Run(":8080")
}
