package main

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
	"pocketbase-ts-generator/internal/cmd"
	"pocketbase-ts-generator/internal/core"
	"pocketbase-ts-generator/internal/credentials"
	"pocketbase-ts-generator/internal/forms"
	"pocketbase-ts-generator/internal/pocketbase_api"
)

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	rootCmd := cmd.GetGenerateTsCommand(false, func(cmd *cobra.Command, args []string, generatorFlags *cmd.GeneratorFlags) {
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

		pocketBase := pocketbase_api.New(pbCredentials)

		err := pocketBase.Authenticate()
		if err != nil {
			log.Fatal().Err(err).Msg("Authentication error")
		}

		collections, err := pocketBase.GetCollections()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not retrieve collections")
		}

		var selectedCollections []*pocketbase_api.Collection
		outputTarget := generatorFlags.Output

		if !generatorFlags.DisableForm {
			selectedCollections = forms.AskCollectionSelection(collections.Items)
			outputTarget = forms.AskOutputTarget(outputTarget)
		} else {
			selectedCollections = forms.GetSelectedCollections(generatorFlags, collections.Items)
		}

		core.ProcessCollections(selectedCollections, collections.Items, outputTarget)
	})

	err := rootCmd.Execute()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed processing command")
	}
}
