package oauth

import (
	"errors"
	"fmt"
	"ibtools_server/models"
	"ibtools_server/util"
	pass "ibtools_server/util/password"
	"io"
	"strings"
	"time"

	"gorm.io/gorm"
)

var (
	// MinPasswordLength defines minimum password length
	MinPasswordLength = 6

	// ErrPasswordTooShort ...
	ErrPasswordTooShort = fmt.Errorf(
		"password must be at least %d characters long",
		MinPasswordLength,
	)
	// ErrUserNotFound ...
	ErrUserNotFound = errors.New("user not found")
	// ErrInvalidUserPassword ...
	ErrInvalidUserPassword = errors.New("invalid user password")
	// ErrCannotSetEmptyUsername ...
	ErrCannotSetEmptyUsername = errors.New("cannot set empty username")

	// ErrCannotSetEmptyUseremail ...
	ErrCannotSetEmptyUseremail = errors.New("cannot set empty user email")

	// ErrCannotSetEmptyUserphone ...
	ErrCannotSetEmptyUserphone = errors.New("cannot set empty user phone")
	// ErrUserPasswordNotSet ...
	ErrUserPasswordNotSet = errors.New("user password not set")
	// ErrUsernameTaken ...
	ErrUsernameTaken = errors.New("username taken")
	// ErrUserphoneTaken ...
	ErrUserphoneTaken = errors.New("user phone taken")
	// ErrUseremailTaken ...
	ErrUseremailTaken = errors.New("user email taken")

	//
	ErrInviteCodeNotFound = errors.New("邀请码未找到")
)

//UserPhoneExists 用户手机号是否存在
func (s *Service) UserPhoneExists(phone string) bool {
	_, err := s.FindUserByPhone(phone)
	return err == nil
}

//UserNameExists 用户名是否存在
func (s *Service) UserNameExists(username string) bool {
	_, err := s.FindUserByUserName(username)
	return err == nil
}

//UserEmailExists 用户email是否存在
func (s *Service) UserEmailExists(email string) bool {
	_, err := s.FindUserByEmail(email)
	return err == nil
}

//FindUserByUserID 通过用户ID找用户
func (s *Service) FindUserByUserID(userid string) (*models.OauthUser, error) {
	user := new(models.OauthUser)
	err := s.db.Where("id = LOWER(?)", userid).
		First(user).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return user, nil
}

//FindUserByEmail 通过email找用户
func (s *Service) FindUserByEmail(email string) (*models.OauthUser, error) {
	user := new(models.OauthUser)
	err := s.db.Where("email = LOWER(?)", email).
		First(user).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	return user, nil
}

//FindUserByPhone 通过手机号找用户
func (s *Service) FindUserByPhone(phone string) (*models.OauthUser, error) {
	user := new(models.OauthUser)
	err := s.db.Where("phone = LOWER(?)", phone).
		First(user).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	return user, nil
}

//FindUserByPhone 通过用户名找用户
func (s *Service) FindUserByUserName(userName string) (*models.OauthUser, error) {
	user := new(models.OauthUser)
	err := s.db.Where("user_name = LOWER(?)", userName).
		First(user).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}

	return user, nil
}

