package model

import (
	"encoding/binary"
	"errors"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
	"github.com/luguohuakai/north/srun"
	"srun/dao/mysql"
	"strings"
)

type User struct {
	gorm.Model
	Name        string
	DisplayName string
	Status      uint8 // 1:未激活 2:注册完成
	credentials []webauthn.Credential
}

func (User) TableName() string {
	return "wa_user"
}

// WebAuthnID User ID according to the Relying on Party
func (u User) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	return buf
}

// WebAuthnName UserName according to the Relying on Party
func (u User) WebAuthnName() string {
	return u.Name
}

// WebAuthnDisplayName Display Name of the user
func (u User) WebAuthnDisplayName() string {
	return u.DisplayName
}

// WebAuthnIcon User's icon url
func (u User) WebAuthnIcon() string {
	return ""
}

// WebAuthnCredentials Credentials owned by the user
func (u User) WebAuthnCredentials() []webauthn.Credential {
	return u.credentials
}

func GetUser(username string, pwd ...string) (user User, err error) {
	if len(pwd) > 0 && pwd[0] != "" {
		// : 跟北向接口交互 判断用户名/密码是否正确
		var httpResult *srun.HttpResult
		httpResult, err = srun.Request("/api/v1/user/validate-users", "post", map[string]string{"user_name": username, "password": pwd[0]})
		if err != nil {
			return
		} else {
			if httpResult.Code != 0 {
				return User{}, errors.New(httpResult.Message)
			}
		}
		if err = mysql.GetDB().First(&user, "name = ?", username).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				user.Name = username
				user.DisplayName = strings.SplitN(username, "@", 2)[0]
				err = mysql.GetDB().Create(&user).Error
			}
		}
	}
	if err == nil {
		var c []Credential
		err = mysql.GetDB().Find(&c, "uid = ?", user.ID).Error
		if err == nil {
			for _, v := range c {
				credential, _ := v.GetCredential()
				user.credentials = append(user.credentials, credential)
			}
		}
	}

	return
}

func GetLoginUser(username string) (user User, err error) {
	if err = mysql.GetDB().First(&user, "name = ? and status = ?", username, 2).Error; err == nil {
		var c []Credential
		if err = mysql.GetDB().Find(&c, "uid = ?", user.ID).Error; err == nil {
			for _, v := range c {
				credential, _ := v.GetCredential()
				user.credentials = append(user.credentials, credential)
			}
		}
	}

	return
}

// UserIsWebAuthn 验证用户是否已注册webauthn
func UserIsWebAuthn(username string) (err error) {
	var user User
	return mysql.GetDB().Select("id").First(&user, "name = ? and status = ?", username, 2).Error
}

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) error {
	u.credentials = append(u.credentials, cred)
	var c Credential
	return c.AddCredential(u.ID, cred)
}
