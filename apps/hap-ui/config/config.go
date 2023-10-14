package config

import (
	"github.com/pelletier/go-toml/v2"
	"hap-ui/ac"
	"os"
)

type Config struct {
	Broker string
	PIN    string
	AC     *ac.Config
}

func LoadConfig(fpath string) (*Config, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var config Config

	d := toml.NewDecoder(f)
	d.DisallowUnknownFields()

	err = d.Decode(&config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