//CreateUser 创建用户
func (s *Service) CreateUser(phone, invitationCode, roleName, username, password, company, title, email string) (*models.OauthUser, error) {

	role, _ := s.FindRoleByName(roleName)

	user := &models.OauthUser{
		RoleID:      role.ID,
		Phone:       phone,
		UserName:    username,
		Password:    util.StringOrNull(""),
		CompanyName: company,
		Title:       title,
		Email:       email,
	}

	// If the password is being set already, create a bcrypt hash
	if password != "" {
		if len(password) < MinPasswordLength {
			return nil, ErrPasswordTooShort
		}
		passwordHash, err := pass.HashPassword(password)
		if err != nil {
			return nil, err
		}
		user.Password = util.StringOrNull(string(passwordHash))
	}

	// Check the useremail is available
	if s.UserPhoneExists(user.Phone) {
		return nil, ErrUserphoneTaken
	}
	// 处理邀请码

	if invitationCode != "" {
		invitecode := new(models.OauthInvitationCode)
		err := s.db.Where("code = ? and expires_at > ?", invitationCode, time.Now().UTC()).First(&invitecode).Error
		if err == nil {
			project := new(models.Project)
			err := s.db.Where("code = ?", invitecode.ProjectCode).First(&project).Error
			if err == nil {
				role, _ = s.FindRoleByName(invitecode.RoleName)
				user.Projects = make([]models.Project, 0)
				user.Projects = append(user.Projects, *project)
				user.RoleID = role.ID
			}
		} else {
			return nil, err
		}
	}

	// Create the user
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

//SetPassword 设置密码
func (s *Service) SetPassword(user *models.OauthUser, password string) error {
	if len(password) < MinPasswordLength {
		return ErrPasswordTooShort
	}

	// Create a bcrypt hash
	passwordHash, err := pass.HashPassword(password)
	if err != nil {
		return err
	}

	// Set the password on the user object
	return s.db.Model(user).UpdateColumns(models.OauthUser{
		Password: util.StringOrNull(string(passwordHash)),
	}).Error
}

//UpdateUseremail 更新用户邮箱
func (s *Service) UpdateUseremail(user *models.OauthUser, email string) error {
	if email == "" {
		return ErrCannotSetEmptyUseremail
	}
	return s.db.Model(user).UpdateColumn("email", strings.ToLower(email)).Error
}

//UpdateUserphone 更新用户电话号码
func (s *Service) UpdateUserphone(user *models.OauthUser, phone string) error {
	if phone == "" {
		return ErrCannotSetEmptyUserphone
	}
	return s.db.Model(user).UpdateColumn("phone", phone).Error
}

//AuthUserByEmail 通过email认证用户
func (s *Service) AuthUserByEmail(email, thePassword string) (*models.OauthUser, error) {
	// Fetch the user
	user, err := s.FindUserByEmail(email)
	if err != nil {
		return nil, err
	}

	// Check that the password is set
	if !user.Password.Valid {
		return nil, ErrUserPasswordNotSet
	}

	// Verify the password
	if pass.VerifyPassword(user.Password.String, thePassword) != nil {
		return nil, ErrInvalidUserPassword
	}

	return user, nil
}

//AuthUserByUserName 通过username认证用户
func (s *Service) AuthUserByUserName(userName, thePassword string) (*models.OauthUser, error) {
	// Fetch the user
	user, err := s.FindUserByUserName(userName)
	if err != nil {
		return nil, err
	}

	// Check that the password is set
	if !user.Password.Valid {
		return nil, ErrUserPasswordNotSet
	}

	// Verify the password
	if pass.VerifyPassword(user.Password.String, thePassword) != nil {
		return nil, ErrInvalidUserPassword
	}

	return user, nil
}

//AuthUserByPhone 通过手机号认证用户
func (s *Service) AuthUserByPhone(phone, thePassword string) (*models.OauthUser, error) {
	// Fetch the user
	user, err := s.FindUserByPhone(phone)
	if err != nil {
		return nil, err
	}

	// Check that the password is set
	if !user.Password.Valid {
		return nil, ErrUserPasswordNotSet
	}

	// Verify the password
	if pass.VerifyPassword(user.Password.String, thePassword) != nil {
		return nil, ErrInvalidUserPassword
	}

	return user, nil
}

//UploadUserImage ...
func (s *Service) UploadUserImage(user *models.OauthUser, filename string, reader *io.Reader) (string, error) {
	_, cdnpath, err := util.UpLoadFileFromByte(filename, "userimage", reader)
	if err != nil {
		return "", err
	}
	if err := s.db.Model(user).UpdateColumn("image_link", cdnpath).Error; err != nil {
		return "", err
	}
	return cdnpath, nil
}
