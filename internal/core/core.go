package core

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"pocketbase-ts-generator/internal/interpreter"
	"pocketbase-ts-generator/internal/pocketbase_api"
	"strings"
)

func ProcessCollections(selectedCollections []*pocketbase_api.Collection, allCollections []pocketbase_api.Collection, outputTarget string) {
	interpretedCollections := interpreter.InterpretCollections(selectedCollections, allCollections)

	output := make([]string, len(interpretedCollections))

	for i, collection := range interpretedCollections {
		output[i] = collection.GetTypescriptInterface()
	}

	joinedData := strings.Join(output, "\n\n")

	if outputTarget == "" {
		fmt.Println(joinedData)
	} else {
		err := os.WriteFile(outputTarget, []byte(joinedData), 0644)
		log.Info().Msgf("Saved generated interfaces to %s", outputTarget)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not output contents")
		}

	}
}
