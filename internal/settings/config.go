package settings

import (
	"context"

	"github.com/goyalmunish/reminder/internal/appinfo"
	"github.com/goyalmunish/reminder/pkg/calendar"
	"github.com/goyalmunish/reminder/pkg/logger"
	"github.com/goyalmunish/reminder/pkg/utils"
	"github.com/spf13/viper"
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

func LoadConfig(ctx context.Context) (*Settings, error) {
	settings := DefaultSettings()
	viper.SetConfigFile("./config/default.yaml")
	if err := viper.ReadInConfig(); err != nil {
		utils.LogError(ctx, err)
		return nil, err
	}
	if err := viper.Unmarshal(settings); err != nil {
		utils.LogError(ctx, err)
		return nil, err
	}
	return settings, nil
}
