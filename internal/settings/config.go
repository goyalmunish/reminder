package settings

import (
	"context"
	"fmt"

	"github.com/goyalmunish/reminder/internal/appinfo"
	"github.com/goyalmunish/reminder/pkg/calendar"
	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const currentConfigPath = "./config/current.yaml" // used for current operational values
// const defaultConfigPath = "./config/default.yaml" // just a representation of current default values

type Settings struct {
	AppInfo  *appinfo.Options
	Log      *logger.Options
	Calendar *calendar.Options
}

func DefaultSettings() *Settings {
	return &Settings{
		AppInfo:  appinfo.DefaultOptions(),
		Log:      logger.DefaultOptions(),
		Calendar: calendar.DefaultOptions(),
	}
}

func (s *Settings) String() string {
	value, _ := yaml.Marshal(s)
	return string(value)
}

func LoadConfig(ctx context.Context) (*Settings, error) {
	// set default settings
	settings := DefaultSettings()
	logger.Debug(ctx, fmt.Sprintf("Default Settings: %q", settings))
	viper.SetConfigType("yaml")
	viper.SetConfigFile(currentConfigPath)

	// override with current settings
	logger.Info(ctx, fmt.Sprintf("Attempt to read the app config %q (on top of default values).", currentConfigPath))
	if err := viper.ReadInConfig(); err != nil {
		// Just log the error, and fall back to default settings.
		utils.LogError(ctx, err)
	}
	// If config file is found, unmarshal those values ontop of default settings struct.
	// Otherwise, do nothing.
	if err := viper.Unmarshal(settings); err != nil {
		utils.LogError(ctx, err)
		return nil, err
	}
	logger.Info(ctx, fmt.Sprintf("Final Settings:\n%v", settings))
	return settings, nil
}
