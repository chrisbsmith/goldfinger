package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v2"
)

func TestLoadConfigs(t *testing.T) {
	dir := t.TempDir()
	scenarios := []struct {
		name           string
		configPath     string            // value to pass as the configPath parameter in LoadConfiguration
		pathAndFiles   map[string]string // files to create in dir
		expectedConfig *Config
		expectedError  error
	}{
		{
			name:       "empty-config-file",
			configPath: filepath.Join(dir, "config.yaml"),
			pathAndFiles: map[string]string{
				"config.yaml": "",
			},
			expectedError: ErrConfigFileNotFound,
		},
		{
			name:          "config-file-that-does-not-exist",
			configPath:    filepath.Join(dir, "config.yaml"),
			expectedError: ErrConfigFileNotFound,
		},
		{
			name:       "config-file-with-no-bonds",
			configPath: filepath.Join(dir, "config.yaml"),
			pathAndFiles: map[string]string{
				"config.yaml": `
bonds:`,
			},
			expectedError: ErrNoBondsInConfig,
		},
		{
			name:       "config-file-with-only-denomination",
			configPath: filepath.Join(dir, "config.yaml"),
			pathAndFiles: map[string]string{
				"config.yaml": `
bonds:
  - denomination: 50`,
			},
			expectedError: ErrConfigIsNotComplete,
		},
		{
			name:       "config-file-with-valid-config",
			configPath: filepath.Join(dir, "config.yaml"),
			pathAndFiles: map[string]string{
				"config.yaml": `
bonds:
  - denomination: 50
    serial: "abcdefg"
    issue_date: "01/2000"
    series: "EE"`,
			},
			expectedConfig: &Config{
				Bonds: []ConfigBond{
					{
						Denomination: 50,
						Serial:       "abcdefg",
						IssueDate:    "01/2000",
						Series:       "EE",
					},
				},
			},
		},
		{
			name:       "config-file-with-extra-bond-field",
			configPath: filepath.Join(dir, "config.yaml"),
			pathAndFiles: map[string]string{
				"config.yaml": `
bonds:
  - denomination: 50
    serial: "abcdefg"
    issue_date: "01/2000"
    series: "EE"
    extra: "field"`,
			},
			expectedConfig: &Config{
				Bonds: []ConfigBond{
					{
						Denomination: 50,
						Serial:       "abcdefg",
						IssueDate:    "01/2000",
						Series:       "EE",
					},
				},
			},
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			for path, content := range scenario.pathAndFiles {
				if err := os.WriteFile(filepath.Join(dir, path), []byte(content), 0644); err != nil {
					t.Fatalf("[%s] failed to write file: %v", scenario.name, err)
				}
			}
			defer func(pathAndFiles map[string]string) {
				for path := range pathAndFiles {
					_ = os.Remove(filepath.Join(dir, path))
				}
			}(scenario.pathAndFiles)
			// config, err := LoadBonds(scenario.configPath)
			config, err := LoadConfig(scenario.configPath)
			if !errors.Is(err, scenario.expectedError) {
				t.Errorf("[%s] expected error %v, got %v", scenario.name, scenario.expectedError, err)
				return
			} else if err != nil && errors.Is(err, scenario.expectedError) {
				return
			}
			// parse the expected output so that expectations are closer to reality (under the right circumstances, even I can be poetic)
			expectedConfigAsYAML, _ := yaml.Marshal(scenario.expectedConfig)
			expectedConfigAfterBeingParsedAndValidated, err := parseAndValidateConfigBytes(expectedConfigAsYAML)
			if err != nil {
				t.Fatalf("[%s] failed to parse expected config: %v", scenario.name, err)
			}
			// Marshal em' before comparing em' so that we don't have to deal with formatting and ordering
			actualConfigAsYAML, err := yaml.Marshal(config)
			if err != nil {
				t.Fatalf("[%s] failed to marshal actual config: %v", scenario.name, err)
			}
			expectedConfigAfterBeingParsedAndValidatedAsYAML, _ := yaml.Marshal(expectedConfigAfterBeingParsedAndValidated)
			if string(actualConfigAsYAML) != string(expectedConfigAfterBeingParsedAndValidatedAsYAML) {
				t.Errorf("[%s] expected config %s, got %s", scenario.name, string(expectedConfigAfterBeingParsedAndValidatedAsYAML), string(actualConfigAsYAML))
			}
		})
	}
}
