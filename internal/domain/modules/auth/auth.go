package auth

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/domain"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
)

func New(conf config.AuthConfig) (*Module, error) {
	return &Module{
		secret: []byte(conf.Secret),
	}, nil
}

type Module struct {
	secret []byte
}

func (m *Module) NewToken(roles ...domain.Role) (Token, error) {
	claims := &authClaims{roles}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return Token{}, errors.Wrap(err, "could not sign token")
	}
	token.Raw = tokenString
	return newToken(token, claims), nil
}

func (m *Module) ValidateToken(token *Token) error {
	claims := &authClaims{}
	jwToken, err := jwt.ParseWithClaims(token.String, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("Unexpected signing method '%v'", token.Header["alg"])
		}
		return m.secret, nil
	})
	token.token = jwToken
	if claims != nil {
		token.Roles = claims.Roles
	}
	return err
}

type Token struct {
	String      string
	Roles       []domain.Role
	Permissions []domain.Permission

	token *jwt.Token
}

func newToken(token *jwt.Token, claims *authClaims) Token {
	return Token{
		String:      token.Raw,
		Roles:       claims.Roles,
		Permissions: PermissionsFromRoles(claims.Roles),
		token:       token,
	}
}

type authClaims struct {
	Roles []domain.Role
}

func (c *authClaims) Valid() error {
	return nil
}

func PermissionsFromRoles(roles []domain.Role) []domain.Permission {
	perms := []domain.Permission{}
	for _, role := range roles {
		for _, perm := range role.Permissions {
			if !slices.Contains(perms, perm) {
				perms = append(perms, perm)
			}
		}
	}
	return perms
}
