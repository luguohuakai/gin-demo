package cfg

import (
	"github.com/duo-labs/webauthn/protocol"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var VP *viper.Viper
var FD *Fido2

type AuthenticatorSelection struct {
	AuthenticatorAttachment protocol.AuthenticatorAttachment     `json:"authenticator_attachment,omitempty" mapstructure:"authenticator_attachment" binding:"oneof=null platform cross-platform ''"`
	UserVerification        protocol.UserVerificationRequirement `json:"user_verification,omitempty" mapstructure:"user_verification" binding:"oneof=required preferred discouraged ''"`
	RequireResidentKey      string                               `json:"require_resident_key,omitempty" mapstructure:"require_resident_key" binding:"oneof=true false ''"`
}

type ExcludeCredentials struct {
	Transports []protocol.AuthenticatorTransport `json:"transports,omitempty" mapstructure:"transports" binding:"inArray=usb nfc internal ble ''"`
}

type Register struct {
	AuthenticatorSelection `json:"authenticator_selection" mapstructure:"authenticator_selection"`
	ExcludeCredentials     `json:"exclude_credentials,omitempty" mapstructure:"exclude_credentials,omitempty"`
	Timeout                uint                          `json:"timeout,omitempty" mapstructure:"timeout"`
	Attestation            protocol.ConveyancePreference `json:"attestation,omitempty" mapstructure:"attestation" binding:"oneof=none indirect direct ''"`
}

type Login struct {
	AllowCredentials ExcludeCredentials                   `json:"allow_credentials,omitempty" mapstructure:"allow_credentials"`
	UserVerification protocol.UserVerificationRequirement `json:"user_verification,omitempty" mapstructure:"user_verification" binding:"oneof=required preferred discouraged ''"`
	Timeout          uint                                 `json:"timeout,omitempty" mapstructure:"timeout"`
}

type Fido2 struct {
	Register Register `json:"register" mapstructure:"register"`
	Login    Login    `json:"login" mapstructure:"login"`
}

func InitFido2() (err error) {
	VP = viper.New()
	VP.SetConfigName("fido2")
	VP.SetConfigType("yaml")
	VP.AddConfigPath("etc")
	VP.AddConfigPath("/srun3/etc")
	VP.AddConfigPath("/srun3/bin/webauthn/etc")
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
