package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/SermoDigital/jose/crypto"
	"github.com/SermoDigital/jose/jws"
)

var hashAlg = crypto.SigningMethodHS256

const scopesClaim = "scopes"

type Scope string

var (
	SkinScope = Scope("skin")
)

type JwtAuth struct {
	Key []byte
}

func (t *JwtAuth) NewToken(scopes ...Scope) ([]byte, error) {
	if len(t.Key) == 0 {
		return nil, errors.New("signing key not available")
	}

	claims := jws.Claims{}
	claims.Set(scopesClaim, scopes)
	claims.SetIssuedAt(time.Now())
	encoder := jws.NewJWT(claims, hashAlg)
	token, err := encoder.Serialize(t.Key)
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (t *JwtAuth) Check(req *http.Request) error {
	if len(t.Key) == 0 {
		return &Unauthorized{"Signing key not set"}
	}

	bearerToken := req.Header.Get("Authorization")
	if bearerToken == "" {
		return &Unauthorized{"Authentication header not presented"}
	}

	if !strings.EqualFold(bearerToken[0:7], "BEARER ") {
		return &Unauthorized{"Cannot recognize JWT token in passed value"}
	}

	tokenStr := bearerToken[7:]
	token, err := jws.ParseJWT([]byte(tokenStr))
	if err != nil {
		return &Unauthorized{"Cannot parse passed JWT token"}
	}

	err = token.Validate(t.Key, hashAlg)
	if err != nil {
		return &Unauthorized{"JWT token have invalid signature. It may be corrupted or expired."}
	}

	return nil
}

type Unauthorized struct {
	Reason string
}

func (e *Unauthorized) Error() string {
	if e.Reason != "" {
		return e.Reason
	}

	return "Unauthorized"
}
