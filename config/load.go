// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/drone/envconfig"
	"github.com/joho/godotenv"
)

// legacy environment variables. the key is the legacy
// variable name, and the value is the new variable name.
var legacy = map[string]string{
	"DRONE_ENABLE_PINGER": "DRONE_PINGER_ENABLED",
	"DRONE_ENABLE_REAPER": "DRONE_REAPER_ENABLED",
}

func init() {
	// loop through legacy environment variable and, if set
	// rewrite to the new variable name.
	for k, v := range legacy {
		if s, ok := os.LookupEnv(k); ok {
			os.Setenv(v, s)
		}
	}
}

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
	// If environment variables don't contain `=`, we consider that it's an environment name, we fetch and expose the value
	for i, env := range config.Agent.Environ {
		if !strings.Contains(env, "=") {
			config.Agent.Environ[i] = fmt.Sprintf("%s=%s", env, os.Getenv(env))
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
