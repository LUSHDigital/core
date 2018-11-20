package keys

var (
	// StagingTokenPublicKeySource represents the url to the JWT public key for the staging environment
	StagingTokenPublicKeySource = HTTPSource("https://api-gateway-staging.platformserviceaccount.com/token/public")

	// ProductionTokenPublicKeySource represents the url to the JWT public key for the production environment
	ProductionTokenPublicKeySource = HTTPSource("https://api-gateway.platformserviceaccount.com/token/public")

	// StagingTokenPublicKeySources represents all staging token public key sources
	StagingTokenPublicKeySources = Sources{
		StagingTokenPublicKeySource,
	}

	// ProductionTokenPublicKeySources represents all production token public key sources
	ProductionTokenPublicKeySources = Sources{
		ProductionTokenPublicKeySource,
	}
)
