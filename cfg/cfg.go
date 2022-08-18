package cfg

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./etc/")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {})
	return
}
