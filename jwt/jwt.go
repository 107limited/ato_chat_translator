package jwt

import (
	"ato_chat/config"
	"ato_chat/dbAto"
	"ato_chat/models"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a JWT token based on the provided user data
func GenerateToken(user *models.User, privateKey *rsa.PrivateKey) (string, error) {
	header := models.JwtHeader{
		Alg: "RS256",
		Typ: "JWT",
	}

	nowTime := time.Now().Add(time.Hour * config.DiffUTC)
	payload := models.JwtPayload{
		UserId:   user.UserID,
		Com:      user.Company,
		ComId:    user.CompanyID,
		Dep:      user.Department,
		Nam:      user.Name,
		Address:  user.Address,
		Lng:      0.0,
		Lat:      0.0,
		Iat:      nowTime,
		Exp:      nowTime.Add(time.Minute * config.ExpireMinute),
		Iss:      "107",
		Aut:      user.Auth,
		SiteType: "orikomi_manage",
	}

	headerBytes, err := json.Marshal(header)
	if err != nil {
		log.Println(err)
		return "", err
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
		return "", err
	}

	encodedHeader := base64.URLEncoding.EncodeToString(headerBytes)
	encodedPayload := base64.URLEncoding.EncodeToString(payloadBytes)
	headPay := encodedHeader + "." + encodedPayload
	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(headPay))
	if err != nil {
		log.Println(err)
		return "", err
	}
	msgHashSum := msgHash.Sum(nil)
	signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
	if err != nil {
		log.Println(err)
		return "", err
	}
	encodedSignature := base64.URLEncoding.EncodeToString(signature)
	token := "Bearer " + headPay + "." + encodedSignature
	log.Println("token======", token)
	return token, nil
}

// func signTokenWithRSAPrivateKey(data string, privateKey *rsa.PrivateKey) (string, error) {
//     msgHash := sha256.New()
//     _, err := msgHash.Write([]byte(data))
//     if err != nil {
//         log.Println(err)
//         return "", err
//     }
//     msgHashSum := msgHash.Sum(nil)
//     signature, err := rsa.SignPSS(rand.Reader, privateKey, crypto.SHA256, msgHashSum, nil)
//     if err != nil {
//         log.Println(err)
//         return "", err
//     }
//     encodedSignature := base64.URLEncoding.EncodeToString(signature)
//     token := data + "." + encodedSignature
//     log.Println("token======", token)
//     return token, nil
// }

func Verify(token string) bool {
	if len(token) < 20 {
		return false
	}
	spaceLastIndex := strings.LastIndex(token, " ")
	if len(token) < spaceLastIndex+2 {

	} else {
		token = token[spaceLastIndex+1:]
	}

	lastCommaIndex := strings.LastIndex(token, ".")
	if len(token) < lastCommaIndex+2 {
		return false
	}

	log.Println("lastCommaIndex=", lastCommaIndex)

	signature := token[lastCommaIndex+1:]
	log.Println(signature)
	verifyPart := token[:lastCommaIndex]

	log.Println("verifyPart======", verifyPart)

	log.Println("signature=====", signature)

	decodedSignature, err := base64.URLEncoding.DecodeString(signature)
	if err != nil {
		// Handle the error here, e.g., log it or return an error response.
		log.Println("Error decoding signature:", err)
		// You may want to return an error response to the client.
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return false
	}

	msgHash := sha256.New()
	_, err = msgHash.Write([]byte(verifyPart))
	if err != nil {
		panic(err)
	}
	msgHashSum := msgHash.Sum(nil)

	privateKey := dbAto.PrivateKey
	err = rsa.VerifyPSS(&privateKey.PublicKey, crypto.SHA256, msgHashSum, decodedSignature, nil)
	if err != nil {
		log.Println("could not verify signature: ", err)
		return false
	}
	// If we don't get any error from the `VerifyPSS` method, that means our
	// signature is valid
	log.Println("signature verified")
	//readToken(verifyPart)
	return true
}

func VerifyExpired() bool {
	return true
}

func ReadToken(token string) models.JwtPayload {
	payloadModel := &models.JwtPayload{}
	tokenArr := strings.Split(token, ".")
	if len(tokenArr) != 3 {
		log.Println("token token!")
		return *payloadModel
	}
	payload := tokenArr[1]

	payloadBytes, err := base64.URLEncoding.DecodeString(payload)

	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(payloadBytes, payloadModel)

	return *payloadModel
}

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
