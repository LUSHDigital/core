package keybroker

var (
	// JWTPublicKeyEnvStringSource represents the source of an RSA public key as a string
	JWTPublicKeyEnvStringSource = EnvStringSource("JWT_PUBLIC_KEY")

	// JWTPublicKeyEnvHTTPSource represents the source of an RSA public key at a HTTP GET destination
	JWTPublicKeyEnvHTTPSource = EnvHTTPSource("JWT_PUBLIC_KEY_URL")

	// JWTPublicKeyEnvFileSource represents the source of an RSA public key on disk
	JWTPublicKeyEnvFileSource = EnvFileSource("JWT_PUBLIC_KEY_PATH")

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
