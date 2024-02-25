package main

import (
	"github.com/joho/godotenv"
	stdLog "log"
	"sso/config"
	"sso/internal/app"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		stdLog.Fatal("can't set env file:", err)
	}

	stdLog.Println("env are set")
}

func main() {
	cfg := config.MustLoad()

	app.Run(cfg)
}
