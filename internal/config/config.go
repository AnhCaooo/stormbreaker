// AnhCao 2024
package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/AnhCaooo/go-goods/crypto"
	"github.com/AnhCaooo/go-goods/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/constants"
	"github.com/AnhCaooo/stormbreaker/internal/models"
)

var Config models.Config

const configPath string = "/internal/config/"

// load the configuration from the encrypted yaml config file
func ReadFile(cfg *models.Config) error {
	currentDir, err := helpers.GetCurrentDir()
	if err != nil {
		return err
	}
	keyFilePath := currentDir + configPath + constants.CryptoKeyFile
	key, err := crypto.ReadEncryptionKey(keyFilePath)
	if err != nil {
		return err
	}

	encryptedConfigFilePath := currentDir + configPath + constants.EncryptedConfigFile
	decryptedConfigFilePath := currentDir + configPath + constants.DecryptedConfigFile
	if err = crypto.DecryptFile(key, encryptedConfigFilePath, decryptedConfigFilePath); err != nil {
		return err
	}

	f, err := os.Open(decryptedConfigFilePath)
	if err != nil {
		return fmt.Errorf("failed to open config.yml: %s", err.Error())
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		return fmt.Errorf("failed to decode config.yml: %s", err.Error())
	}
	return nil
}
