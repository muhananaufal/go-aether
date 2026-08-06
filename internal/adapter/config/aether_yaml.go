package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

const ConfigFileName = ".aether.yaml"

type ProjectPreferences struct {
	Architecture string `yaml:"architecture,omitempty"`
	ORM          string `yaml:"orm,omitempty"`
	Engine       string `yaml:"engine,omitempty"`
}

type AetherConfig struct {
	Version     int                `yaml:"version"`
	Preferences ProjectPreferences `yaml:"preferences"`
}

// LoadConfig reads the .aether.yaml file from the current directory if it exists.
func LoadConfig() (*AetherConfig, error) {
	data, err := os.ReadFile(ConfigFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No config file exists, this is fine
		}
		return nil, err
	}

	var cfg AetherConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// SaveConfig writes the given config to .aether.yaml in the current directory.
func SaveConfig(cfg *AetherConfig) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(ConfigFileName, data, 0644)
}
