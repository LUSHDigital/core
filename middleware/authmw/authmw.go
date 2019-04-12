package authmw

import (
	"crypto/rsa"
)

// RSAPublicKeyCopierRenewer represents the combination of a Copier and Renewer interface
type RSAPublicKeyCopierRenewer interface {
	Copy() rsa.PublicKey
	Renew()
}
