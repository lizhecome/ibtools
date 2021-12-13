package oauth

import (
	"errors"
	"ibtools_server/models"
)

var (
	// ErrInvalidUsernameOrPassword ...
	ErrInvalidUsernameOrPassword = errors.New("invalid username or password")
)

//Login 登陆
func (s *Service) Login(user *models.OauthUser, scope string) (*models.OauthAccessToken, *models.OauthRefreshToken, error) {
	// Return error if user's role is not allowed to use this service
	// if !s.IsRoleAllowed(user.Role.Name) {
	// 	// For security reasons, return a general error message
	// 	return nil, nil, ErrInvalidUsernameOrPassword
	// }

	// Create a new access token
	accessToken, err := s.GrantAccessToken(
		user,
		s.cnf.Oauth.AccessTokenLifetime, // expires in
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	// Create or retrieve a refresh token
	refreshToken, err := s.GetOrCreateRefreshToken(
		user,
		s.cnf.Oauth.RefreshTokenLifetime, // expires in
		scope,
	)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}
