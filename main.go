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

	// Remove for Docker Image
	err := LoadDotEnv(); if err != nil{
		slog.Error("program encountered a fatal error", slog.Any("error", err))
		panic(err)
	}

	// Get file path for reading services
	FILE_PATH := os.Getenv("FILE_PATH")

	// Read all services
	services := modules.ReadServices(FILE_PATH)
	
	// Probe all services
	services_resp := modules.ProbeAllServices(services)

	// Print
	fmt.Println(services_resp.ToStr())

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