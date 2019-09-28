package auth

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

// PrivateKeyFromPEM will take a private key PEM and derive the private key from it.
func PrivateKeyFromPEM(key []byte) (crypto.PrivateKey, error) {
	var err error
	var block *pem.Block
	var parsed interface{}
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}
	encodings := []encodingFunc{
		func(der []byte) (interface{}, error) { return x509.ParsePKCS1PrivateKey(der) },
		func(der []byte) (interface{}, error) { return x509.ParsePKCS8PrivateKey(der) },
		func(der []byte) (interface{}, error) { return x509.ParseECPrivateKey(der) },
	}
	for _, encoding := range encodings {
		parsed, err = encoding(block.Bytes)
		if err == nil {
			break
		}
	}
	var private crypto.PrivateKey
	switch key := parsed.(type) {
	case *rsa.PrivateKey:
		private = key
	case *ecdsa.PrivateKey:
		private = key
	default:
		return nil, ErrNotPrivateKey
	}
	return private, nil
}

// PrivateKeyFromPEMWithPassword will take a private key PEM with a password and derive the private key from it.
func PrivateKeyFromPEMWithPassword(key []byte, password string) (crypto.PrivateKey, error) {
	var err error
	var block *pem.Block
	var parsed interface{}
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}
	var decrypted []byte
	if decrypted, err = x509.DecryptPEMBlock(block, []byte(password)); err != nil {
		return nil, err
	}
	encodings := []encodingFunc{
		func(der []byte) (interface{}, error) { return x509.ParsePKCS1PrivateKey(der) },
		func(der []byte) (interface{}, error) { return x509.ParsePKCS8PrivateKey(der) },
	}
	for _, encoding := range encodings {
		parsed, err = encoding(decrypted)
		if err == nil {
			break
		}
	}
	var private crypto.PrivateKey
	switch key := parsed.(type) {
	case *rsa.PrivateKey:
		private = key
	default:
		return nil, ErrNotRSAPrivateKey
	}
	return private, nil
}

// PublicKeyFromPEM will take a public key PEM and derive the public key from it.
func PublicKeyFromPEM(key []byte) (crypto.PublicKey, error) {
	var err error
	var block *pem.Block
	var parsed interface{}
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}
	if parsed, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsed = cert.PublicKey
		} else {
			return nil, err
		}
	}
	var public crypto.PublicKey
	switch key := parsed.(type) {
	case *rsa.PublicKey:
		public = key
	case *ecdsa.PublicKey:
		public = key
	default:
		return nil, ErrNotPublicKey
	}
	return public, nil
}
