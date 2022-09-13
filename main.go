package main

import (
	"github.com/joho/godotenv"
	"github.com/rizface/golang-api-template/database"
	"github.com/rizface/golang-api-template/database/postgresql"
	"github.com/rizface/golang-api-template/system/logger"
	"github.com/rizface/golang-api-template/system/router"
	"github.com/rizface/golang-api-template/system/server"
)

func main() {
	godotenv.Load()

	// Create new logger
	log := logger.CreateErrorLogger()

	// Create new database instance
	db := database.New()

	// Setup database
	db.Setup(log)
	postgres := postgresql.New(db)
	err := postgres.AutoMigrate(&postgresql.Event{})
	if err != nil {
		panic(err)
	}

	// Create new router instance
	chiRouter := router.CreateRouter()

	// Setup controllers
	server.SetupController(chiRouter, postgres)

	// Create & Setup http server
	httpServer := server.CreateHttpServer(chiRouter)

	// Start http server
	err = httpServer.ListenAndServe()

	// Error handling when http server fail to start
	if err == nil {
		log.Error(err.Error())
	}
}
