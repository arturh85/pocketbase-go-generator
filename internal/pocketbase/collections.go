package pocketbase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type CollectionField struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type Collection struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
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
