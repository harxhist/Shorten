package main

import (
	"be/config"
	"be/constant"
	"be/handler"
	"be/logger"
	"be/middleware"
	"net/http"
	"be/storage"
	"fmt"
	"os"
)

func main() {
	// Load the configuration
	rootConfig, err := config.InitialiseAndGetRootConfig()
	if err != nil {
		// log.Fatalf("Unable to load configuration: %v", err)
		fmt.Printf("Unable to load configuration: %v\n", err)
        os.Exit(1)
	}
	config.SetConfig(rootConfig)
	//Initializing logger
	e := logger.InitLogger()
	if(e != nil){
		fmt.Printf("Error initializing logger: %v\n", err)
        os.Exit(1)
	}
	log := logger.Logger
	
	//Connect to database
	if err := storage.StartDB(); err != nil {
        log.Error("Failed to initialize database: ", err)
    }
	//Connect to S3
	if err := storage.S3Init(); err != nil {
        log.Error("Failed to initialize S3: ", err)
    }
	
	http.Handle("/feedback", (middleware.CorsMiddleware(middleware.VerifyToken((http.HandlerFunc(handler.HandleFeedback))))));
	http.HandleFunc("/ws", handler.HandleWebSocket)
	log.Info(`Server running on port`, constant.PORT)
	if err := http.ListenAndServe(constant.PORT, nil); err != nil {
		log.Error("Failed to start server: ", err)
	}
}