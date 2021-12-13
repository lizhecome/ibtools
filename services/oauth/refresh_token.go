package oauth

import (
	"errors"
	"time"

	"ibtools_server/models"
	"ibtools_server/util"

	"gorm.io/gorm"
)

var (
	// ErrRefreshTokenNotFound ...
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	// ErrRefreshTokenExpired ...
	ErrRefreshTokenExpired = errors.New("refresh token expired")
	// ErrRequestedScopeCannotBeGreater ...
	ErrRequestedScopeCannotBeGreater = errors.New("requested scope cannot be greater")
)

// GetOrCreateRefreshToken retrieves an existing refresh token, if expired,
// the token gets deleted and new refresh token is created
func (s *Service) GetOrCreateRefreshToken(user *models.OauthUser, expiresIn int, scope string) (*models.OauthRefreshToken, error) {
	// Try to fetch an existing refresh token first
	refreshToken := new(models.OauthRefreshToken)
	query := models.OauthRefreshTokenPreload(s.db)
	if user != nil && user.ID > 0 {
		query = query.Where("user_id = ?", user.ID)
	} else {
		query = query.Where("user_id IS NULL")
	}
	err := query.First(refreshToken).Error

	// Check if the token is expired, if found
	var expired bool
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		expired = time.Now().UTC().After(refreshToken.ExpiresAt)
	}

	// If the refresh token has expired, delete it
	if expired {
		s.db.Unscoped().Delete(refreshToken)
	}

	// Create a new refresh token if it expired or was not found
	if expired || errors.Is(err, gorm.ErrRecordNotFound) {
		refreshToken = models.NewOauthRefreshToken(user, expiresIn, scope)
		if err := s.db.Create(refreshToken).Error; err != nil {
			return nil, err
		}
		refreshToken.User = user
	}

	return refreshToken, nil
}

// GetValidRefreshToken returns a valid non expired refresh token
func (s *Service) GetValidRefreshToken(token string) (*models.OauthRefreshToken, error) {
	// Fetch the refresh token from the database
	refreshToken := new(models.OauthRefreshToken)
	err := models.OauthRefreshTokenPreload(s.db).
		Where("token = ?", token).First(refreshToken).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRefreshTokenNotFound
	}

	// Check the refresh token hasn't expired
	if time.Now().UTC().After(refreshToken.ExpiresAt) {
		return nil, ErrRefreshTokenExpired
	}

	return refreshToken, nil
}

// getRefreshTokenScope returns scope for a new refresh token
func (s *Service) getRefreshTokenScope(refreshToken *models.OauthRefreshToken, requestedScope string) (string, error) {
	var (
		scope = refreshToken.Scope // default to the scope originally granted by the resource owner
		err   error
	)

	// If the scope is specified in the request, get the scope string
	if requestedScope != "" {
		scope, err = s.GetScope(requestedScope)
		if err != nil {
			return "", err
		}
	}

	// Requested scope CANNOT include any scope not originally granted
	if !util.SpaceDelimitedStringNotGreater(scope, refreshToken.Scope) {
		return "", ErrRequestedScopeCannotBeGreater
	}

	return scope, nil
}
