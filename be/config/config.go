package config

import (
    "be/constant"
    "be/logger"
    "be/model"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"
    "github.com/harsh-side/keyrotator"
)

var log = logger.Logger
var LimitedConfig *keyrotator.APIKeyConfig

func InitialiseAndGetRootConfig() (*model.Config, error) {
    configPath := os.Getenv(constant.SH10_CONFIG_PATH)
    if configPath == "" {
        configPath = constant.DEF_APP_DIR
    }
    
    // Read main config
    data, err := os.ReadFile(filepath.Join(configPath, constant.APP_FILE_NAME))
    if err != nil {
        log.Error("Couldn't read resource file")
        return nil, fmt.Errorf("couldn't read resource file: %v", err)
    }
    
    config := new(model.Config)
    if err := json.Unmarshal(data, config); err != nil {
        log.Error("Couldn't unmarshal configurations")
        return nil, fmt.Errorf("couldn't unmarshal configurations: %v", err)
    }

    // Initialize API key config
    limitedConfigPath := filepath.Join(configPath, constant.LIMITED_FILE_NAME)
    LimitedConfig, err = keyrotator.NewAPIKeyConfig(limitedConfigPath)
    if err != nil {
        log.Error("Couldn't initialize API key config")
        return nil, fmt.Errorf("couldn't initialize API key config: %v", err)
    }
    log.Info("API key config initialized")
    // Start midnight reset in UTC
    LimitedConfig.StartMidnightReset(time.UTC)
    log.Info("Midnight reset goroutine started")
    return config, nil
}