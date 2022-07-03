package main

import (
	app "spender/v1/api/app"
	config "spender/v1/api/config"
)

func main() {
	config := config.GetConfig()

	app := &app.App{}
	app.Initialize(config)
	app.Run(":3000")
}
