package database

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Auth struct {
	ID       int    `gorm:"primary_key" json:"-"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Role     *int   `json:"role,omitempty"`
}

type TokenDetails struct {
	//berisi akes token
	AccessToken string
	//berisi refrsh token
	RefreshToken string
	//waktu expired dari accestoken
	AtExpires int64
	//waktu expired dari refreshtoken
	RtExpires int64
}

func Validate(token string, db *gorm.DB) (*Auth, error) {
	var auth Auth
	// if err := db.Where(&Auth{Token: token}).First(&auth).Error;err != nil{
	// 	if err == gorm.ErrRecordNotFound{
	// 		return nil,errors.Errorf("invalid token")
	// 	}
	// }

	return &auth, nil
}

func (auth *Auth) SignUp(db *gorm.DB) error {

	if err := db.Where(&Auth{Username: auth.Username}).First(auth).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			if err := db.Create(auth).Error; err != nil {
				return err
			}
		} else {
			return err
		}
	}

	return nil
}

func (auth *Auth) Login(db *gorm.DB) (*Auth, error) {
	if err := db.Where(&Auth{Username: auth.Username, Password: auth.Password}).First(auth).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.Errorf("wrong email/password")
		}
	}

	return auth, nil
}
