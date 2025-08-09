package config

import (
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"strings"
)

var k = koanf.New(".")

func LoadConfig() (*Config, error) {
	var (
		k   = koanf.New(".")
		err error
	)
	//TODO : read from yaml then read from env
	if err = k.Load(env.Provider("HAMRAVESH_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "HAMRAVESH_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	if err = k.Load(env.Provider("HAMRAVESH_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "HAMRAVESH_")), "_", ".", -1)
	}), nil); err != nil {
		return nil, err
	}

	conf := &Config{}
	if err = k.Unmarshal("", conf); err != nil {
		return nil, err
	}

	return conf, nil
}
