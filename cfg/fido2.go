package cfg

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var VP *viper.Viper
var FD *Fido2

type AuthenticatorSelection struct {
	AuthenticatorAttachment string `json:"authenticator_attachment,omitempty" binding:"oneof=null platform cross-platform ''"`
	UserVerification        string `json:"user_verification,omitempty" binding:"oneof=required preferred discouraged ''"`
	RequireResidentKey      string `json:"require_resident_key,omitempty" binding:"oneof=true false ''"`
}

type ExcludeCredentials struct {
	Transports []string `json:"transports,omitempty" binding:"inArray=usb nfc internal ble ''"`
}

type Register struct {
	AuthenticatorSelection
	ExcludeCredentials
	Timeout     uint   `json:"timeout,omitempty"`
	Attestation string `json:"attestation,omitempty" binding:"oneof=none indirect direct ''"`
}

type Login struct {
	ExcludeCredentials
	UserVerification string `json:"user_verification,omitempty" binding:"oneof=required preferred discouraged ''"`
	Timeout          uint   `json:"timeout,omitempty"`
}

type Fido2 struct {
	Register Register `json:"register"`
	Login    Login    `json:"login"`
}

func InitFido2() (err error) {
	VP = viper.New()
	VP.SetConfigName("fido2")
	VP.SetConfigType("yaml")
	VP.AddConfigPath("etc")
	VP.AddConfigPath("/srun3/etc")
	VP.AddConfigPath("/srun3/bin/etc")
	VP.AddConfigPath(".")
	err = VP.ReadInConfig()
	if err != nil {
		return
	}
	err = VP.Unmarshal(&FD)
	if err != nil {
		return
	}
	VP.WatchConfig()
	VP.OnConfigChange(func(in fsnotify.Event) {
		_ = VP.Unmarshal(&FD)
	})
	return
}
