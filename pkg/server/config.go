package application

import (
	"fmt"
	"os"
)

// ChatServerConfig ...
type ChatServerConfig struct {
	Server   ServerConfig
	Postgres PostgresConfig
}

// ServerConfig ...
type ServerConfig struct {
	Port string
}

// PostgresConfig ...
type PostgresConfig struct {
	Hostname string
	Port     string
	Username string
	Password string
	Database string
}

// Configure ...
func (chat *ChatServer) configure() error {

	config := overrideConfig(getDefaultConfig(), getConfigFromEnvironment())

	if err := checkRequiredConfigsPresent(config); err != nil {
		return err
	}

	chat.config = config

	return nil
}

func getRequiredConfig() []string {
	return []string{
		"server.port",
		"postgres.hostname",
		"postgres.port",
		"postgres.username",
		"postgres.password",
		"postgres.database",
	}
}

func getDefaultConfig() map[string]string {
	return map[string]string{
		"server.port":       "8080",
		"postgres.hostname": "chatserver-postgresql",
		"postgres.port":     "5432",
		"postgres.username": "admin",
		"postgres.database": "chat",
	}
}

func getConfigFromEnvironment() map[string]string {
	return map[string]string{
		"server.port":       os.Getenv("CHATSERVER_SERVICE_PORT_HTTP"),
		"postgres.password": os.Getenv("POSTGRES_PASSWORD"),
	}
}

func overrideConfig(init map[string]string, override map[string]string) map[string]string {
	for config, val := range override {
		init[config] = val
	}
	return init
}

func checkRequiredConfigsPresent(config map[string]string) error {

	var missingConfigs []string
	requiredConfigs := getRequiredConfig()
	for _, required := range requiredConfigs {
		if _, ok := config[required]; !ok {
			missingConfigs = append(missingConfigs, required)
		}
	}

	if len(missingConfigs) > 0 {
		return fmt.Errorf("Missing required configs: %s", missingConfigs)
	}

	return nil
}
