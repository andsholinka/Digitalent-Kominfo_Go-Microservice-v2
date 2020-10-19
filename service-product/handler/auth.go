package handler

import (
	"encoding/json"

	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/service-product/config"
	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/service-product/entity"
	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/utils"
	"github.com/pkg/errors"

	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/context"
)

type AuthMiddleware struct {
	AuthService config.AuthService
}

func (auth *AuthMiddleware) ValidateAuth(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest("POST", auth.AuthService.Host+"/auth/validate", nil)
		if err != nil {
			utils.WrapAPIError(w, r, "failed to create request : "+err.Error(), http.StatusInternalServerError)
			return
		}

		request.Header = r.Header
		authResponse, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Println("ERROR DISINI 1", auth.AuthService.Host)
			utils.WrapAPIError(w, r, "validate auth failed : "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer authResponse.Body.Close()

		body, err := ioutil.ReadAll(authResponse.Body)
		if err != nil {
			log.Println("ERROR DISINI 2")
			utils.WrapAPIError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		var authResult entity.AuthResponse
		err = json.Unmarshal(body, &authResult)

		if authResponse.StatusCode != 200 {

			utils.WrapAPIError(w, r, authResult.ErrorDetails, authResponse.StatusCode)
			return
		}

		context.Set(r, "user", authResult.Data.Username)
		nextHandler(w, r)
	}
}

func (auth *AuthMiddleware) ValidateAuthAdmin(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := http.NewRequest("POST", auth.AuthService.Host+"/auth/validate", nil)
		if err != nil {
			utils.WrapAPIError(w, r, "failed to create request : "+err.Error(), http.StatusInternalServerError)
			return
		}

		request.Header = r.Header
		authResponse, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Println("ERROR DISINI 1", auth.AuthService.Host)
			utils.WrapAPIError(w, r, "validate auth failed : "+err.Error(), http.StatusInternalServerError)
			return
		}
		defer authResponse.Body.Close()

		body, err := ioutil.ReadAll(authResponse.Body)
		if err != nil {
			log.Println("ERROR DISINI 2")
			utils.WrapAPIError(w, r, err.Error(), http.StatusInternalServerError)
			return
		}

		var authResult entity.AuthResponse
		err = json.Unmarshal(body, &authResult)

		if authResponse.StatusCode != 200 {

			utils.WrapAPIError(w, r, authResult.ErrorDetails, authResponse.StatusCode)
			return
		}

		if authResult.Data.Role != 0 {
			utils.WrapAPIError(w, r, errors.Errorf("Not Authorized").Error(), http.StatusBadRequest)
			return
		}

		context.Set(r, "user", authResult.Data.Username)
		nextHandler(w, r)
	}
}
