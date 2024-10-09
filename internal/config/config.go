package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/AnhCaooo/stormbreaker/internal/helpers"
	"github.com/AnhCaooo/stormbreaker/internal/models"
)

var Config models.Config

// load the configuration from the yaml config file
func ReadFile(cfg *models.Config) error {
	currentDir, err := helpers.GetCurrentDir()
	if err != nil {
		return err
	}
	configPath := fmt.Sprintf("%s/internal/config/config.yml", currentDir)
	f, err := os.Open(configPath)
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
