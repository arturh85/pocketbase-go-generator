package pocketbase_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/arturh85/pocketbase-go-generator/internal/credentials"
	"github.com/rs/zerolog/log"
)

type PocketBase struct {
	credentials *credentials.Credentials
	token       string
	client      *http.Client
}

func New(Credentials *credentials.Credentials) *PocketBase {
	pocketBase := &PocketBase{
		credentials: Credentials,
		client:      &http.Client{},
	}

	return pocketBase
}

func (pocketBase *PocketBase) GetApiUrl(suffix string) string {
	return fmt.Sprintf("%s/api/%s", pocketBase.credentials.Host, suffix)
}

type pocketBaseAuthResponse struct {
	Token string `json:"token"`
}

func (pocketBase *PocketBase) Authenticate() error {
	log.Info().Msgf("Authenticating with %s...", pocketBase.credentials.Host)

	body := []byte(fmt.Sprintf(`{
		"identity": "%s",
		"password": "%s"
	}`, pocketBase.credentials.Email, pocketBase.credentials.Password))

	request, err := http.NewRequest("POST", pocketBase.GetApiUrl("collections/_superusers/auth-with-password"), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	request.Header.Add("Content-Type", "application/json")

	response, err := pocketBase.client.Do(request)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	log.Debug().Msgf("Got status code %d", response.StatusCode)

	if response.StatusCode != http.StatusOK {
		return errors.New("invalid status code, expected 200")
	}

	authResponse := &pocketBaseAuthResponse{}
	err = json.NewDecoder(response.Body).Decode(authResponse)
	if err != nil {
		return err
	}

	if authResponse.Token == "" {
		return errors.New("token is missing")
	}

	log.Debug().Msgf("Got token %s", authResponse.Token)

	log.Info().Msgf("Authentication successful")

	pocketBase.token = authResponse.Token

	return nil
}

func (pocketBase *PocketBase) DoWithAuth(request *http.Request) (*http.Response, error) {
	if pocketBase.token != "" {
		request.Header.Set("Authorization", pocketBase.token)
	}

	return pocketBase.client.Do(request)
}
