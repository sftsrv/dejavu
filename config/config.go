package config

import (
	"encoding/json"
	"os"
	"strings"
)

type Config struct {
	Docs    string   `json:"docs"`
	Types   []string `json:"types"`
	Summary bool     `json:"summary"`
}

func defaultConfig() Config {
	return Config{
		Docs:    "./",
		Summary: false,
	}
}

func loadConfigFile(path string) Config {
	config := defaultConfig()

	file, fileErr := os.ReadFile(path)
	if fileErr != nil {
		return config
	}

	jsonErr := json.Unmarshal(file, &config)
	if jsonErr != nil {
		return config
	}

	return config
}

type Flags struct {
	Path    string
	Docs    string
	Types   string
	Summary bool
}

func Load(flags Flags) Config {

	config := loadConfigFile(flags.Path)

	if flags.Docs != "" {
		config.Docs = flags.Docs
	}

	if flags.Types != "" {
		config.Types = strings.Split(flags.Types, ",")
	}

	if flags.Summary == true {
		config.Summary = true
	}

	return config
}
