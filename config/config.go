package config

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
	"os"
)

type Storage struct {
	Verbose    bool
	Debug      bool
	StorageDir string
	IndexDir   string
	Server     struct {
		Host string
		Port int
	}
}

// Configuration in-memory storage, holds all the configuration information parsed from configfiles and flags
var Store Storage = Storage{}

// Loads the Config Store
func (s *Storage) Load(configFile string) {
	if configFile == "" {
		configFile = "docloco"
	}
	s.LoadConfig(configFile)
}


// Load the config Store from a file
func (s *Storage) LoadConfig(configName string) {
	viper.SetConfigName(configName) // name of config file (without extension)
	viper.AddConfigPath(".")        // optionally look for config in the working directory
	err := viper.ReadInConfig()     // Find and read the config file
	if err != nil {                 // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}
	// Verbose & Debug
	s.Verbose = viper.Get("verbose").(bool)
	s.Debug = viper.Get("debug").(bool)

	// Paths
	s.StorageDir = viper.Get("storage_dir").(string)
	s.IndexDir = viper.Get("index_dir").(string)

	if stat, err := os.Stat(s.StorageDir); err != nil || !stat.IsDir() {
		log.Fatal("Storage Directory does not exist, exiting...")
	}

	// Server Config
	s.Server.Host = viper.Get("server.host").(string)
	s.Server.Port = viper.Get("server.port").(int)
}
