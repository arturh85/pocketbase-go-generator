package pocketbase_core

import (
	"github.com/arturh85/pocketbase-go-generator/internal/pocketbase_api"
	"github.com/pocketbase/pocketbase"
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
