package main

import (
	"bitmex-gorilla/bitmex"
	"bitmex-gorilla/server"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"
	"github.com/gin-gonic/gin"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	endpoint, ok := os.LookupEnv("API_ENDPOINT")

	if !ok {
		log.Fatal("No wss endpoint specified")
	}

	client := bitmex.NewClient(endpoint)

	go client.Start()

	broadcaster := server.NewBroadcaster(client.DataStream)

	go broadcaster.Start()

	//gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()
	engine.GET("/ws", broadcaster.NewClient)

	engine.StaticFile("/test","console_test.html")
	go engine.Run("127.0.0.1:8081")

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			client.Stop()
			return
		}
	}
}
