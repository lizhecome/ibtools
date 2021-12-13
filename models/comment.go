package models

type Comment struct {
	DDItemID uint `json:"-"`
	UserID   uint `json:"-"`
	User     OauthUser
	Text     string
}
