package settings

import (
	"context"
	"fmt"

	"github.com/goyalmunish/reminder/internal/appinfo"
	"github.com/goyalmunish/reminder/pkg/calendar"
	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

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
	appConfigPath := "./config/default.yaml"
	settings := DefaultSettings()
	viper.SetConfigType("yaml")
	viper.SetConfigFile(appConfigPath)
	if err := viper.ReadInConfig(); err != nil {
		utils.LogError(ctx, err)
		return nil, err
	}
	if err := viper.Unmarshal(settings); err != nil {
		utils.LogError(ctx, err)
		return nil, err
	}
	logger.Info(ctx, fmt.Sprintf("Read the app config %q (on top of default values).", appConfigPath))
	return settings, nil
}
