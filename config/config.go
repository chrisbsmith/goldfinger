package config

import (
	"errors"
	"log"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

var (
	// ErrConfigFileNotFound is an error returned when a configuration file could not be found
	ErrConfigFileNotFound = errors.New("configuration file not found")

	// ErrNoBondsInConfig is an error returned when no bonds are found in the configuration file
	ErrNoBondsInConfig = errors.New("no bonds found in the configuration")

	// ErrConfigIsNil is an error returned when no bonds are found in the configuration file
	ErrConfigIsNil = errors.New("configuration is nil")

	// ErrConfigIsNotInt is an error returned when a denomination value is not an int value
	ErrConfigIsNotInt = errors.New("denomination is expected to be an int value")

	// ErrConfigIsNotString is an error returned when a serial, issue date, or series is not a string value
	ErrConfigIsNotString = errors.New("serial, issue date and series are expected to be string values")

	// ErrConfigIsNotComplete is an error returned when the configuration is not complete
	ErrConfigIsNotComplete = errors.New("full configuration not provided")
)

// Bond
type ConfigBond struct {
	// Denomination of the Bond
	Denomination int `yaml:"denomination"`

	// Serial number for the bond
	Serial string `yaml:"serial"`

	// Issue Date of the bond
	IssueDate string `yaml:"issue_date"`

	// Series of the bond
	Series string `yaml:"series"`
}

// Config
type Config struct {
	Bonds []ConfigBond `yaml:"bonds"`
}

func LoadConfig(configPath string) (*Config, error) {
	var config *Config
	var configBytes []byte
	log.Printf("Reading configuration from configFile=%s", configPath)

	if data, err := os.ReadFile(configPath); err != nil {
		return config, ErrConfigFileNotFound
	} else {
		configBytes = data
	}

	if len(configBytes) == 0 {
		return config, ErrConfigFileNotFound
	}
	config, err := parseAndValidateConfigBytes(configBytes)
	if err != nil {
		return config, err
	}

	log.Printf("Loaded %d bonds\n", len(config.Bonds))
	return config, nil
}

// parseAndValidateConfigBytes parses the config file
func parseAndValidateConfigBytes(yamlBytes []byte) (config *Config, err error) {
	// Parse configuration file
	if err = yaml.Unmarshal(yamlBytes, &config); err != nil {
		return
	}
	// Check if the configuration file at least has bonds configured
	if config == nil {
		err = ErrConfigIsNil
	} else if config.Bonds == nil || len(config.Bonds) == 0 {
		err = ErrNoBondsInConfig
	} else {
		for _, b := range config.Bonds {
			if reflect.TypeOf(b.Denomination).Kind() != reflect.Int {
				return nil, ErrConfigIsNotInt
			}
			if reflect.TypeOf(b.IssueDate).Kind() != reflect.String {
				return nil, ErrConfigIsNotString
			}
			if reflect.TypeOf(b.Serial).Kind() != reflect.String {
				return nil, ErrConfigIsNotString
			}
			if reflect.TypeOf(b.Series).Kind() != reflect.String {
				return nil, ErrConfigIsNotString
			}
			if b.Denomination == 0 || b.Denomination < 0 {
				return nil, ErrConfigIsNotComplete
			}
			if b.IssueDate == "" || b.Serial == "" || b.Series == "" {
				return nil, ErrConfigIsNotComplete
			}

		}
	}

	return
}
