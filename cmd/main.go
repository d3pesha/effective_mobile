package main

import (
	"em/internal/config"
	"em/internal/infrastructure"
	"time"
)

func main() {
	cfg, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}

	app := infrastructure.NewConfig(cfg).
		Logger().
		Database().
		ContextTimeout(10 * time.Second)

	app.WebServer().Start()
}
