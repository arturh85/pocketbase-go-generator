package main

import (
	pocketbase_go_generator "github.com/arturh85/pocketbase-go-generator/pkg/pocketbase-go-generator"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
)

func main() {
	app := pocketbase.New()

	pocketbase_go_generator.RegisterHook(app, &pocketbase_go_generator.GeneratorOptions{
		Output: "test.go",
	})

	if err := app.Start(); err != nil {
		log.Fatal().Err(err)
	}
}
