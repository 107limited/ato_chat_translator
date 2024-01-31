package dbAto

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// OpenDB membuka koneksi ke database
func OpenDB() (*sql.DB, error) {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }

    dbURL := os.Getenv("DB_URL")

    db, err := sql.Open("mysql", dbURL)
    if err != nil {
        return nil, err
    }

    return db, nil
}

func LoadKey() (key *rsa.PrivateKey, err error) {
    db, err := OpenDB()
    if err != nil {
        log.Println(err)
        return nil, err
    }
    defer db.Close()

    q := "select private_key from sec_m"
    var keyBytes []byte
    rows, err := db.Query(q)
    if err != nil {
        log.Println(err)
        return nil, err
    }

    if rows.Next() {
        err2 := rows.Scan(&keyBytes)
        if err2 != nil {
            log.Println(err2)
        }
        log.Println("Load success")
        rows.Close()
        return x509.ParsePKCS1PrivateKey(keyBytes)
    }
    rows.Close()
    return &rsa.PrivateKey{}, err
}
