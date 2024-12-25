package forms

import (
	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"
	"pocketbase-ts-generator/internal/pocketbase"
	"sort"
)

func AskCollectionSelection(collections []pocketbase.Collection) []*pocketbase.Collection {
	options := make([]huh.Option[*pocketbase.Collection], len(collections))

	for i, collection := range collections {
		options[i] = huh.NewOption(collection.String(), &collection).Selected(!collection.System)
	}

	sort.SliceStable(options, func(i, j int) bool {
		return options[i].Value.System != options[j].Value.System
	})

	var output []*pocketbase.Collection

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[*pocketbase.Collection]().
				Options(options...).
				Title("Select collections to generate types from").
				Value(&output),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal().Err(err).Msg("Form error")
	}

	return output
}
