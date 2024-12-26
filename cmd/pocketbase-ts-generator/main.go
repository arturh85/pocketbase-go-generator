package main

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"pocketbase-ts-generator/internal/cli"
	"pocketbase-ts-generator/internal/credentials"
	"pocketbase-ts-generator/internal/forms"
	"pocketbase-ts-generator/internal/interpreter"
	"pocketbase-ts-generator/internal/pocketbase"
	"strings"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	generatorFlags, _, err := cli.ParseArgs()
	if err != nil {
		return
	}

	if generatorFlags.DisableLogs {
		zerolog.SetGlobalLevel(4)
	} else {
		zerolog.SetGlobalLevel(1)
	}

	pbCredentials := &credentials.Credentials{
		Host:     generatorFlags.Host,
		Email:    generatorFlags.Email,
		Password: generatorFlags.Password,
	}

	if !generatorFlags.DisableForm {
		storeCredentials := forms.AskCredentials(pbCredentials)

		if storeCredentials {
			forms.AskStoreCredentials(pbCredentials)
		}
	} else {
		credentialExist, isEncrypted, err := credentials.CheckExistingCredentials()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not check for credentials")
		}

		if credentialExist {
			if isEncrypted {
				err = pbCredentials.Decrypt(generatorFlags.EncryptionPassword)
				if err != nil {
					log.Fatal().Err(err).Msg("Could not decrypt stored credentials")
				}
			} else {
				err = pbCredentials.Load()
				if err != nil {
					log.Fatal().Err(err).Msg("Could not load stored credentials")
				}
			}
		}
	}

	pocketBase := pocketbase.New(pbCredentials)

	err = pocketBase.Authenticate()
	if err != nil {
		log.Fatal().Err(err).Msg("Authentication error")
	}

	collections, err := pocketBase.GetCollections()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not retrieve collections")
	}

	var selectedCollections []*pocketbase.Collection
	outputTarget := generatorFlags.Output

	if !generatorFlags.DisableForm {
		selectedCollections = forms.AskCollectionSelection(collections.Items)
		outputTarget = forms.AskOutputTarget(outputTarget)
	} else {
		selectedCollections = forms.GetSelectedCollections(generatorFlags, collections.Items)
	}

	interpretedCollections := interpreter.InterpretCollections(selectedCollections, collections.Items)

	output := make([]string, len(interpretedCollections))

	for i, collection := range interpretedCollections {
		output[i] = collection.GetTypescriptInterface()
	}

	joinedData := strings.Join(output, "\n\n")

	if outputTarget == "" {
		fmt.Println(joinedData)
	} else {
		err := os.WriteFile(outputTarget, []byte(joinedData), 0644)
		log.Info().Msgf("Saved generated interfaces to %s", outputTarget)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not output contents")
		}

	}
}
