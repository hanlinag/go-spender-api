package main

import (
	"os"

	app "spender/v1/app"
	config "spender/v1/config"
)

func main() {
	config := config.GetConfig()

	port := os.Getenv("PORT")

	app := &app.App{}
	app.Initialize(config)
	app.Run(":"+port)
}
