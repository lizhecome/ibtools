package oauth

import (
	"errors"
	"sort"
	"strings"

	"ibtools_server/models"
)

var (
	// ErrInvalidScope ...
	ErrInvalidScope = errors.New("invalid scope")
)

// GetScope takes a requested scope and, if it's empty, returns the default
// scope, if not empty, it validates the requested scope
func (s *Service) GetScope(requestedScope string) (string, error) {
	// Return the default scope if the requested scope is empty
	if requestedScope == "" {
		return s.GetDefaultScope(), nil
	}

	// If the requested scope exists in the database, return it
	if s.ScopeExists(requestedScope) {
		return requestedScope, nil
	}

	// Otherwise return error
	return "", ErrInvalidScope
}

// GetDefaultScope returns the default scope
func (s *Service) GetDefaultScope() string {
	// Fetch default scopes
	var scopes []string
	s.db.Model(new(models.OauthScope)).Where("is_default = ?", true).Pluck("scope", &scopes)

	// Sort the scopes alphabetically
	sort.Strings(scopes)

	// Return space delimited scope string
	return strings.Join(scopes, " ")
}

// ScopeExists checks if a scope exists
func (s *Service) ScopeExists(requestedScope string) bool {
	// Split the requested scope string
	scopes := strings.Split(requestedScope, " ")

	// Count how many of requested scopes exist in the database
	var count int64
	s.db.Model(new(models.OauthScope)).Where("scope in (?)", scopes).Count(&count)

	// Return true only if all requested scopes found
	return int(count) == len(scopes)
}
