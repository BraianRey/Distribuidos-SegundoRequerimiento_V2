package utilidades

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	UserID            string `json:"user_id"`
	ServidorCanciones string `json:"servidor_canciones"`
	ServidorStreaming string `json:"servidor_streaming"`
}

var configInstance *Config

func GetConfig() (*Config, error) {
	if configInstance != nil {
		return configInstance, nil
	}

	// Buscar el archivo config.json desde el directorio actual hacia arriba
	configPath := "config.json"
	for i := 0; i < 3; i++ { // Buscar hasta 3 niveles arriba
		if _, err := os.Stat(configPath); err == nil {
			break
		}
		configPath = filepath.Join("..", configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, err
	}

	configInstance = config
	return config, nil
}

func UpdateUserID(userID string) error {
	config, err := GetConfig()
	if err != nil {
		return err
	}

	config.UserID = userID

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return err
	}

	// Buscar el archivo config.json desde el directorio actual hacia arriba
	configPath := "config.json"
	for i := 0; i < 3; i++ { // Buscar hasta 3 niveles arriba
		if _, err := os.Stat(configPath); err == nil {
			break
		}
		configPath = filepath.Join("..", configPath)
	}

	return os.WriteFile(configPath, data, 0644)
}
