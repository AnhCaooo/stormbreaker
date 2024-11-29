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

// load the configuration from the encrypted yaml config file
func LoadFile(cfg *models.Config) error {
	if err := readFile(cfg); err != nil {
		return err
	}
	// Validate configuration data
	if err := cfg.Validate(); err != nil {
		return err
	}
	return nil
}

// ReadFile reads the encrypted configuration file then decrypt and read it
func readFile(cfg *models.Config) error {
	currentDir, err := helpers.GetCurrentDir()
	if err != nil {
		return err
	}
	keyFilePath := currentDir + constants.CryptoKeyFile
	key, err := crypto.ReadEncryptionKey(keyFilePath)
	if err != nil {
		return err
	}

	encryptedConfigFilePath := currentDir + constants.EncryptedConfigFile
	decryptedConfigFilePath := currentDir + constants.DecryptedConfigFile

	// *Note*: uncomment this to get the encrypted config file update if there is some configs updated in config.yml
	// originalConfig := currentDir + constants.ConfigFile
	// if err = crypto.EncryptFile(key, originalConfig, encryptedConfigFilePath); err != nil {
	// 	return err
	// }

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
