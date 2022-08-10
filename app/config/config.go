package config

import (
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/mitchellh/go-homedir"
)

type Config struct {
	Service    ServiceConfig
	HttpServer HttpServerConfig
	//HttpClient HttpClientConfig
	Database DatabaseConfig
}

// Returns config from a file or environment variables.
// Env vars have priority if they exist.
func GetConfig(configFileDirName string) *Config {

	var configData Config

	// File - Env
	configDir, fileErr := getDefaultConfigDir(configFileDirName)
	if fileErr != nil {
		log.Println("Error getting default config directory:", fileErr)
	} else {
		configFilePath := filepath.Join(configDir, "config.toml")

		// Reads both file and env vars in that order, but fails when the file
		// does not exist or the requested parameters do not exist in it.
		// Also Env values overwrite those from the file.
		fileErr = cleanenv.ReadConfig(configFilePath, &configData)
		if fileErr == nil {
			return &configData
		}
		// If there's an err will be logged later only if the app cannot run.
	}

	// Env
	// library note: env-required - flag to mark a field as required. If set will
	// return an error during environment parsing when the flagged as required
	// field is empty (default Go value).
	envErr := cleanenv.ReadEnv(&configData)
	if envErr != nil {
		log.Println("Error reading config file:", fileErr)
		log.Panicln("Error reading config environment variables:", envErr)
	}
	return &configData
}

func getDefaultConfigDir(configDirName string) (string, error) {
	var configDirLocation string

	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}

	switch runtime.GOOS {
	case "linux":
		// Use the XDG_CONFIG_HOME variable if it is set, otherwise
		// $HOME/.config/example
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
			configDirLocation = xdgConfigHome
		} else {
			configDirLocation = filepath.Join(homeDir, ".config", configDirName)
		}

	default:
		// On other platforms we just use $HOME/.example
		hiddenConfigDirName := "." + configDirName
		configDirLocation = filepath.Join(homeDir, hiddenConfigDirName)
	}

	return configDirLocation, nil
}
