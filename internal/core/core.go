package core

import (
	"fmt"
	"github.com/Vogeslu/pocketbase-ts-generator/internal/cmd"
	"github.com/Vogeslu/pocketbase-ts-generator/internal/interpreter"
	"github.com/Vogeslu/pocketbase-ts-generator/internal/pocketbase_api"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

func ProcessCollections(selectedCollections []*pocketbase_api.Collection, allCollections []pocketbase_api.Collection, generatorFlags *cmd.GeneratorFlags) {
	interpretedCollections := interpreter.InterpretCollections(selectedCollections, allCollections)

	output := make([]string, len(interpretedCollections))

	for i, collection := range interpretedCollections {
		output[i] = collection.GetTypescriptInterface(generatorFlags)
	}

	joinedData := strings.Join(output, "\n\n")

	if generatorFlags.Output == "" {
		fmt.Println(joinedData)
	} else {
		err := os.WriteFile(generatorFlags.Output, []byte(joinedData), 0644)
		log.Info().Msgf("Saved generated interfaces to %s", generatorFlags.Output)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not output contents")
		}

	}
}
