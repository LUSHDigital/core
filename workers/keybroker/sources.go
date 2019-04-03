package keybroker

import (
	"os"
)

var (
	// JWTPublicKeyEnvStringSource represents the source of an RSA public key as a string
	JWTPublicKeyEnvStringSource = StringSource(os.Getenv("JWT_PUBLIC_KEY"))

	// JWTPublicKeyEnvHTTPSource represents the source of an RSA public key at a HTTP GET destination
	JWTPublicKeyEnvHTTPSource = HTTPSource(os.Getenv("JWT_PUBLIC_KEY_URL"))

	// JWTPublicKeyEnvFileSource represents the source of an RSA public key on disk
	JWTPublicKeyEnvFileSource = FileSource(os.Getenv("JWT_PUBLIC_KEY_PATH"))

	// JWTPublicKeyDefaultFileSource represents the source of an RSA public key on disk
	JWTPublicKeyDefaultFileSource = FileSource("/usr/local/var/jwt.pub")

	// JWTPublicKeySources represents a chain of sources for JWT Public Keys in order of priority
	JWTPublicKeySources = Sources{
		JWTPublicKeyEnvStringSource,
		JWTPublicKeyEnvFileSource,
		JWTPublicKeyEnvHTTPSource,
		JWTPublicKeyDefaultFileSource,
	}
)
