package main

import (
	"github.com/arturh85/pocketbase-go-generator/pkg/pocketbase-go-generator"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
)

func main() {
	app := pocketbase.New()

	pocketbase_go_generator.RegisterCommand(app)

	if err := app.Start(); err != nil {
		log.Fatal().Err(err)
	}
}
