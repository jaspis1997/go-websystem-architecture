package entities

type TokenService interface {
	// Authorize a token with a matching token name
	// Returns true on successful authentication
	// errors
	//  - Token not found: not found token(%s=name)
	Authorize(name, token string) (ok bool, err error)

	// Retrieve the token to be sent externally from token name
	// errors
	//  - Token not found: not found token(%s=name)
	GetToken(name string) (token string, err error)
}

type BatchService interface {
	Start() error
	Stop() error
}
