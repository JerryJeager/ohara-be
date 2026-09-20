package main

import (
	"log"

	"github.com/JerryJeager/ohara-be/cmd"
	"github.com/JerryJeager/ohara-be/config"
)

func init() {
	config.LoadEnv()
	config.ConnectToDB()
	config.NewAIClient()
}

func main() {
	log.Println("Starting Server")

	cmd.ExecuteApiRoutes()
}
