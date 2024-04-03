package config

import (
	"io"
	"lda/logging"
	"os"

	"github.com/spf13/viper"
)

// Config config definition
type Config struct {
	// Debug persists the debug mode so we don't have to pass it via flag, flag will override this
	Debug bool `mapstructure:"debug"`
	// ProcessInterval interval in seconds to tick and collect general information about processes - defaults to 120 seconds
	ProcessInterval int `mapstructure:"process_interval"`
	// CommandInterval interval in which to collect process information when command has been executed - defaults to 1 second
	CommandInterval int `mapstructure:"command_interval"`
	// CommandIntervalMultiplier multiplier for the command interval - defaults to 5 seconds
	CommandIntervalMultiplier int `mapstructure:"command_interval_multiplier"`
	// MaxConcurrentCommands maximum number of concurrent commands to collect - defaults to 20
	MaxConcurrentCommands int `mapstructure:"max_concurrent_commands"`
	// RemoteCollection flag to enable remote collection - defaults to false
	RemoteCollection bool `mapstructure:"remote_collection"`
	// ServerAddress host to connect to for remote collection
	ServerAddress string `mapstructure:"server_host"`
	// SecureConnection flag to enable secure connection to the server
	SecureConnection bool `mapstructure:"secure_connection"`
	// CertFile path to the certificate file
	CertFile string `mapstructure:"cert_file"`
	// ExcludeRegex regular expression to exclude processes from collection
	ExcludeRegex string `mapstructure:"exclude_regex"`
	// ProcessCollectionType type of process collection to use, ps or psutil
	ProcessCollectionType string `mapstructure:"process_collection_type"`

	// Configuration that is not available via the configuration file

	// Out is the output writer for printing information
	Out io.Writer
	// ErrOut is the error output writer for printing errors
	ErrOut io.Writer
}

// AppConfig is the global configuration instance
var AppConfig *Config

// SetupConfig initialize the configuration instance
func SetupConfig() {

	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(LdaDir)

	// Set default configuration values
	var config = &Config{
		Debug:                     false,
		RemoteCollection:          false,
		ProcessInterval:           3600,
		CommandInterval:           1,
		CommandIntervalMultiplier: 5,
		MaxConcurrentCommands:     20,
		ProcessCollectionType:     "ps",
		// Set default output writers - currently we only use stdout and stderr, but this could potentially be changed
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	if err := viper.ReadInConfig(); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to read config file")
	}

	if err := viper.Unmarshal(config); err != nil {
		logging.Log.Error().Err(err).Msg("Failed to unmarshal config file")
	}

	AppConfig = config
}
