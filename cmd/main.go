package main

import (
	"os"

	"github.com/DavidKlz/homeserver-backend/pkgs/logger"
	"github.com/DavidKlz/homeserver-backend/server"
	"github.com/DavidKlz/homeserver-backend/storage"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Fatal("The .env file could not be loaded: %s", err.Error())
	}

	err = os.MkdirAll(os.Getenv("INPUT_FOLDER"), os.ModePerm)
	err = os.MkdirAll(os.Getenv("OUTPUT_FOLDER"), os.ModePerm)

	if err != nil {
		logger.Fatal("The input and output directory could not be created: %s", os.Getenv("INPUT_FOLDER"), os.Getenv("OUTPUT_FOLDER"))
	}

	store, err := storage.CreatePostgresStore()
	if err != nil {
		logger.Fatal("Error while trying to create postgres store:\n %s", err.Error())
	}

	server := server.CreateNewServer(":3000", store)
	server.Run()
}
