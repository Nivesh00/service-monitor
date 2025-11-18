package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"

	"github.com/Nivesh00/service-monitor/modules"
)

func main() {

	err := godotenv.Load(".env.public")
	if err != nil {
		slog.Error("cannot load .env.public file")
	}

	err = godotenv.Load(".env.private")
	if err != nil {
		slog.Error("cannot load .env.private file")
	}

	FILE_PATH := os.Getenv("FILE_PATH")

	services := modules.ReadServices(FILE_PATH)

	fmt.Println(services.ToStr())
	fmt.Println(services.Services[0].ToStr())
}