package helper

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/auth/constant"
	"github.com/FadhlanHawali/Digitalent-Kominfo_Pendalaman-Rest-API/auth/database"
	"github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

//create token
func CreateToken(role int, idUser string) (error, *database.TokenDetails) {
	var roleStr string

	if role == constant.ADMIN {
		roleStr = "admin"
	} else if role == constant.CONSUMER {
		roleStr = "consumer"
	}

	//token details initialization
	td := &database.TokenDetails{}
	//set waktu access token expiry
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	//set waktu refresh token expiry
	td.RtExpires = time.Now().Add(time.Hour).Unix()

	//set header + payload access token
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_user": idUser,
		"role":    role,
		"exp":     td.RtExpires,
	})

	//set salt acces token
	//admin salt -> secret_admin_digitalent
	//consumer salt -> secret_consumer_digitalent
	var err error
	td.AccessToken, err = at.SignedString([]byte(fmt.Sprintf("secret_%s_digitalent", roleStr)))
	if err != nil {
		return err, &database.TokenDetails{}
	}

	//set header + payload refresh token
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id_user": idUser,
		"role":    role,
		"exp":     td.RtExpires,
	})

	//set salt refresh token
	//admin salt -> refresh_secret_admin_digitalent
	//consumer salt -> refresh_secret_consumer_digitalent

	td.RefreshToken, err = rt.SignedString([]byte(fmt.Sprintf("refresh_secret_%s_digitalent", roleStr)))
	if err != nil {
		return err, &database.TokenDetails{}
	}
	return nil, td
}

//create token
func ExtractToken(roles int, r *http.Request) string {
	var bearToken string

	//ambil dari key header
	if roles == constant.ADMIN {
		bearToken = r.Header.Get("digitalent-admin")
	} else if roles == constant.CONSUMER {
		bearToken = r.Header.Get("digitalent-consumer")
	}

	//nge split bearer xxx-xxx-xxx -> array of string
	//array[0]= bearer
	//array[1]= xxx-xxx-xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""

}

//verfifikasi jenis token
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	var roleStr string
	var roles int

	if r.Header.Get("digitalent-admin") != "" {
		roleStr = "admin"
		roles = constant.ADMIN
	} else if r.Header.Get("digitalent-consumer") != "" {
		roleStr = "consumer"
		roles = constant.CONSUMER
	} else {
		return nil, errors.Errorf("Session Invalid")
	}

	tokenString := ExtractToken(roles, r)
	log.Println(tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//cek signIn header apakah hs256
		if jwt.GetSigningMethod("HS256") != token.Method {
			return nil, errors.Errorf("Unexpected signin method: %v", token.Header["alg"])
		}
		return []byte(fmt.Sprintf("secret_%s_digitalent", roleStr)), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil
}

//token validation
func TokenValid(r *http.Request) (string, int, error) {
	//manggil fungsi verifikasi
	token, err := VerifyToken(r)
	if err != nil {
		return "", 0, err
	}

	//proses claim payoad data dari token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		idUser, ok := claims["id_user"].(string)
		role, ok := claims["role"]
		if !ok {
			return "", 0, nil
		}

		return idUser, int(role.(float64)), nil
	}

	return "", 0, nil
}
