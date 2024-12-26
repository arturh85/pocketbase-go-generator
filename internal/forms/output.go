package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"
)

func AskOutputTarget(inputValue string) string {
	var outputTarget string = inputValue

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Output target").
				Description("Target file for generated interfaces, keep empty to print results directly in console").
				Value(&outputTarget),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Output target form error")
	}

	return outputTarget
}
