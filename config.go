package main

import (
	"github.com/kovetskiy/ko"
)

type Config struct {
	IssuerID string `toml:"issuer_id"`
	KeyID    string `toml:"key_id"`
	VendorID string `toml:"vendor_id"`
	KeyPath  string `toml:"key_path"`
}

func loadConfig(path string) (*Config, error) {
	config := &Config{}
	err := ko.Load(path, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
