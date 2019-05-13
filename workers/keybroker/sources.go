package keybroker

var (
	// JWTPublicKeyEnvStringSource represents the source of an RSA public key as a string
	JWTPublicKeyEnvStringSource = EnvStringSource("JWT_PUBLIC_KEY")

	// JWTPublicKeyEnvHTTPSource represents the source of an RSA public key at a HTTP GET destination
	JWTPublicKeyEnvHTTPSource = EnvHTTPSource("JWT_PUBLIC_KEY_URL")

	// JWTPublicKeyEnvFileSource represents the source of an RSA public key on disk
	JWTPublicKeyEnvFileSource = EnvFileSource("JWT_PUBLIC_KEY_PATH")

	// JWTPublicKeyDefaultFileSource represents the source of an RSA public key on disk
	JWTPublicKeyDefaultFileSource = FileSource("/usr/local/var/jwt.pub.pem")

	// JWTPublicKeySources represents a chain of sources for JWT Public Keys in order of priority
	JWTPublicKeySources = Sources{
		JWTPublicKeyEnvStringSource,
		JWTPublicKeyEnvFileSource,
		JWTPublicKeyEnvHTTPSource,
		JWTPublicKeyDefaultFileSource,
	}

	// JWTPrivateKeyEnvStringSource represents the source of an RSA public key as a string
	JWTPrivateKeyEnvStringSource = EnvStringSource("JWT_PRIVATE_KEY")

	// JWTPrivateKeyEnvHTTPSource represents the source of an RSA public key at a HTTP GET destination
	JWTPrivateKeyEnvHTTPSource = EnvHTTPSource("JWT_PRIVATE_KEY_URL")

	// JWTPrivateKeyEnvFileSource represents the source of an RSA public key on disk
	JWTPrivateKeyEnvFileSource = EnvFileSource("JWT_PRIVATE_KEY_PATH")

	// JWTPrivateKeyDefaultFileSource represents the source of an RSA public key on disk
	JWTPrivateKeyDefaultFileSource = FileSource("/usr/local/var/jwt.pem")

	// JWTPrivateKeySources represents a chain of sources for JWT Public Keys in order of priority
	JWTPrivateKeySources = Sources{
		JWTPrivateKeyEnvStringSource,
		JWTPrivateKeyEnvFileSource,
		JWTPrivateKeyEnvHTTPSource,
		JWTPrivateKeyDefaultFileSource,
	}
)
