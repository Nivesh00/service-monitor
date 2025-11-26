package main

import (
	"errors"
	// "fmt"
	"github.com/joho/godotenv"
	"log/slog"
	"os"
	"time"

	"github.com/Nivesh00/service-monitor/modules"
)

func main() {

	// start := time.Now()

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
	for {
		go func() {
			slog.Info("Starting probe")
			time.Sleep(5 * time.Second)
			modules.ProbeAllServices(services)
		}()

		go func() {
			slog.Info("Other task")
			time.Sleep(10 * time.Second)
		}()
	}

	// fmt.Printf("Program duration: %6s\n", time.Since(start))
	// os.Exit(0)
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