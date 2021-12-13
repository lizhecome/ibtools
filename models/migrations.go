package models

import (
	"ibtools_server/util/migrations"

	"gorm.io/gorm"
)

var (
	list = []migrations.MigrationStage{
		{
			Name:     "001",
			Function: migrate,
		},
	}
)

// MigrateAll executes all migrations
func MigrateAll(db *gorm.DB) error {
	return migrations.Migrate(db, list)
}

//create extension "uuid-ossp"
func migrate(db *gorm.DB, name string) error {
	//-------------
	// OAUTH models
	//-------------
	err := db.SetupJoinTable(&OauthUser{}, "Projects", &ProjectUser{})
	db.SetupJoinTable(&Project{}, "Users", &ProjectUser{})
	if err != nil {
		return nil
	}
	db.AutoMigrate(
		OauthRole{},
		&OauthScope{},
		&OauthUser{},
		&OauthRefreshToken{},
		&OauthAccessToken{},
		&OauthPhoneNumValidate{},
		&OauthInvitationCode{},
		&Comment{},
		&Project{},
		&DDModel{},
		&DDItem{},
		&DDFile{},
		&DDEvent{},
		&ReviewQuestion{},
		&ProjectFile{},
	)
	//TODO 添加测试角色
	roleA := &OauthRole{
		Name: "管理员",
	}
	roleA.ID = 1

	roleB := &OauthRole{
		Name: "普通用户",
	}
	roleB.ID = 2

	db.Create(roleA)
	db.Create(roleB)
	return nil
}
