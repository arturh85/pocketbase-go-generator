package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"pocketbase-ts-generator/internal/credentials"
	"pocketbase-ts-generator/internal/forms"
	"pocketbase-ts-generator/internal/interpreter"
	"pocketbase-ts-generator/internal/pocketbase"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	pbCredentials := &credentials.Credentials{}

	storeCredentials := forms.AskCredentials(pbCredentials)

	if storeCredentials {
		forms.AskStoreCredentials(pbCredentials)
	}

	pocketBase := pocketbase.New(pbCredentials)

	err := pocketBase.Authenticate()
	if err != nil {
		log.Fatal().Err(err).Msg("Authentication error")
	}

	collections, err := pocketBase.GetCollections()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not retrieve collections")
	}

	selectedCollections := forms.AskCollectionSelection(collections.Items)

	interpretedCollections := interpreter.InterpretCollections(selectedCollections, collections.Items)

	for _, collection := range interpretedCollections {
		fmt.Println(collection.GetTypescriptInterface())
	}
}
