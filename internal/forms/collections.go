package forms

import (
	"github.com/Vogeslu/pocketbase-ts-generator/internal/cmd"
	"github.com/Vogeslu/pocketbase-ts-generator/internal/pocketbase_api"
	"github.com/charmbracelet/huh"
	"github.com/rs/zerolog/log"
	"sort"
	"strings"
)

func AskCollectionSelection(collections []pocketbase_api.Collection) []*pocketbase_api.Collection {
	options := make([]huh.Option[*pocketbase_api.Collection], len(collections))

	for i, collection := range collections {
		options[i] = huh.NewOption(collection.String(), &collection).Selected(!collection.System)
	}

	sort.SliceStable(options, func(i, j int) bool {
		return options[i].Value.System != options[j].Value.System
	})

	var output []*pocketbase_api.Collection

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[*pocketbase_api.Collection]().
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

func GetSelectedCollections(generatorFlags *cmd.GeneratorFlags, collections []pocketbase_api.Collection) []*pocketbase_api.Collection {
	var output []*pocketbase_api.Collection

	checkInclude := len(generatorFlags.CollectionsInclude) > 0

	for _, collection := range collections {
		if checkInclude && !containsIgnoreCase(generatorFlags.CollectionsInclude, collection.Name) {
			continue
		} else if containsIgnoreCase(generatorFlags.CollectionsExclude, collection.Name) {
			continue
		} else if !generatorFlags.AllCollections && collection.System && !(checkInclude && containsIgnoreCase(generatorFlags.CollectionsInclude, collection.Name)) {
			continue
		}

		output = append(output, &collection)
	}

	return output
}

func containsIgnoreCase(list []string, value string) bool {
	for _, v := range list {
		if strings.ToLower(v) == strings.ToLower(value) {
			return true
		}
	}

	return false
}
