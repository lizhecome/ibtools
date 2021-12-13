package oauth

import (
	"errors"

	"ibtools_server/models"

	"gorm.io/gorm"
)

var (
	// ErrRoleNotFound ...
	ErrRoleNotFound = errors.New("role not found")
)

// FindRoleByID looks up a role by ID and returns it
func (s *Service) FindRoleByID(id uint) (*models.OauthRole, error) {
	role := new(models.OauthRole)
	err := s.db.Where("id = ?", id).First(role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRoleNotFound
	}
	return role, nil
}

// FindRoleByID looks up a role by ID and returns it
func (s *Service) FindRoleByName(roleName string) (*models.OauthRole, error) {
	role := new(models.OauthRole)
	err := s.db.Where("name = ?", roleName).First(role).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRoleNotFound
	}
	return role, nil
}
