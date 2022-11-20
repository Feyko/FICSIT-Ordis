package auth

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/internal/id"
	"FICSIT-Ordis/internal/ports/repos"
	"FICSIT-Ordis/internal/ports/repos/repo"
	"FICSIT-Ordis/internal/util"
	"context"
	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"
	"net/http"
	"strings"
)

type AuthedConfig struct {
	NoAuth     bool
	AuthModule *Module
}

type Config struct {
	Secret string
}

func New[T id.IDer](conf Config, repository repo.Repository[T]) (*Module, error) {
	collection, err := repos.GetOrCreateCollection[id.Wrapper[nothing]](repository, "Auth")
	if err != nil {
		return nil, errors.Wrap(err, "could not get or create the collection")
	}

	return &Module{
		secret: []byte(conf.Secret),
		coll:   collection,
	}, nil
}

type nothing struct{}

type Module struct {
	secret []byte
	coll   repo.Collection[id.Wrapper[nothing]]
}

func (m *Module) Middleware() func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
			r = r.WithContext(CtxWithTokenString(r.Context(), tokenString))
			next.ServeHTTP(w, r)
		})
	}
}

func CtxWithTokenString(ctx context.Context, token string) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, "ordis-string-token", token)
}

// Authorize requires a valid JWT token passed through the context. Use CtxWithTokenString to pass a token.
// Use Middleware for ease of use with http.
func (m *Module) Authorize(ctx *context.Context, perms ...domain.Permission) error {
	if ctx == nil || *ctx == nil {
		return errors.New("nil context")
	}

	userToken := (*ctx).Value("ordis-user-token")
	token, ok := userToken.(*Token)
	if ok && token.HasPermissions(perms...) {
		return nil
	}

	authorization := (*ctx).Value("ordis-string-token")
	tokenString, ok := authorization.(string)
	if !ok {
		return errors.New("missing token string")
	}

	token, valid, err := m.ValidateTokenString(tokenString)
	if err != nil {
		return errors.Wrap(err, "error validating token")
	}

	if !valid {
		return errors.New("invalid token")
	}

	*ctx = context.WithValue(*ctx, "ordis-user-token", token)
	if !token.HasPermissions(perms...) {
		return errors.New("insufficient permissions")
	}

	return nil
}

func (m *Module) NewToken(ctx context.Context, roleIDs ...int) (Token, error) {
	err := m.Authorize(&ctx, domain.PermissionTokenCreation)
	if err != nil {
		return Token{}, errors.Wrap(err, "error authorizing")
	}

	return m.NewTokenNoAuth(roleIDs...)
}

func (m *Module) NewTokenNoAuth(roleIDs ...int) (Token, error) {
	tokenID := util.GenerateUniqueID()

	claims := &authClaims{roleIDs, tokenID}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString(m.secret)
	if err != nil {
		return Token{}, errors.Wrap(err, "could not sign token")
	}

	err = m.coll.Create(nil, id.Wrap(nothing{}, tokenID))
	if err != nil {
		return Token{}, errors.Wrap(err, "error registering token id")
	}

	token.Raw = tokenString
	return newToken(token, claims), nil
}

func (m *Module) ValidateTokenString(tokenString string) (*Token, bool, error) {
	claims := &authClaims{}
	jwToken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.Errorf("Unexpected signing method '%v'", token.Header["alg"])
		}
		return m.secret, nil
	})

	var token Token

	token.token = jwToken

	if claims != nil {
		token.Roles = rolesFromIDs(claims.RoleIDs)
		token.Permissions = PermissionsFromRoles(token.Roles)
		token.TokenID = claims.ID
	}

	_, err = m.coll.Get(nil, token.TokenID)
	if errors.Is(err, repo.ErrElementNotFound) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, errors.Wrap(err, "error verifying token id")
	}

	return &token, true, nil
}

func (m *Module) RevokeCurrentToken(ctx context.Context) error {
	err := m.Authorize(&ctx)
	if err != nil {
		return errors.Wrap(err, "error authorizing")
	}

	userToken := ctx.Value("ordis-user-token")
	token, ok := userToken.(*Token)
	if !ok {
		return errors.New("token not resolved")
	}

	return m.RevokeTokenNoAuth(token)
}

func (m *Module) RevokeToken(ctx context.Context, token *Token) error {
	return m.RevokeTokenID(ctx, token.TokenID)
}

func (m *Module) RevokeTokenNoAuth(token *Token) error {
	return m.RevokeTokenIDNoAuth(token.TokenID)
}

func (m *Module) RevokeTokenID(ctx context.Context, id string) error {
	err := m.Authorize(&ctx, domain.PermissionTokenCreation)
	if err != nil {
		return errors.Wrap(err, "error authorizing")
	}
	return m.RevokeTokenIDNoAuth(id)
}

func (m *Module) RevokeTokenIDNoAuth(id string) error {
	err := m.coll.Delete(nil, id)
	if errors.Is(err, repo.ErrElementNotFound) {
		return errors.New("token id is not valid")
	}
	if err != nil {
		return errors.Wrap(err, "error deleting id")
	}

	return nil
}

func newToken(token *jwt.Token, claims *authClaims) Token {
	roles := rolesFromIDs(claims.RoleIDs)
	return Token{
		String:      token.Raw,
		Roles:       roles,
		Permissions: PermissionsFromRoles(roles),
		TokenID:     claims.ID,
		token:       token,
	}
}

type Token struct {
	TokenID     string
	String      string
	Roles       []domain.Role
	Permissions []domain.Permission

	token *jwt.Token
}

func (t Token) ID() string {
	return t.TokenID
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
	ID      string
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
