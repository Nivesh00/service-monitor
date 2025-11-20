package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"

	"github.com/Nivesh00/service-monitor/modules"
)

func main() {

	err := LoadDotEnv(); if err != nil{
		slog.Error("Fatal error occured", slog.Any("error", err))
	}

	FILE_PATH := os.Getenv("FILE_PATH")

	services := modules.ReadServices(FILE_PATH)

	fmt.Println(services.ToStr())

	os.Exit(0)
}

// Remove for Docker Image
func LoadDotEnv() error {
	err := godotenv.Load(".env.public")
	if err != nil {
		return errors.New("cannot load .env.public file")
	}

	err = godotenv.Load(".env.private")
	if err != nil {
		return errors.New("cannot load .env.private file")
	}

	return nil
}