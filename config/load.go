// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

// Load loads the configuration from the environment.
func Load() (Config, error) {
	config := Config{}
	if err := envconfig.Process("DRONE", &config); err != nil {
		return config, err
	}
	if path := config.Agent.EnvironFile; path != "" {
		envs, _ := godotenv.Read(path)
		for k, v := range envs {
			config.Agent.Environ = append(
				config.Agent.Environ,
				fmt.Sprintf("%s=%s", k, v),
			)
		}
	}
	godotenv.Load()
	return config, nil
}

// MustLoad loads the configuration from the environmnet
// and panics if an error is encountered.
func MustLoad() Config {
	config, err := Load()
	if err != nil {
		panic(err)
	}
	return config
}
