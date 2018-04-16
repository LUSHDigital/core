package errors

// LoginUnauthorisedError - Error to throw when a login was unauthorised.
type LoginUnauthorisedError struct{}

// Error - Error string for login not unauthorised.
func (e LoginUnauthorisedError) Error() string {
	return "unauthorised"
}

// ConsumerHasNoTokensError - Error to throw when a consumer has no tokens.
type ConsumerHasNoTokensError struct{}

// Error - Error string when a consumer has no tokens.
func (e ConsumerHasNoTokensError) Error() string {
	return "consumer has no tokens"
}
