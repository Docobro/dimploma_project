package main

import (
	"log"

	"github.com/docobro/dimploma_project/internal/app"
	"github.com/docobro/dimploma_project/internal/db/queries"
)

func main() {
	a := app.NewApp()
	err := a.Run()
	if err != nil {
		log.Fatalf("failed to run app:%v", err)
	}
	res, err := queries.User.FilterWithNameAndRole("john", "admim")
	if err != nil {
		log.Println("error:", err.Error())
	}
	log.Println(res)
	log.Println("program has started")
}
