package model

import (
	"encoding/binary"
	"errors"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
	"srun/dao/mysql"
)

func (Admin) TableName() string {
	return "wa_admin"
}

type Admin struct {
	gorm.Model
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Status      uint8  `json:"status,omitempty"`
	credentials []webauthn.Credential
}

type QueryAdmin struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Avatar   string `json:"avatar,omitempty"`
	Status   uint8  `json:"status,omitempty"`
	Page     int    `json:"page,omitempty"`
	Size     int    `json:"size,omitempty"`
}

func (a QueryAdmin) GetAdminLst() (lst []AdminLst, total int, err error) {
	if a.Page == 0 {
		a.Page = 1
	}
	if a.Size == 0 {
		a.Size = 20
	}
	db := mysql.GetDB().Model(&Admin{})
	if a.Username != "" {
		db = db.Where("username like %?%", a.Username)
	}
	if a.Status != 0 {
		db = db.Where("status = ?", a.Status)
	}
	if err = db.Count(&total).Error; err != nil {
		return
	}
	err = db.Order("id DESC").Offset((a.Page - 1) * a.Size).Limit(a.Size).Find(&lst).Error

	return
}

type AdminLst struct {
	Id          uint
	CreatedAt   Date
	Name        string
	DisplayName string
	Status      uint8
}

func (*AdminLst) TableName() string {
	return "wa_admin"
}

// WebAuthnID Admin ID according to the Relying on Party
func (a Admin) WebAuthnID() []byte {
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(a.ID))
	return buf
}

// WebAuthnName AdminName according to the Relying on Party
func (a Admin) WebAuthnName() string {
	return a.Username
}

// WebAuthnDisplayName Display Name of the admin
func (a Admin) WebAuthnDisplayName() string {
	return a.Username
}

// WebAuthnIcon Admins icon url
func (a Admin) WebAuthnIcon() string {
	return a.Avatar
}

// WebAuthnCredentials Credentials owned by the admin
func (a Admin) WebAuthnCredentials() []webauthn.Credential {
	return a.credentials
}

func GetAdmin(username, action string, pwd ...string) (admin Admin, err error) {
	if action == "begin" {
		if len(pwd) != 1 || pwd[0] == "" {
			return Admin{}, errors.New("password can not be empty")
		}
	}
	if len(pwd) > 0 && pwd[0] != "" {
		if err = mysql.GetDB().First(&admin, "username = ?", username).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				admin.Username = username
				err = mysql.GetDB().Create(&admin).Error
			}
		}
	}
	err = mysql.GetDB().First(&admin, "username = ?", username).Error
	if err == nil {
		var c []AdminCredential
		err = mysql.GetDB().Find(&c, "uid = ?", admin.ID).Error
		if err == nil {
			for _, v := range c {
				credential, _ := v.GetCredential()
				admin.credentials = append(admin.credentials, credential)
			}
		}
	}

	return
}

func GetLoginAdmin(username string) (admin Admin, err error) {
	if err = mysql.GetDB().First(&admin, "name = ? and status = ?", username, 2).Error; err == nil {
		var c []Credential
		if err = mysql.GetDB().Find(&c, "uid = ?", admin.ID).Error; err == nil {
			if len(c) == 0 {
				return Admin{}, errors.New("no credentials found, please register first")
			}
			for _, v := range c {
				credential, _ := v.GetCredential()
				admin.credentials = append(admin.credentials, credential)
			}
		}
	}

	return
}

// AdminIsWebAuthn 验证用户是否已注册webauthn
func AdminIsWebAuthn(username string) (err error) {
	var admin Admin
	return mysql.GetDB().Select("id").First(&admin, "username = ? and status = ?", username, 3).Error
}

// AddCredential associates the credential to the admin
func (a *Admin) AddCredential(cred webauthn.Credential) error {
	a.credentials = append(a.credentials, cred)
	var c Credential
	return c.AddCredential(a.ID, cred)
}
