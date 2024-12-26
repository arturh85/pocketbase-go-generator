package pocketbase_ts_generator

import (
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"pocketbase-ts-generator/internal/cmd"
)

func RegisterCommand(app *pocketbase.PocketBase) {
	app.RootCmd.AddCommand(cmd.GetGenerateTsCommand(true, func(cmd *cobra.Command, args []string, generatorFlags *cmd.GeneratorFlags) {
		err := processFileGeneration(app, generatorFlags)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not process file generation")
		}
	}))
}
