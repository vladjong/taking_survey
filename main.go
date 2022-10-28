package main

import (
	"log"

	"github.com/vladjong/taking_survey/config"
	"github.com/vladjong/taking_survey/workerpool"
)

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}
	workerpool.StartWorkerpool()
	// client := client.NewClinet(context.Background())
	// client.Run()
}
