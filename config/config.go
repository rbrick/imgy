package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	defaultConfig = &Config{
		Host:           "localhost",
		Port:           "8080",
		CookieStoreKey: "",
		DatabaseConfig: &DBConfig{
			Path: "imgy.db",
		},
		TLSEnabled: false,
		TLSConfig: &TLSConfig{
			CertPath: "imgy.cert",
			KeyPath:  "imgy.key",
		},
		AWSConfig: &AWSConfig{
			Bucket: "imgy-s3",
			Region: "us-west-2",
		},
		GoogleAuth: &GoogleAuthConfig{
			JsonPath: "auth_creds.json",
		},
	}
)

type Config struct {
	Host              string            `json:"host"`
	Port              string            `json:"port"`
	CookieStoreKey    string            `json:"cookieStoreKey"`
	DatabaseConfig    *DBConfig         `json:"database"`
	AWSConfig         *AWSConfig        `json:"aws"`
	TLSEnabled        bool              `json:"tlsEnabled"`
	TLSConfig         *TLSConfig        `json:"tls"`
	GoogleAuthEnabled bool              `json:"googleAuthEnabled"`
	GoogleAuth        *GoogleAuthConfig `json:"googleAuth"`
}

type DBConfig struct {
	Path string `json:"path"`
}

type TLSConfig struct {
	CertPath string `json:"cert"`
	KeyPath  string `json:"path"`
}

type AWSConfig struct {
	Bucket string `json:"bucket"`
	Region string `json:"region"`
}

type GoogleAuthConfig struct {
	JsonPath string `json:"jsonPath"`
}

// Open opens a config file
func Open(path string) (*Config, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return writeDefault(path)
	}

	// open the file
	f, err := ioutil.ReadFile(path)

	if err != nil {
		return nil, err
	}

	var config Config

	if err = json.Unmarshal(f, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

func writeDefault(path string) (*Config, error) {
	// write the default config
	log.Println("Path not found, writing default config")

	data, err := json.MarshalIndent(defaultConfig, "", "    ")

	if err != nil {
		log.Fatalln("Failed to create config:")
		log.Fatalln(err)
	}

	ioutil.WriteFile(path, data, os.ModePerm)
	return defaultConfig, nil
}

// Default returns the default config
func Default() *Config {
	return defaultConfig
}
