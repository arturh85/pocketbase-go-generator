package pocketbase_ts_generator

import (
	"github.com/pocketbase/pocketbase"
	"pocketbase-ts-generator/internal/cmd"
	"pocketbase-ts-generator/internal/core"
	"pocketbase-ts-generator/internal/forms"
	"pocketbase-ts-generator/internal/pocketbase_api"
	"pocketbase-ts-generator/internal/pocketbase_core"
)

func processFileGeneration(app *pocketbase.PocketBase, generatorFlags *cmd.GeneratorFlags) error {
	collections, err := pocketbase_core.GetCollections(app)
	if err != nil {
		return err
	}

	var selectedCollections []*pocketbase_api.Collection
	outputTarget := generatorFlags.Output

	selectedCollections = forms.GetSelectedCollections(generatorFlags, collections.Items)

	core.ProcessCollections(selectedCollections, collections.Items, outputTarget)

	return nil
}
