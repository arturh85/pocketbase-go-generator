package pocketbase_api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type CollectionField struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	CollectionId string   `json:"collectionId"`
	MaxSelect    int      `json:"maxSelect"`
	Required     bool     `json:"required"`
	Hidden       bool     `json:"hidden"`
	Values       []string `json:"values"`
}

type Collection struct {
	Id     string            `json:"id"`
	Name   string            `json:"name"`
	Type   string            `json:"type"`
	System bool              `json:"system"`
	Fields []CollectionField `json:"fields"`
}

type CollectionsResponse struct {
	Items []Collection `json:"items"`
}

func (pocketBase *PocketBase) GetCollections() (*CollectionsResponse, error) {
	request, err := http.NewRequest("GET", pocketBase.GetApiUrl("collections?perPage=500"), nil)
	if err != nil {
		return nil, err
	}

	response, err := pocketBase.DoWithAuth(request)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status code, expected 200")
	}

	collectionResponse := &CollectionsResponse{}
	err = json.NewDecoder(response.Body).Decode(collectionResponse)
	if err != nil {
		return nil, err
	}

	return collectionResponse, nil
}

func (collection Collection) String() string {
	if collection.System {
		return fmt.Sprintf("%s (System, %d fields)", collection.Name, len(collection.Fields))
	}

	return fmt.Sprintf("%s (%d fields)", collection.Name, len(collection.Fields))
}
