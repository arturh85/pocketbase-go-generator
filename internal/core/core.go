package core

import (
	"fmt"
	"os"
	"strings"

	"github.com/arturh85/pocketbase-go-generator/internal/cmd"
	"github.com/arturh85/pocketbase-go-generator/internal/interpreter"
	"github.com/arturh85/pocketbase-go-generator/internal/pocketbase_api"
	"github.com/rs/zerolog/log"
)

func ProcessCollections(selectedCollections []*pocketbase_api.Collection, allCollections []pocketbase_api.Collection, generatorFlags *cmd.GeneratorFlags) {
	interpretedCollections := interpreter.InterpretCollections(selectedCollections, allCollections)

	output := make([]string, len(interpretedCollections)+1)

	for i, collection := range interpretedCollections {
		output[i] = collection.GetGoStruct(generatorFlags) + "\n" + collection.GetGoRecord(generatorFlags)
	}

	collectionDefinitions := make([]string, len(interpretedCollections))
	for i, collection := range interpretedCollections {
		collectionDefinitions[i] = collection.GetGoCollectionEntry(generatorFlags)
	}

	output[len(interpretedCollections)] = fmt.Sprintf("const (\n%s\n)", strings.Join(collectionDefinitions, "\n"))

	imports := `
import (
    "github.com/pocketbase/pocketbase/core"
    "github.com/pocketbase/pocketbase/tools/types"
)
`

	joinedData := "package collections\n\n" + imports + "\n" + strings.Join(output, "\n\n")

	helper_funcs := make([]string, len(interpretedCollections))
	for i, collection := range interpretedCollections {
		helper_funcs[i] = collection.GetGoCollectionHelperFuncs(generatorFlags)
	}
	joinedData += strings.Join(helper_funcs, "\n\n")

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
