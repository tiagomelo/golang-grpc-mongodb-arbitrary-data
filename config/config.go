// Copyright (c) 2023 Tiago Melo. All rights reserved.
// Use of this source code is governed by the MIT License that can be found in
// the LICENSE file.
//
// Package config provides functions for reading and processing the application configuration.
// It reads environment variables from a file and populates a Config struct with the values.
// The configuration struct holds all the necessary configuration values needed by the application.
package config

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

// Config holds all the configuration needed by the application.
type Config struct {
	MongodbDatabase     string `envconfig:"MONGODB_DATABASE" required:"true"`
	MongodbHostName     string `envconfig:"MONGODB_HOST_NAME" required:"true"`
	MongodbPort         int    `envconfig:"MONGODB_PORT" required:"true"`
	MongodbTestDatabase string `envconfig:"MONGODB_TEST_DATABASE" required:"true"`
	MongodbTestHostName string `envconfig:"MONGODB_TEST_HOST_NAME" required:"true"`
	MongodbTestPort     int    `envconfig:"MONGODB_TEST_PORT" required:"true"`
	GrpcServerṔort      int    `envconfig:"GRPC_SERVER_PORT" required:"true"`
}

// For ease of unit testing.
var (
	godotenvLoad     = godotenv.Load
	envconfigProcess = envconfig.Process
)

// Read reads the environment variables from the given file and returns a Config.
func Read(envFilePath string) (*Config, error) {
	if err := godotenvLoad(envFilePath); err != nil {
		return nil, errors.Wrap(err, "loading env vars")
	}
	config := new(Config)
	if err := envconfigProcess("", config); err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}
	return config, nil
}
