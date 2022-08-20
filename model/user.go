package model

import (
	"encoding/binary"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
	"srun/dao/mysql"
	"strings"
)

type User struct {
	gorm.Model
	Name        string
	DisplayName string
	credentials []webauthn.Credential
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

func GetUser(username string) (user User, err error) {
	if err = mysql.GetDB().First(&user, "name = ?", username).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			user.Name = username
			user.DisplayName = strings.SplitN(username, "@", 2)[0]
			err = mysql.GetDB().Create(&user).Error
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

// AddCredential associates the credential to the user
func (u *User) AddCredential(cred webauthn.Credential) error {
	u.credentials = append(u.credentials, cred)
	var c Credential
	return c.AddCredential(u.ID, cred)
}
