package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"

	config "spender/v1/api/config"
)

type Payload struct {
	Name string
	Email    string
	Id       uint
}

type Claims struct {
	Name	 string `json:"name"`
	Email    string `json:"email"`
	Id       uint   `json:"id"`
	jwt.StandardClaims
}

var JWT_SECRET string

func GenerateJwtToken(payload Payload) (string, error) {
	// if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
	// 	log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	// }
	JWT_SECRET = config.JWT_SECRET

	key := []byte(JWT_SECRET)

	//7 * 24 * 60 : 7 days 
	//in mins
	expirationTime := config.TokenExpiredTime

	claims := &Claims{
		Id:       payload.Id,
		Name: payload.Name,
		Email:    payload.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	UnsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	SignedToken, err := UnsignedToken.SignedString(key)
	if err != nil {
		return "", err
	}

	return SignedToken, nil
}

func VerifyJwtToken(strToken string) (*Claims, error) {
	// if JWT_SECRET = os.Getenv("JWT_SECRET"); JWT_SECRET == "" {
	// 	log.Fatal("[ ERROR ] JWT_SECRET environment variable not provided!\n")
	// }
	JWT_SECRET = config.JWT_SECRET

	key := []byte(JWT_SECRET)

	claims := &Claims{}

	token, err := jwt.ParseWithClaims(strToken, claims, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return claims, fmt.Errorf("invalid token signature")
		}
		
	}

	if token == nil {
		return claims, fmt.Errorf("token error")
	} else {
		if !token.Valid {
			return claims, fmt.Errorf("invalid token")
		}
	}


	return claims, nil
}