package dbAto

import (
	"crypto/rsa"
	"log"
)

var PrivateKey *rsa.PrivateKey

func Init() {
	var err error
	PrivateKey, err = LoadKey()
	if err != nil {
		log.Println(err)
	}

}
