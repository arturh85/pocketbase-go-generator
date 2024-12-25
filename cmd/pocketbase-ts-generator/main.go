package main

import (
	"fmt"
	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"pocketbase-ts-generator/internal/pocketbase"
)

var (
	hostname string
	email    string
	password string
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Hostname").
				Value(&hostname),
			huh.NewInput().
				Title("Email address").
				Value(&email),
			huh.NewInput().
				Title("Password").
				Value(&password).
				EchoMode(huh.EchoModePassword),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Form error")
	}

	pocketBase := pocketbase.New(&pocketbase.PocketBaseCredentials{
		Host:     hostname,
		Email:    email,
		Password: password,
	})

	err = pocketBase.Authenticate()
	if err != nil {
		log.Fatal().Err(err).Msg("Authentication error")
	}

	collections, err := pocketBase.GetCollections()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not retrieve collections")
	}

	fmt.Println(collections)
}
