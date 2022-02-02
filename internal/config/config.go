package config

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type Config struct {
	Database map[string]string `json:"databaseMapping"`
	Color    map[string]string `json:"colorMapping"`
}

func New(configPath string) *Config {
	jsonFile, err := os.Open(configPath)
	if err != nil {
		log.Fatalf("Error occurred: %v\n", err)
	}

	defer jsonFile.Close()

	data, _ := io.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(data, &config)

	return &config
}

func (config *Config) DatabaseNames() []string {
	var names []string

	for _, databaseMapping := range config.Database {
		names = append(names, databaseMapping)
	}

	return names
}
