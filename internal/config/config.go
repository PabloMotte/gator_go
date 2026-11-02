package config

import (
	"os"
	"io"
	"fmt"
	"encoding/json"
	"path/filepath"
)
const configFileName = ".gatorconfig.json"
type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	var config Config
	config_file, err := getConfigFilePath()
	if err != nil {
		err2 := fmt.Errorf("Error getting home directory: %w", err)
		return config, err2
	}
	// Open our jsonFile
	jsonFile, err := os.Open(config_file)
	// if we os.Open returns an error then handle it
	if err != nil {
		err2 := fmt.Errorf("Error opening config file for read: %w", err)
		return config, err2
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		err2 := fmt.Errorf("Error reading config file: %w", err)
		return config, err2
	}

	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		err2 := fmt.Errorf("Error unmarshalling JSON: %w", err)
		return config, err2
	}
	
	return config, nil
	
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user
	err := write(*c)
	if err != nil {
		err2 := fmt.Errorf("Error writing config to file: %w", err)
		return err2
	}
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, configFileName), nil
}

func write(cfg Config) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		err2 := fmt.Errorf("Error marshalling Config data: %w", err)
		return err2
	}
	config_file, err := getConfigFilePath()
	if err != nil {
		err2 := fmt.Errorf("Error getting home directory: %w", err)
		return err2
	}
	// Overwrite our Config File 
	err = os.WriteFile(config_file, data, 0755)
	if err != nil {
		err2 := fmt.Errorf("Error writing config file: %w", err)
		return err2
	}

	return nil
}
