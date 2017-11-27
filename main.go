package main

import (
	"log"
	"os"
	"utsusemi"
)

func main() {
	logger := log.New(os.Stdout, "[utsusemi] ", log.Ldate|log.Ltime)

	flags := utsusemi.ParseFlag()
	config, err := utsusemi.LoadConfig(flags)

	if err != nil {
		logger.Fatalf("Load config failed: %s", err)
	}

	server, err := utsusemi.NewServer(config, logger)

	if err != nil {
		logger.Fatalf("Create server failed: %s", err)
	}

	err = server.Run()

	if err != nil {
		logger.Fatal(err)
	}
}
