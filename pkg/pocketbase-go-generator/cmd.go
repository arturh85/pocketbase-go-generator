package pocketbase_go_generator

import (
	"github.com/arturh85/pocketbase-go-generator/internal/cmd"
	"github.com/pocketbase/pocketbase"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func RegisterCommand(app *pocketbase.PocketBase) {
	app.RootCmd.AddCommand(cmd.GetGenerateGoCommand(true, func(cmd *cobra.Command, args []string, generatorFlags *cmd.GeneratorFlags) {
		err := processFileGeneration(app, generatorFlags)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not process file generation")
		}
	}))
}
