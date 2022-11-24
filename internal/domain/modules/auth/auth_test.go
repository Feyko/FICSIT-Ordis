package auth

import (
	"FICSIT-Ordis/internal/domain/domain"
	"FICSIT-Ordis/test"
	"context"
	"github.com/stretchr/testify/suite"
	"golang.org/x/exp/slices"
	"testing"
)

var adminToken string

func TestAuthModuleTestSuite(t *testing.T) {
	suite.Run(t, new(AuthModuleTestSuite))
}

type AuthModuleTestSuite struct {
	suite.Suite
	mod          *Module
	tokenToClean Token
}

func (s *AuthModuleTestSuite) TearDownTest() {
	if s.tokenToClean.TokenID != "" {
		err := s.mod.RevokeTokenIDNoAuth(s.tokenToClean.TokenID)
		s.tokenToClean = Token{}
		s.Require().NoError(err)
	}
}

func (s *AuthModuleTestSuite) SetupTest() {
	rep, err := test.GetRepo()
	s.Require().NoError(err)

	mod, err := New(Config{Secret: "test-secret"}, rep)
	s.Require().NoError(err)
	token, err := mod.NewTokenNoAuth(domain.RoleAdmin.ID)
	s.Require().NoError(err)
	adminToken = token.String

	s.mod = mod
}

func (s *AuthModuleTestSuite) TestNewToken() {
	token, err := s.mod.NewTokenNoAuth()
	s.tokenToClean = token
	s.Require().NoError(err)

}

func (s *AuthModuleTestSuite) TestNewTokenIsValid() {
	token, err := s.mod.NewTokenNoAuth()
	s.tokenToClean = token
	s.Require().NoError(err)
	_, _, err = s.mod.ValidateTokenString(token.String)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestInvalidToken() {
	_, valid, err := s.mod.ValidateTokenString("")
	s.Require().NoError(err)
	s.Require().False(valid)
}

func (s *AuthModuleTestSuite) TestNewTokenWithRolesIsValid() {
	token, err := s.mod.NewTokenNoAuth(domain.RoleAdmin.ID, domain.RoleModerator.ID)
	s.tokenToClean = token
	s.Require().NoError(err)
	_, _, err = s.mod.ValidateTokenString(token.String)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestNewTokenWithRolesHasRoles() {
	roleIDs := []int{domain.RoleAdmin.ID, domain.RoleModerator.ID}
	roles := []domain.Role{domain.Roles[roleIDs[0]], domain.Roles[roleIDs[1]]}
	token, err := s.mod.NewTokenNoAuth(roleIDs...)
	s.tokenToClean = token
	s.Require().NoError(err)
	s.Equal(roles, token.Roles)
}

func (s *AuthModuleTestSuite) TestNewTokenWithRolesHasRolesValidated() {
	roleIDs := []int{domain.RoleAdmin.ID, domain.RoleModerator.ID}
	roles := []domain.Role{domain.Roles[roleIDs[0]], domain.Roles[roleIDs[1]]}
	token, err := s.mod.NewTokenNoAuth(roleIDs...)
	s.tokenToClean = token
	s.Require().NoError(err)
	_, _, err = s.mod.ValidateTokenString(token.String)
	s.Require().NoError(err)
	s.Equal(roles, token.Roles)
}

func (s *AuthModuleTestSuite) TestTokenCorrectPermissionsOneRole() {
	perms := domain.RoleAdmin.Permissions
	token, err := s.mod.NewTokenNoAuth(domain.RoleAdmin.ID)
	s.tokenToClean = token
	s.Require().NoError(err)
	s.Equal(perms, token.Permissions)
}

func (s *AuthModuleTestSuite) TestTokenCorrectPermissionsMultipleRoles() {
	allPerms := []domain.Permission{domain.PermissionContentEditing, domain.PermissionTicketManagement}
	domain.Roles[len(domain.Roles)] = domain.Role{ID: len(domain.Roles), Permissions: []domain.Permission{domain.PermissionContentEditing}}
	domain.Roles[len(domain.Roles)] = domain.Role{ID: len(domain.Roles), Permissions: []domain.Permission{domain.PermissionTicketManagement}}
	defer func() {
		delete(domain.Roles, len(domain.Roles)-1)
		delete(domain.Roles, len(domain.Roles)-1)
	}()

	token, err := s.mod.NewTokenNoAuth(len(domain.Roles)-1, len(domain.Roles)-2)
	s.tokenToClean = token
	s.Require().NoError(err)
	slices.Sort(token.Permissions)
	slices.Sort(allPerms)
	s.Equal(allPerms, token.Permissions)
}

func (s *AuthModuleTestSuite) TestTokenHasPermission() {
	domain.Roles[len(domain.Roles)] = domain.Role{ID: len(domain.Roles), Permissions: []domain.Permission{domain.PermissionTicketManagement}}
	defer func() {
		delete(domain.Roles, len(domain.Roles)-1)
	}()
	token, err := s.mod.NewTokenNoAuth(len(domain.Roles) - 1)
	s.tokenToClean = token
	s.Require().NoError(err)
	s.Require().True(token.HasPermissions(domain.PermissionTicketManagement))
	s.Require().False(token.HasPermissions(domain.PermissionContentEditing))
}

func (s *AuthModuleTestSuite) TestAuthorize() {
	ctx := context.WithValue(context.Background(), "ordis-string-token", adminToken)
	err := s.mod.Authorize(&ctx)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestAuthorizePermissions() {
	ctx := context.WithValue(context.Background(), "ordis-string-token", adminToken)
	err := s.mod.Authorize(&ctx, domain.PermissionTokenCreation)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestAuthorizeWrongPermissions() {
	ctx := context.WithValue(context.Background(), "ordis-string-token", adminToken)
	err := s.mod.Authorize(&ctx, "NonexistentPermission")
	s.Require().Error(err)
}

func (s *AuthModuleTestSuite) TestReauthorize() {
	ctx := context.WithValue(context.Background(), "ordis-string-token", adminToken)
	err := s.mod.Authorize(&ctx, domain.PermissionTokenCreation)
	s.Require().NoError(err)
	err = s.mod.Authorize(&ctx, domain.PermissionTokenCreation)
	s.Require().NoError(err)
}

func (s *AuthModuleTestSuite) TestReauthorizeWrongPermissions() {
	ctx := context.WithValue(context.Background(), "ordis-string-token", adminToken)
	err := s.mod.Authorize(&ctx, "NonexistentPermission")
	s.Require().Error(err)
	err = s.mod.Authorize(&ctx, "NonexistentPermission")
	s.Require().Error(err)
}
