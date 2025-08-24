package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type ServerConfigs struct {
	Port         string
	ServiceName  string
	LogLevel     string
	OTLPEndpoint string
	JWTSecret    string
}

type DatabaseConfigs struct {
	Host         string
	Port         string
	User         string
	Password     string
	DBName       string
	SSLMode      string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
}

type RedisConfigs struct {
	Host     string        `mapstructure:"host"`
	Port     string        `mapstructure:"port"`
	Password string        `mapstructure:"password"`
	DB       int           `mapstructure:"db"`
	Prefix   string        `mapstructure:"prefix"`
	TTL      time.Duration `mapstructure:"ttl"`
}

// AppConfigs holds all configs for the service
type AppConfigs struct {
	Server   *ServerConfigs
	Database *DatabaseConfigs
	Redis    *RedisConfigs
}

// GetAppConfigs loads all configs (server + db) and validates them
func GetAppConfigs() (*AppConfigs, error) {
	serverCfg := &ServerConfigs{
		Port:         fmt.Sprintf(":%s", getEnvOrDefault("PORT", "8080")),
		ServiceName:  getEnvOrDefault("SERVICE_NAME", "go-chi-boilerplate"),
		LogLevel:     getEnvOrDefault("LOG_LEVEL", "info"),
		OTLPEndpoint: getEnvOrDefault("OTLP_ENDPOINT", "otelcollector:4317"),
		JWTSecret:    getEnvOrDefault("JWT_SECRET", "supersecret"),
	}

	dbCfg := &DatabaseConfigs{
		Host:         getEnvOrDefault("DB_HOST", ""),
		Port:         getEnvOrDefault("DB_PORT", "5432"),
		User:         getEnvOrDefault("DB_USER", ""),
		Password:     getEnvOrDefault("DB_PASSWORD", ""),
		DBName:       getEnvOrDefault("DB_NAME", ""),
		SSLMode:      getEnvOrDefault("DB_SSLMODE", "disable"),
		MaxOpenConns: getEnvOrDefaultInt("DB_MAX_OPEN_CONNS", 25),
		MaxIdleConns: getEnvOrDefaultInt("DB_MAX_IDLE_CONNS", 25),
		MaxLifetime:  getEnvOrDefaultDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
	}

	if err := dbCfg.Validate(); err != nil {
		return nil, err
	}

	redisCfg := &RedisConfigs{
		Host:     getEnvOrDefault("REDIS_HOST", "redis"),
		Port:     getEnvOrDefault("REDIS_PORT", "6379"),
		Password: getEnvOrDefault("REDIS_PASSWORD", ""),
		DB:       getEnvOrDefaultInt("REDIS_DB", 0),
		Prefix:   getEnvOrDefault("REDIS_PREFIX", "go-chi-boilerplate:"),
		TTL:      getEnvOrDefaultDuration("REDIS_TTL", time.Hour), // default 1h
	}

	if err := redisCfg.Validate(); err != nil {
		return nil, err
	}

	return &AppConfigs{
		Server:   serverCfg,
		Database: dbCfg,
		Redis:    redisCfg,
	}, nil
}

// Validate checks if required DB configs are present
func (d *DatabaseConfigs) Validate() error {
	missing := []string{}

	if d.Host == "" {
		missing = append(missing, "DB_HOST")
	}
	if d.User == "" {
		missing = append(missing, "DB_USER")
	}
	if d.Password == "" {
		missing = append(missing, "DB_PASSWORD")
	}
	if d.DBName == "" {
		missing = append(missing, "DB_NAME")
	}

	if len(missing) > 0 {
		return errors.New("database configuration is incomplete, missing: " + strings.Join(missing, ", "))
	}
	return nil
}

// Validate checks if required Redis configs are present
func (r *RedisConfigs) Validate() error {
	var missing []string

	if r.Host == "" {
		missing = append(missing, "REDIS_HOST")
	}
	if r.Port == "" {
		missing = append(missing, "REDIS_PORT")
	}
	// Password can be optional depending on deployment, so don’t force it
	if r.Prefix == "" {
		missing = append(missing, "REDIS_PREFIX")
	}
	if r.TTL <= 0 {
		missing = append(missing, "REDIS_TTL (must be > 0)")
	}

	if len(missing) > 0 {
		return errors.New("redis configuration is incomplete, missing: " + strings.Join(missing, ", "))
	}
	return nil
}

// getEnvOrDefault returns the value of an environment variable or a default
func getEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// getEnvOrDefaultInt reads an int from env or returns a default
func getEnvOrDefaultInt(key string, defaultValue int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}

// getEnvOrDefaultDuration reads a duration from env or returns a default
func getEnvOrDefaultDuration(key string, defaultValue time.Duration) time.Duration {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultValue
	}
	val, err := time.ParseDuration(valStr)
	if err != nil {
		return defaultValue
	}
	return val
}
