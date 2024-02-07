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

func CreateTokenOrSession(email string, companyId int) (string, error) {
    // Atur waktu kedaluwarsa untuk token
    expirationTime := time.Now().Add(24 * time.Hour) // Contoh: 24 jam kedaluwarsa

    // Konversi companyId dari int ke string
    companyIdStr := strconv.Itoa(companyId)

    // Buat klaim dengan email pengguna dan CompanyID sebagai string
    claims := &CustomClaims{
        Email:     email,
        CompanyID: companyIdStr, // CompanyID sebagai string
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    // Buat token dengan klaim yang telah ditetapkan
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Tandatangani token dengan kunci rahasia Anda
    tokenString, err := token.SignedString([]byte("your_secret_key")) // Ganti "your_secret_key" dengan kunci rahasia yang sesungguhnya
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



// Claims adalah struktur yang akan digunakan untuk menambahkan klaim ke token.
type Claims struct {
    Email string `json:"email"`
    UserID int `json:"user_id"`
    jwt.StandardClaims
}

// CreateToken menghasilkan token JWT untuk pengguna dengan email dan ID tertentu.
func CreateToken(email string, userID int) (string, error) {
    // Waktu kedaluwarsa token, misalnya 24 jam dari saat ini
    expirationTime := time.Now().Add(24 * time.Hour)

    // Buat klaim yang berisi informasi pengguna dan waktu kedaluwarsa
    claims := &Claims{
        Email: email,
        UserID: userID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }

    // Deklarasikan token dengan algoritma HS256 dan klaim
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    // Tandatangani token dengan kunci rahasia
    tokenString, err := token.SignedString(jwtKey)
    if err != nil {
        return "", err
    }

    return tokenString, nil
}
