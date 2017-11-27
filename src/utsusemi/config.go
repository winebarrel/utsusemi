package utsusemi

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

const (
	DefaultPort = 11080
	DefaultOk   = 200
)

type Config struct {
	Port    int
	Backend []BackendConfig
}

type BackendConfig struct {
	Target string
	Ok     []int
}

func LoadConfig(flags *Flags) (config *Config, err error) {
	config = &Config{}
	_, err = toml.DecodeFile(flags.Config, config)

	if err != nil {
		return
	}

	if config.Port == 0 {
		config.Port = DefaultPort
	}

	if len(config.Backend) == 0 {
		err = fmt.Errorf("No backends")
		return
	}

	for i := 0; i < len(config.Backend); i++ {
		backend := &config.Backend[i]

		if len(backend.Ok) == 0 {
			backend.Ok = []int{DefaultOk}
		}
	}

	return
}
