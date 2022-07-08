package main

import (
	"os"

	app "spender/v1/app"
	config "spender/v1/config"
)

func main() {
	config := config.GetConfig()

	port := os.Getenv("PORT") //heroku
	port = "3000"//local

	app := &app.App{}
	app.Initialize(config)
	app.Run(":"+port)
}

