package oauth

import (
	"errors"
	"time"

	"ibtools_server/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	// "drpiggy_server/session"
)

var (
	// ErrAccessTokenNotFound ...
	ErrAccessTokenNotFound = errors.New("access token not found")
	// ErrAccessTokenExpired ...
	ErrAccessTokenExpired = errors.New("access token expired")
)

// Authenticate checks the access token is valid
func (s *Service) Authenticate(token string) (*models.OauthAccessToken, error) {
	// Fetch the access token from the database
	accessToken := new(models.OauthAccessToken)
	err := s.db.Preload(clause.Associations).Where("token = ?", token).First(accessToken).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAccessTokenNotFound
	}

	// Check the access token hasn't expired
	if time.Now().UTC().After(accessToken.ExpiresAt) {
		return nil, ErrAccessTokenExpired
	}

	// Extend refresh token expiration database
	query := s.db.Model(new(models.OauthRefreshToken))
	if accessToken.UserID > 0 {
		query = query.Where("user_id = ?", accessToken.UserID)
	} else {
		query = query.Where("user_id IS NULL")
	}
	increasedExpiresAt := s.db.NowFunc().Add(
		time.Duration(s.cnf.Oauth.RefreshTokenLifetime) * time.Second,
	)
	if err := query.UpdateColumn("expires_at", increasedExpiresAt).Error; err != nil {
		return nil, err
	}

	return accessToken, nil
}

// ClearUserTokens deletes the user's access and refresh tokens associated with this client id
func (s *Service) ClearUserTokens(userrefreshToken, useraccessToken string) {
	// Clear all refresh tokens with user_id and client_id
	refreshToken := new(models.OauthRefreshToken)
	err := models.OauthRefreshTokenPreload(s.db).Where("token = ?", userrefreshToken).First(refreshToken).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.db.Unscoped().Where("user_id = ?", refreshToken.UserID).Delete(models.OauthRefreshToken{})
	}

	// Clear all access tokens with user_id and client_id
	accessToken := new(models.OauthAccessToken)
	err = models.OauthAccessTokenPreload(s.db).Where("token = ?", useraccessToken).First(accessToken).Error
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		s.db.Unscoped().Where("user_id = ?", accessToken.UserID).Delete(models.OauthAccessToken{})
	}
}
