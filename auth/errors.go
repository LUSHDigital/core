package auth

import "errors"

var (
	// ErrKeyMustBePEMEncoded happens when the PEM format is not valid.
	ErrKeyMustBePEMEncoded = errors.New("invalid key: must be PEM encoded PKCS1 or PKCS8 private key")
	// ErrNotRSAPrivateKey happens when the key is not a valid RSA private key.
	ErrNotRSAPrivateKey = errors.New("invalid key: must be a valid RSA private key")
	// ErrNotPrivateKey happens when the key is neither an RSA or ECDSA private key.
	ErrNotPrivateKey = errors.New("invalid key: must be either an RSA or ECDSA private key.")
	// ErrNotPublicKey happens when the key is neither an RSA or ECDSA public key.
	ErrNotPublicKey = errors.New("invalid key: must be either an RSA or ECDSA public key.")
)
