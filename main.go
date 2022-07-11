package main

import (
	"os"

	app "spender/v1/app"
	configs "spender/v1/config"
)

func main() {
	config := configs.GetConfig()

	var port = ""

	if configs.ISLOCAL {
		port = "3000"
	} else {
		port = os.Getenv("PORT")
	}

	app := &app.App{}
	app.Initialize(config)
	app.Run(":"+port)
}

