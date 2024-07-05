package main

import (
	"github.com/joho/godotenv"
	"go-authentication/config"
	"go-authentication/internal/app"
	stdLog "log"
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
