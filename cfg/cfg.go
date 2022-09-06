package cfg

import (
	"fmt"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func Init() (err error) {
	viper.SetConfigName("webauthn")
	viper.SetConfigType("yaml")
	//viper.SetConfigFile("/srun3/bin/etc/prod/webauthn.yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./etc/dev/")
	viper.AddConfigPath("./etc/prod/")
	viper.AddConfigPath("/srun3/etc/")
	err = viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.WatchConfig()
	viper.OnConfigChange(func(in fsnotify.Event) {})

	fmt.Println(fmt.Sprintf("%s starting..., port: %d, mode: %s ....", viper.GetString("app.name"), viper.GetInt("app.port"), viper.GetString("app.mode")))

	return
}

var WAWeb *webauthn.WebAuthn

func InitWebAuthn() (err error) {
	WAWeb, err = webauthn.New(&webauthn.Config{
		RPDisplayName: viper.GetString("app.name"), // Display Name for your site
		RPID:          viper.GetString("app.host"), // Generally the FQDN for your site
		RPOrigin: fmt.Sprintf(
			"%s://%s:%d",
			viper.GetString("app.protocol"),
			viper.GetString("app.host"),
			viper.GetInt("app.port")), // The origin URL for WebAuthn requests
		//RPIcon: "https://duo.com/logo.png", // Optional icon URL for your site
	})
	return
}
