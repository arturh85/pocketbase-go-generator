package pocketbase_go_generator

import (
	"github.com/arturh85/pocketbase-go-generator/internal/cmd"
	"github.com/pocketbase/pocketbase"
	pbcore "github.com/pocketbase/pocketbase/core"
)

type GeneratorOptions struct {
	AllCollections     bool
	CollectionsInclude []string
	CollectionsExclude []string

	Output string
}

func RegisterHook(app *pocketbase.PocketBase, options *GeneratorOptions) {
	generatorFlags := &cmd.GeneratorFlags{
		AllCollections:     options.AllCollections,
		CollectionsInclude: options.CollectionsInclude,
		CollectionsExclude: options.CollectionsExclude,

		Output: options.Output,
	}

	app.OnCollectionAfterCreateSuccess().BindFunc(func(e *pbcore.CollectionEvent) error {
		_ = processFileGeneration(app, generatorFlags)

		return e.Next()
	})

	app.OnCollectionAfterUpdateSuccess().BindFunc(func(e *pbcore.CollectionEvent) error {
		_ = processFileGeneration(app, generatorFlags)

		return e.Next()
	})

	app.OnCollectionAfterDeleteSuccess().BindFunc(func(e *pbcore.CollectionEvent) error {
		_ = processFileGeneration(app, generatorFlags)

		return e.Next()
	})
}
