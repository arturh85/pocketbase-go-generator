package forms

import (
	"errors"
	"github.com/Vogeslu/pocketbase-ts-generator/internal/credentials"
	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"
)

func AskCredentials(pbCredentials *credentials.Credentials) bool {
	credentialExist, isEncrypted, err := credentials.CheckExistingCredentials()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not check for credentials")
	}

	if credentialExist {
		if isEncrypted {
			var encryptionPassword string

			credentialsForm := huh.NewForm(
				huh.NewGroup(
					huh.NewInput().
						Title("Encryption password").
						Description("Used to decrypt the stored credentials.env file. Delete the file or enter nothing to enter new credentials.").
						Value(&encryptionPassword).
						EchoMode(huh.EchoModePassword),
				),
			)

			err := credentialsForm.Run()
			if err != nil {
				log.Fatal().Err(err).Msg("Credentials form error")
			}

			if encryptionPassword != "" {
				err = pbCredentials.Decrypt(encryptionPassword)
				if err != nil {
					log.Fatal().Err(err).Msg("Could not decrypt stored credentials")
				}

				return false
			}
		} else {
			err = pbCredentials.Load()
			if err != nil {
				log.Fatal().Err(err).Msg("Could not load stored credentials")
			}

			return false
		}
	}

	var storeCredentials bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Hostname").
				Value(&pbCredentials.Host),
			huh.NewInput().
				Title("Email address").
				Value(&pbCredentials.Email),
			huh.NewInput().
				Title("Password").
				Value(&pbCredentials.Password).
				EchoMode(huh.EchoModePassword),
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to store the credentials?").
				Value(&storeCredentials),
		),
	)

	err = form.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Form error")
	}

	return storeCredentials
}

func AskStoreCredentials(pbCredentials *credentials.Credentials) {
	var encryptCredentials bool

	useEncryptionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title("Do you want to encrypt the credentials?").
				Value(&encryptCredentials),
		),
	)

	err := useEncryptionForm.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Use encryption form error")
	}

	if encryptCredentials {
		var encryptionPassword string
		var encryptionPasswordRepeat string

		credentialsForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Encryption password").
					Value(&encryptionPassword).
					EchoMode(huh.EchoModePassword).
					Validate(func(str string) error {
						if str == "" {
							return errors.New("password cannot be empty")
						}

						return nil
					}),
				huh.NewInput().
					Title("Repeat encryption password").
					Value(&encryptionPasswordRepeat).
					EchoMode(huh.EchoModePassword).
					Validate(func(str string) error {
						if str != encryptionPassword {
							return errors.New("passwords do not match")
						}

						return nil
					}),
			),
		)

		err = credentialsForm.Run()
		if err != nil {
			log.Fatal().Err(err).Msg("Form error")
		}

		err = pbCredentials.Encrypt(encryptionPassword)
		if err != nil {
			log.Fatal().Err(err).Msg("Encrypt error")
		}
	} else {
		err = pbCredentials.Save()
		if err != nil {
			log.Fatal().Err(err).Msg("Save error")
		}
	}
}
