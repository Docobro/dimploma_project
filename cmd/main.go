package main

import (
	"log"

	"github.com/docobro/dimploma_project/internal/app"
	"github.com/hanagantig/gracy"
)

var (
	version    = "dev"
	configPath = "./conf.yaml"
)

func main() {
	_, err := app.New(configPath)
	if err != nil {
		log.Fatal(err)
	}
	err = gracy.Wait()
	if err != nil {
		log.Fatal("failed to gracefully shutdown")
	}
}
