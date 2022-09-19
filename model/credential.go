package model

import (
	"encoding/json"
	"errors"
	"github.com/duo-labs/webauthn/webauthn"
	"github.com/jinzhu/gorm"
	"srun/dao/mysql"
)

type Credential struct {
	gorm.Model
	Uid uint   `json:"uid"`
	Cid []byte `json:"cid"`
	//PublicKey  []byte `json:"public_key"`
	Credential string `json:"credential"`
}

func (Credential) TableName() string {
	return "wa_credential"
}

func (c Credential) AddCredential(uid uint, cred webauthn.Credential) error {
	var one Credential
	err := mysql.GetDB().First(&one, "uid = ? and cid = ?", uid, cred.ID).Error
	//err := mysql.GetDB().First(&one, "uid = ? ", uid).Error
	if err != nil && !gorm.IsRecordNotFoundError(err) {
		return err
	}
	if err == nil {
		return errors.New("record existed, can not add again")
	}
	c.Uid = uid
	c.Cid = cred.ID
	//c.PublicKey = cred.PublicKey
	marshal, err := json.Marshal(cred)
	if err != nil {
		return err
	}
	c.Credential = string(marshal)
	return mysql.GetDB().Create(&c).Error
}

func (c Credential) UpdateCredential(newCredential webauthn.Credential) error {
	marshal, err := json.Marshal(newCredential)
	if err != nil {
		return err
	}
	return mysql.GetDB().Model(&c).Update(Credential{Credential: string(marshal)}).Error
}

func (c Credential) GetCredential() (cred webauthn.Credential, err error) {
	err = json.Unmarshal([]byte(c.Credential), &cred)
	return
}
