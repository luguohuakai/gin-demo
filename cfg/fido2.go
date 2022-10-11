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
	RequireResidentKey      bool   `json:"require_resident_key,omitempty" binding:"oneof=true false ''"`
}

type ExcludeCredentials struct {
	Transports Transports `json:"transports" binding:"required,inArray=usb nfc internal ble"`
}

type Transports []string

type Register struct {
	AuthenticatorSelection AuthenticatorSelection `json:"authenticator_selection,omitempty"`
	ExcludeCredentials     ExcludeCredentials     `json:"exclude_credentials,omitempty"`
	Timeout                uint                   `json:"timeout,omitempty"`
	Attestation            string                 `json:"attestation,omitempty" binging:"oneof=none indirect direct ''"`
}

type Login struct {
	UserVerification   string             `json:"user_verification,omitempty" binding:"oneof=required preferred discouraged ''"`
	ExcludeCredentials ExcludeCredentials `json:"exclude_credentials"`
	Timeout            uint               `json:"timeout"`
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
