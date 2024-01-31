package models

import (
	"time"
)

type JwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}


type JwtPayload struct {
	UserId   string    `json:"uid"`
	Com      string    `json:"company"`
	ComId    int     `json:"company_id"`
	Dep      string    `json:"department"`
	Nam      string    `json:"name"`
	Address  string    `json:"add"`
	Lng      float64   `json:"lng"`
	Lat      float64   `json:"lat"`
	Iat      time.Time `json:"iat"`
	Exp      time.Time `json:"exp"`
	Iss      string    `json:"iss"`
	Aut      int       `json:"aut"`
	SiteType string    `json:"site_type"`
}

type JwtSignature struct {
	Signature []byte `json:"signature"`
}



