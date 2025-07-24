package utils

import (
	"crypto/rand"
	"crypto/rsa"
)

// GenerateRSAKeyPair generates an RSA key pair for JWT signing
func GenerateRSAKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// RSAPublicKeyToJWK converts an RSA public key to JWK format
// TODO: Implement proper RSA to JWK conversion
func RSAPublicKeyToJWK(pubKey *rsa.PublicKey) map[string]string {
	// Placeholder implementation
	return map[string]string{
		"kty": "RSA",
		"use": "sig",
		"kid": "mock-key-id",
		"n":   "mock-modulus",
		"e":   "AQAB",
	}
}
