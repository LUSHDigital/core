# Auth
The `core/auth` package provides functions for services to issue and sign api consumer tokens.

## Using the issuer

### Start a new issuer
You can initiate an token issuer by passing a valid RSA or ECDSA PEM block.

```go
var private = []byte(`... private key ...`)
issuer := auth.NewIssuerFromPEM(private, jwt.SigningMethodRS256)
```

### Issue new tokens
A token can be issued with any struct that follows the [`jwt.Claims`](https://github.com/dgrijalva/jwt-go/blob/master/claims.go#L11-L13) interface.

```go
claims := jwt.StandardClaims{
	Id:        "1234",
	Issuer:    "Tests",
	Audience:  "Developers",
	Subject:   "Example",
	ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	IssuedAt:  time.Now().Unix(),
	NotBefore: time.Now().Unix(),
}
raw, err := issuer.Issue(&claims)
if err != nil {
	return
}
```

## Using the parser

### Start a new parser
You can initiate an token parser by passing a valid RSA or ECDSA PEM block.

```go
var public = []byte(`... public key ...`)
var fn func(pk crypto.PublicKey) jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return pk, fmt.Errorf("unknown algorithm: %v", token.Header["alg"])
		}
		return pk, nil
	}
}
parser := auth.NewParserFromPEM(public, fn)
```

### Parse existing tokens
Now you can parse any token that is signed with the public key provided to the parser.

```go
var claims jwt.StandardClaims
err := parser.Parse(`... jwt ...`, &claims)
if err != nil {
	return
}
```

## Mocking the issuer & parser
An issuer can be mocked with a temporary key pair for testing.

```go
issuer, parser, err := authmock.NewRSAIssuerAndParser()
if err != nil {
	log.Fatalln(err)
}
```