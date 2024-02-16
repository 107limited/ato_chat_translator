package jwtforreg

import (
	
	"time"

	"github.com/dgrijalva/jwt-go"
)
var jwtKey = []byte("your_secret_key") // Ganti dengan kunci rahasia Anda

type CustomClaims struct {
	Email             string `json:"email"`
	CompanyID         int    `json:"company_id"`
	HashedPasswordKey string `json:"hashed_password_key"`
	jwt.StandardClaims
}

func CreateTokenOrSessionWithHashedPasswordKey(email string, companyID int, hashedPasswordKey string) (string, error) {
	// Define your secret key, this should be in an environment variable or configuration file
	var jwtKey = []byte("your_secret_key")

	// Set expiration time for the token
	expirationTime := time.Now().Add(24 * time.Hour)

	// Create the JWT claims, which includes the email, companyID and the hashedPasswordKey alongside the expiration time
	claims := &CustomClaims{
		Email:             email,
		CompanyID:         companyID,
		HashedPasswordKey: hashedPasswordKey,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}



func ValidateTokenAndGetHashedPasswordKey(tokenString string) (string, int, string, error) {
	// Initialize a new instance of `Claims`
	claims := &CustomClaims{}

	// Parse the JWT string and store the result in `claims`.
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// Check if token is valid
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		// Return the extracted info from the token
		return claims.Email, claims.CompanyID, claims.HashedPasswordKey, nil
	} else {
		return "", 0, "", err
	}
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
