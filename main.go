package main

import (
	"log"

	"github.com/vladjong/taking_survey/client"
	"github.com/vladjong/taking_survey/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	client := client.NewClinet()
	client.Run()
}
