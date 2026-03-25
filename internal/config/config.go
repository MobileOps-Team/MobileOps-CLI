package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	EnvProduction  = "production"
	EnvStaging     = "gamma"
	EnvDevelopment = "development"
	EnvTest        = "test"
)

var hostDefaults = map[string]string{
	EnvProduction:  "https://www.mobileops.at",
	EnvStaging:     "https://gamma.mobileops.at/",
	EnvDevelopment: "http://localhost:3000",
	EnvTest:        "http://localhost:3001",
}

// Credentials holds the stored authentication credentials.
type Credentials struct {
	Host        string `json:"host"`
	AccessKey   string `json:"access_key"`
	SecretKey   string `json:"secret_key"`
	Environment string `json:"environment"`
}

// ConfigDir returns the path to the config directory.
func ConfigDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "mobileops")
}

// HostFor returns the default host for the given environment.
func HostFor(environment string) string {
	if h, ok := hostDefaults[environment]; ok {
		return h
	}
	return hostDefaults[EnvProduction]
}

// CurrentEnvironment returns the current environment from the env var or default.
func CurrentEnvironment() string {
	if env := os.Getenv("MOBILEOPS_ENV"); env != "" {
		return env
	}
	return EnvProduction
}

// ResolveEnvironment resolves the environment from flag, env var, or default.
func ResolveEnvironment(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	return CurrentEnvironment()
}

// CredentialsFile returns the path to the credentials file for an environment.
func CredentialsFile(environment string) string {
	if environment == EnvProduction {
		return filepath.Join(ConfigDir(), "credentials.json")
	}
	return filepath.Join(ConfigDir(), fmt.Sprintf("credentials.%s.json", environment))
}

// Load reads credentials for the given environment.
func Load(environment string) (*Credentials, error) {
	file := CredentialsFile(environment)

	data, err := os.ReadFile(file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no credentials found for %s. Run `mobileops auth login --env %s` first", environment, environment)
		}
		return nil, err
	}

	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, fmt.Errorf("credentials file is corrupted. Run `mobileops auth login --env %s` to reconfigure", environment)
	}

	creds.Environment = environment
	return &creds, nil
}

// Save writes credentials for the given environment.
func Save(host, accessKey, secretKey, environment string) error {
	dir := ConfigDir()
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	creds := Credentials{
		Host:        host,
		AccessKey:   accessKey,
		SecretKey:   secretKey,
		Environment: environment,
	}

	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	file := CredentialsFile(environment)
	if err := os.WriteFile(file, data, 0600); err != nil {
		return err
	}

	return nil
}

// Delete removes the credentials file for the given environment.
func Delete(environment string) error {
	file := CredentialsFile(environment)
	err := os.Remove(file)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

// Exists checks whether credentials exist for the given environment.
func Exists(environment string) bool {
	file := CredentialsFile(environment)
	_, err := os.Stat(file)
	return err == nil
}

// Path returns the credentials file path for the given environment.
func Path(environment string) string {
	return CredentialsFile(environment)
}
