package main

import (
	"fmt"
	"syscall"

	"os"

	"log"
	"myproject/pkg/config"
	db "myproject/pkg/database"
	"myproject/pkg/di"
	"os/signal"
)

func main() {
	// Load configuration
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", cwd)
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	errs := db.InitRedis()
	if errs != nil {
		fmt.Printf("Error initializing Redis: %s\n", err.Error())
	} else {
		fmt.Println("Redis connection successful!")
	}

	server, err := di.InitializeEvent(conf)
	if err != nil {
		log.Fatal("failed to initialize the files")
	}

	server.Start(conf)
	c := make(chan os.Signal, 1) // Create channel to signify a signal being sent
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	_ = <-c
	_ = server.Kill()

}
