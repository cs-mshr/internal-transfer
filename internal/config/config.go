package config

import (
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	_ "github.com/joho/godotenv/autoload"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
)

type Config struct {
	Primary  Primary        `koanf:"primary" validate:"required"`
	Server   ServerConfig   `koanf:"server" validate:"required"`
	Database DatabaseConfig `koanf:"database" validate:"required"`
}

type Primary struct {
	Env string `koanf:"env" validate:"required"`
}

type ServerConfig struct {
	Port               string   `koanf:"port" validate:"required"`
	ReadTimeout        int      `koanf:"read_timeout" validate:"required"`
	WriteTimeout       int      `koanf:"write_timeout" validate:"required"`
	IdleTimeout        int      `koanf:"idle_timeout" validate:"required"`
	CORSAllowedOrigins []string `koanf:"cors_allowed_origins" validate:"required"`
}

type DatabaseConfig struct {
	Host            string `koanf:"host" validate:"required"`
	Port            int    `koanf:"port" validate:"required"`
	User            string `koanf:"user" validate:"required"`
	Password        string `koanf:"password"`
	Name            string `koanf:"name" validate:"required"`
	SSLMode         string `koanf:"ssl_mode" validate:"required"`
	MaxOpenConns    int    `koanf:"max_open_conns" validate:"required"`
	MaxIdleConns    int    `koanf:"max_idle_conns" validate:"required"`
	ConnMaxLifetime int    `koanf:"conn_max_lifetime" validate:"required"`
	ConnMaxIdleTime int    `koanf:"conn_max_idle_time" validate:"required"`
}

func LoadConfig() (*Config, error) {
	logger := zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()

	k := koanf.New(".")

	err := k.Load(env.Provider("INTERNAL_TRANSFERS_", ".", func(s string) string {
		result := strings.ToLower(strings.TrimPrefix(s, "INTERNAL_TRANSFERS_"))
		// Use dots to separate namespace levels, keep underscores within field names
		parts := strings.SplitN(result, "_", 2)
		if len(parts) == 2 {
			result = parts[0] + "." + parts[1]
		}
		return result
	}), nil)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not load initial env variables")
	}

	mainConfig := &Config{}

	// Unmarshal into nested structs separately
	err = k.Unmarshal("primary", &mainConfig.Primary)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal primary config")
	}
	
	err = k.Unmarshal("server", &mainConfig.Server)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal server config")
	}
	
	err = k.Unmarshal("database", &mainConfig.Database)
	if err != nil {
		logger.Fatal().Err(err).Msg("could not unmarshal database config")
	}

	validate := validator.New()

	err = validate.Struct(mainConfig)
	if err != nil {
		logger.Fatal().Err(err).Msg("config validation failed")
	}

	return mainConfig, nil
}
