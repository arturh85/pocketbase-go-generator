package main

import (
	"github.com/Vogeslu/pocketbase-ts-generator/pkg/pocketbase-ts-generator"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
)

func main() {
	app := pocketbase.New()

	pocketbase_ts_generator.RegisterCommand(app)

	if err := app.Start(); err != nil {
		log.Fatal().Err(err)
	}
}
