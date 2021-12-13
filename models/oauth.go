package models

import (
	"database/sql"

	"github.com/RichardKnop/uuid"
	"gorm.io/gorm"

	"time"
)

// OauthClient ...
type OauthClient struct {
	gorm.Model  `json:"-"`
	Key         string         `sql:"type:varchar(254);unique;not null"`
	Secret      string         `sql:"type:varchar(60);not null"`
	RedirectURI sql.NullString `sql:"type:varchar(200)"`
}

// TableName specifies table name
func (c *OauthClient) TableName() string {
	return "oauth_clients"
}

// OauthUser ...
type OauthUser struct {
	gorm.Model  `json:"-"`
	UserName    string
	Phone       string `sql:"type:varchar(20);unique;not null"`
	CompanyName string
	Title       string
	//Age          sql.NullInt64
	//Birthday     *time.Time
	Email    string
	RoleID   uint           `json:"-"`
	Role     *OauthRole     `gorm:"foreignkey:RoleID"`
	Password sql.NullString `sql:"type:varchar(60)" json:"-"`
	Code     uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4()" json:"code"`
	//Address      string         `gorm:"index:addr"`     // 给address字段创建名为addr的索引
	//ImageLink    string         `json:"image_link"`
	Projects      []Project             `gorm:"many2many:project_users;" json:"projects"`
	Invites       []OauthInvitationCode `json:"invites"`
	MyOwnProjects []Project             `json:"my_own_projects"`
}

// TableName specifies table name
func (u *OauthUser) TableName() string {
	return "oauth_users"
}

// OauthPhoneNumValidate ...
type OauthPhoneNumValidate struct {
	gorm.Model `json:"-"`
	Phone      string `sql:"type:varchar(50);unique;not null"`
	Code       string
	ExpiresAt  time.Time `sql:"not null"`
}

// OauthScope ...
type OauthScope struct {
	gorm.Model  `json:"-"`
	Scope       string `sql:"type:varchar(200);unique;not null"`
	Description sql.NullString
	IsDefault   bool `sql:"default:false"`
}

// OauthRole ...
type OauthRole struct {
	gorm.Model `json:"-"`
	Name       string
}

// TableName specifies table name
func (r *OauthRole) TableName() string {
	return "oauth_roles"
}

type OauthInvitationCode struct {
	gorm.Model  `json:"-"`
	OauthUserID uint `json:"-"`
	Code        string
	Phone       string
	UserCode    string
	ProjectCode string
	RoleName    string
	Status      string
	ExpiresAt   time.Time
}

// TableName specifies table name
func (r *OauthInvitationCode) TableName() string {
	return "oauth_invitation_code"
}

// OauthRefreshToken ...
type OauthRefreshToken struct {
	gorm.Model `json:"-"`
	UserID     uint `sql:"index" json:"-"`
	User       *OauthUser
	Token      string    `sql:"type:varchar(40);unique;not null"`
	ExpiresAt  time.Time `sql:"not null"`
	Scope      string    `sql:"type:varchar(200);not null"`
}

// TableName specifies table name
func (rt *OauthRefreshToken) TableName() string {
	return "oauth_refresh_tokens"
}

// OauthAccessToken ...
type OauthAccessToken struct {
	gorm.Model `json:"-"`
	UserID     uint       `sql:"index;not null" json:"-"`
	User       *OauthUser `gorm:"foreignkey:UserID"`
	Token      string     `sql:"type:varchar(40);unique;not null"`
	ExpiresAt  time.Time  `sql:"not null"`
	Scope      string     `sql:"type:varchar(200);not null"`
}

// TableName specifies table name
func (at *OauthAccessToken) TableName() string {
	return "oauth_access_tokens"
}

// NewOauthRefreshToken creates new OauthRefreshToken instance
func NewOauthRefreshToken(user *OauthUser, expiresIn int, scope string) *OauthRefreshToken {
	refreshToken := &OauthRefreshToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		refreshToken.UserID = user.ID
	}
	return refreshToken
}

// NewOauthAccessToken creates new OauthAccessToken instance
func NewOauthAccessToken(user *OauthUser, expiresIn int, scope string) *OauthAccessToken {
	accessToken := &OauthAccessToken{
		Token:     uuid.New(),
		ExpiresAt: time.Now().UTC().Add(time.Duration(expiresIn) * time.Second),
		Scope:     scope,
	}
	if user != nil {
		accessToken.UserID = user.ID
	}
	return accessToken
}

// OauthAccessTokenPreload sets up Gorm preloads for an access token object
func OauthAccessTokenPreload(db *gorm.DB) *gorm.DB {
	return OauthAccessTokenPreloadWithPrefix(db, "")
}

// OauthAccessTokenPreloadWithPrefix sets up Gorm preloads for an access token object,
// and prefixes with prefix for nested objects
func OauthAccessTokenPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.
		Preload(prefix + "User")
}

// OauthRefreshTokenPreload sets up Gorm preloads for a refresh token object
func OauthRefreshTokenPreload(db *gorm.DB) *gorm.DB {
	return OauthRefreshTokenPreloadWithPrefix(db, "")
}

// OauthRefreshTokenPreloadWithPrefix sets up Gorm preloads for a refresh token object,
// and prefixes with prefix for nested objects
func OauthRefreshTokenPreloadWithPrefix(db *gorm.DB, prefix string) *gorm.DB {
	return db.
		Preload(prefix + "User")
}
