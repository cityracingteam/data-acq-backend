package jwt

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"log"

	"github.com/cityracingteam/data-acq-backend/models"
	"gorm.io/gorm"
)

var (
	keys []*models.JwtKey
)

func Init(db *gorm.DB) {
	// Init an empty array of pointers to jwtkey
	keys = []*models.JwtKey{}

	// Populate the array from the database
	result := db.Find(&keys)

	if result.Error != nil {
		log.Fatalln("[jwt/Init]: Error fetching jwt keys from database.")
	}

	// Check if we have at least one existing key
	if !(len(keys) > 0) {
		// We have less than one key (not greater than 0)
		key := createKey()

		result := db.Create(&key)

		if result.Error != nil {
			log.Fatalln("[jwt/Init]: Error adding new key to database")
		}

		// Add the new key to the array
		keys = append(keys, &key)

		// Now we have at least one existing key.
	}
}

func createKey() models.JwtKey {
	privKey, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)

	if err != nil {
		log.Fatalln("[jwt/createKey]: Failed to create new jwt key")
	}

	// Created sucessfully, convert the private key to ASN.1 (like in a certificate)
	bytes, marshalErr := x509.MarshalECPrivateKey(privKey)

	if marshalErr != nil {
		log.Fatalln("[jwt/createKey]: failed to marshal generated private key")
	}

	key := models.JwtKey{Data: bytes}

	return key

}
