package auth

import (
	"FICSIT-Ordis/internal/config"
	"FICSIT-Ordis/internal/domain/domain"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"strings"
)

func New(conf config.AuthConfig) (*Module, error) {
	return &Module{
		secret: []byte(conf.Secret),
	}, nil
}

type Module struct {
	secret []byte
}

func (m *Module) Authorize(ctx *context.Context, perms ...domain.Permission) error {
	if ctx == nil || *ctx == nil {
		return errors.New("nil context")
	}

	userToken := (*ctx).Value("ordis-user-token")
	token, ok := userToken.(*Token)
	if ok && token.HasPermissions(perms...) {
		return nil
	}

	authorization := (*ctx).Value("Authorization")
	tokenString, ok := authorization.(string)
	if !ok {
		return errors.New("invalid or missing Authorization header")
	}
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")

	token = &Token{String: tokenString}
	err := m.ValidateToken(token)
	if err != nil {
		return errors.Wrap(err, "invalid token")
	}
	*ctx = context.WithValue(*ctx, "ordis-user-token", token)
	if !token.HasPermissions(perms...) {
		return errors.New("insufficient permissions")
	}

	return nil
}

func (m *Module) NewToken(roles ...domain.Role) (Token, error) {
	claims := &authClaims{extractRoleIDs(roles)}
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
		token.Roles = rolesFromIDs(claims.RoleIDs)
		token.Permissions = PermissionsFromRoles(token.Roles)
	}
	return err
}

func newToken(token *jwt.Token, claims *authClaims) Token {
	roles := rolesFromIDs(claims.RoleIDs)
	return Token{
		String:      token.Raw,
		Roles:       roles,
		Permissions: PermissionsFromRoles(roles),
		token:       token,
	}
}

type Token struct {
	String      string
	Roles       []domain.Role
	Permissions []domain.Permission

	token *jwt.Token
}

func (t *Token) HasPermissions(perms ...domain.Permission) bool {
	for _, perm := range perms {
		if !slices.Contains(t.Permissions, perm) {
			return false
		}
	}
	return true
}

type authClaims struct {
	RoleIDs []int
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

func extractRoleIDs(roles []domain.Role) []int {
	ids := make([]int, len(roles))
	for i, role := range roles {
		ids[i] = role.ID
	}
	return ids
}

func rolesFromIDs(ids []int) []domain.Role {
	roles := make([]domain.Role, 0, len(ids))
	for _, id := range ids {
		role, ok := domain.Roles[id]
		if ok {
			roles = append(roles, role)
		}
	}
	return roles
}
