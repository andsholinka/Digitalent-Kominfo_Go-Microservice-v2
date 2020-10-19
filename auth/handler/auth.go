package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/auth/database"
	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/auth/helper"
	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/utils"
	"gorm.io/gorm"
)

type Auth struct {
	Db *gorm.DB
}

func (db *Auth) ValidateAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	idUser, role, err := helper.TokenValid(r)
	if err != nil {
		utils.WrapAPIError(w, r, err.Error(), http.StatusUnauthorized)

		return
	}

	utils.WrapAPIData(w, r, database.Auth{
		Username: idUser,
		Role:     &role,
	}, http.StatusOK, "success")

	return
}

func (db *Auth) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.WrapAPIError(w, r, "can't read body", http.StatusBadRequest)
		return
	}

	var signup database.Auth

	err = json.Unmarshal(body, &signup)
	if err != nil {
		utils.WrapAPIError(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = signup.SignUp(db.Db)
	if err != nil {
		utils.WrapAPIError(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	utils.WrapAPISuccess(w, r, "success", http.StatusOK)
}

func (db *Auth) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		utils.WrapAPIError(w, r, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		utils.WrapAPIError(w, r, "can't read body", http.StatusBadRequest)
		return
	}

	var login database.Auth

	err = json.Unmarshal(body, &login)
	if err != nil {
		utils.WrapAPIError(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	res, err := login.Login(db.Db)
	if err != nil {
		utils.WrapAPIError(w, r, "error unmarshal : "+err.Error(), http.StatusInternalServerError)
		return
	}

	err, token := helper.CreateToken(*res.Role, res.Username)
	utils.WrapAPIData(w, r, token, http.StatusOK, "success")
	return
}
