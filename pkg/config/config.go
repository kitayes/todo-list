package config

import (
	"github.com/caarlos0/env/v11"
)

// TODO: разобраться с библиотекой для конфигов
func ReadEnvConfig(cfg any) error {
	err := env.Parse(cfg)
	if err != nil {
		return err
	}
	return nil
}
