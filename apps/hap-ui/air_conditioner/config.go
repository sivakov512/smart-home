package airconditioner

import (
	"github.com/pelletier/go-toml/v2"
	"os"
)

type Config struct {
	Manufacturer string
	Name         string
	PIN          string
	Temperature  Temperature
}

type Temperature struct {
	Min  float64
	Max  float64
	Step float64
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
