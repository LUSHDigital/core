// Package keybrokermock implements no-op mocks for the keys package
package keybrokermock

import (
	"crypto/ecdsa"
	"crypto/rsa"
)

// MockRSAPublicKey resolves any source and returns a mocked RSAPublicKey Copier and Renewer
func MockRSAPublicKey(key *rsa.PublicKey) *RSAPublicKeyMock {
	return &RSAPublicKeyMock{
		key: key,
	}
}

// RSAPublicKeyMock defines the implementation for brokering an RSA public key during testing
type RSAPublicKeyMock struct {
	key *rsa.PublicKey
}

// Copy returns a shallow copy o the RSA public key
func (b *RSAPublicKeyMock) Copy() rsa.PublicKey {
	return *b.key
}

// Renew is a no-op
func (b *RSAPublicKeyMock) Renew() {
	// no-op
}

// Close is a no-op
func (b *RSAPublicKeyMock) Close() {
	// no-op
}

// MockECDSAPublicKey resolves any source and returns a mocked ECDSAPublicKey Copier and Renewer
func MockECDSAPublicKey(key *ecdsa.PublicKey) *ECDSAPublicKeyMock {
	return &ECDSAPublicKeyMock{
		key: key,
	}
}

// ECDSAPublicKeyMock defines the implementation for brokering an ECDSA public key during testing
type ECDSAPublicKeyMock struct {
	key *ecdsa.PublicKey
}

// Copy returns a shallow copy o the ECDSA public key
func (b *ECDSAPublicKeyMock) Copy() ecdsa.PublicKey {
	return *b.key
}

// Renew is a no-op
func (b *ECDSAPublicKeyMock) Renew() {
	// no-op
}

// Close is a no-op
func (b *ECDSAPublicKeyMock) Close() {
	// no-op
}
