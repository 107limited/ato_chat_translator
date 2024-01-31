package jwtforreg

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)
var jwtKey = []byte("your_secret_key") // Ganti dengan kunci rahasia Anda

type CustomClaims struct {
	Email     string `json:"email"`
	CompanyID string `json:"company_id"`
	jwt.StandardClaims
}

func CreateTokenOrSession(email, hashedPassword string, companyId int) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token akan berlaku selama 24 jam
	claims := &CustomClaims{
		Email:     email,
		CompanyID: strconv.Itoa(companyId), // Konversi companyId ke string
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateTokenOrSession(tokenString string) (email string, companyId int, err error) {
	claims := &CustomClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return "", 0, err
	}

	if !token.Valid {
		return "", 0, errors.New("invalid token")
	}

	email = claims.Email
	companyId, err = strconv.Atoi(claims.CompanyID) // Konversi kembali ke integer
	if err != nil {
		return "", 0, err // Jika konversi gagal, kembalikan error
	}

	return email, companyId, nil
}
