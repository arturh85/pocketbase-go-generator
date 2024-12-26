package pocketbase_core

import (
	"github.com/pocketbase/pocketbase"
	"pocketbase-ts-generator/internal/pocketbase_api"
)

func GetCollections(app *pocketbase.PocketBase) (*pocketbase_api.CollectionsResponse, error) {
	pbCollections, err := app.App.FindAllCollections()
	if err != nil {
		return nil, err
	}

	output := &pocketbase_api.CollectionsResponse{
		Items: convertPBCollections(pbCollections),
	}

	return output, nil
}
