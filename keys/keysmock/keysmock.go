package keysmock

import "crypto/rsa"

// MockRSAPublicKey resolves any source and returns a mocked RSAPublicKeyCopier and Renewer
func MockRSAPublicKey(key rsa.PublicKey) *RSAPublicKeyMock {
	return &RSAPublicKeyMock{
		key: &key,
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
