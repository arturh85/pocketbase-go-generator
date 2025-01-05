package credentials

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/scrypt"
)

var (
	encryptedFileName string = "credentials.enc.env"
	rawFileName       string = "credentials.env"
)

type Credentials struct {
	Host     string
	Email    string
	Password string
}

func CheckExistingCredentials() (bool, bool, error) {
	_, err := os.Stat(encryptedFileName)

	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Stat(rawFileName)

		if errors.Is(err, os.ErrNotExist) {
			return false, false, nil
		} else if err != nil {
			return false, false, err
		}

		return true, false, nil
	} else if err != nil {
		return false, false, err
	}

	return true, true, nil
}

func (credentials *Credentials) Encrypt(encryptionPassword string) error {
	log.Info().Msg("Encrypting data...")

	key, salt, err := deriveKey(encryptionPassword, nil)
	if err != nil {
		return err
	}

	credentialsData := fmt.Sprintf("%s;%s;%s",
		url.QueryEscape(credentials.Host),
		url.QueryEscape(credentials.Email),
		url.QueryEscape(credentials.Password),
	)

	encryptedCredentialsData, err := encryptString(credentialsData, key)
	if err != nil {
		return err
	}

	_ = os.Remove(encryptedFileName)
	_ = os.Remove(rawFileName)

	data := []byte(fmt.Sprintf("SALT=%s\nDATA=%s",
		base64.URLEncoding.EncodeToString(salt),
		base64.URLEncoding.EncodeToString(encryptedCredentialsData),
	))

	err = os.WriteFile(encryptedFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (credentials *Credentials) Save() error {
	log.Info().Msg("Saving data...")

	_ = os.Remove(encryptedFileName)
	_ = os.Remove(rawFileName)

	data := []byte(fmt.Sprintf("HOST=%s\nEMAIL=%s\nPASSWORD=%s",
		credentials.Host,
		credentials.Email,
		credentials.Password,
	))

	err := os.WriteFile(rawFileName, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (credentials *Credentials) Decrypt(encryptionPassword string) error {
	log.Info().Msg("Decrypting data...")

	encryptedData := make(map[string][]byte)

	file, err := os.Open(encryptedFileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Failed closing credentials file")
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) < 2 {
			return errors.New("invalid credentials file content")
		}

		encryptedData[parts[0]], err = base64.URLEncoding.DecodeString(parts[1])
		if err != nil {
			return err
		}
	}

	salt, ok := encryptedData["SALT"]
	if !ok {
		return errors.New("salt is missing")
	}

	encryptedCredentials, ok := encryptedData["DATA"]
	if !ok {
		return errors.New("data is missing")
	}

	key, salt, err := deriveKey(encryptionPassword, salt)
	if err != nil {
		return err
	}

	decryptedCredentials, err := decryptBytes(encryptedCredentials, key)
	if err != nil {
		return err
	}

	splitCredentials := strings.Split(decryptedCredentials, ";")

	credentials.Host, err = url.QueryUnescape(splitCredentials[0])
	if err != nil {
		return err
	}

	credentials.Email, err = url.QueryUnescape(splitCredentials[1])
	if err != nil {
		return err
	}

	credentials.Password, err = url.QueryUnescape(splitCredentials[2])
	if err != nil {
		return err
	}

	return nil
}

func (credentials *Credentials) Load() error {
	log.Info().Msg("Loading data...")

	data := make(map[string]string)

	file, err := os.Open(rawFileName)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Warn().Err(err).Msg("Failed closing credentials file")
		}
	}(file)

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)

		if len(parts) < 2 {
			return errors.New("invalid credentials file")
		}

		data[parts[0]] = parts[1]
	}

	var ok bool

	credentials.Host, ok = data["HOST"]
	if !ok {
		return errors.New("host is missing")
	}

	credentials.Email, ok = data["EMAIL"]
	if !ok {
		return errors.New("email is missing")
	}

	credentials.Password, ok = data["PASSWORD"]
	if !ok {
		return errors.New("password is missing")
	}

	return nil
}

func encryptString(data string, key []byte) ([]byte, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = rand.Read(nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(data), nil)

	return ciphertext, nil
}

func decryptBytes(data []byte, key []byte) (string, error) {
	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(blockCipher)
	if err != nil {
		return "", err
	}

	nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func deriveKey(password string, salt []byte) ([]byte, []byte, error) {
	if salt == nil {
		salt = make([]byte, 32)
		if _, err := rand.Read(salt); err != nil {
			return nil, nil, err
		}
	}

	key, err := scrypt.Key([]byte(password), salt, 1048576, 8, 1, 32)
	if err != nil {
		return nil, nil, err
	}

	return key, salt, nil
}
