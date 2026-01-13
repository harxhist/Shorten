package config

import (
	"be/constant"
	"be/model"
	"errors"
	"fmt"
)

// ValidateConfig checks if the provided Config object has necessary values set.
func ValidateConfig(cfg *model.Config) error {
	if cfg == nil {
		return errors.New("configuration cannot be nil")
	}
	if cfg.LokiURL == "" {
		return errors.New("LOKI_URL cannot be empty")
	}
	if cfg.AppName == "" {
		return errors.New("APP_NAME cannot be empty")
	}
	if cfg.Environment == "" {
		return errors.New("ENVIRONMENT cannot be empty")
	}
	if len(cfg.AllowedOrigin) == 0 {
		return errors.New("ALLOWED_ORIGIN must have at least one entry")
	}
	return nil
}

// SetConfig sets the application configuration, validating it beforehand.
func SetConfig(cfg *model.Config) error {
	if err := ValidateConfig(cfg); err != nil {
		return fmt.Errorf("couldn't validate configuration: %v", err)
	}
	constant.APPCONFIG = cfg
	log.Info("Configuration set successfully")
	return nil
}
