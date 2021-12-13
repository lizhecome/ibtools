package oauth

import (
	"errors"
	"ibtools_server/models"
	"time"
)

func (s *Service) CreateInvite(createUser *models.OauthUser, projectCode, invitecode, userCode, phone, role string) (*models.OauthInvitationCode, error) {
	inviteCode := &models.OauthInvitationCode{
		OauthUserID: createUser.ID,
		Code:        invitecode,
		ProjectCode: projectCode,
		Phone:       phone,
		UserCode:    userCode,
		RoleName:    role,
		ExpiresAt:   time.Now().UTC().Add(time.Duration(s.cnf.Oauth.InviteCodeLifetime) * time.Second),
	}
	err := s.db.Create(&inviteCode).Error
	if err != nil {
		return nil, err
	} else {
		return inviteCode, nil
	}
}

func (s *Service) ProcessInvite(processUser *models.OauthUser, invitecode, accept string) error {
	invite := new(models.OauthInvitationCode)
	err := s.db.Where("code = ? and status ='' and expires_at > ?", invitecode, time.Now().UTC()).First(&invite).Error

	if err == nil {
		if invite.Phone != "" || invite.UserCode != "" {
			if invite.Phone != "" && invite.Phone != processUser.Phone {
				return errors.New("只能处理和自己相关的邀请")
			}
			if invite.UserCode != "" && invite.UserCode != processUser.Code.String() {
				return errors.New("只能处理和自己相关的邀请")
			}
		}
		if accept == "0" {
			invite.Status = accept
			s.db.Save(&invite)
		} else if accept == "1" {
			invite.Status = accept
			s.db.Save(&invite)
			project := new(models.Project)
			err := s.db.Where("code = ?", invite.ProjectCode).First(&project).Error
			if err == nil {
				processUser.Projects = append(processUser.Projects, *project)
				if err := s.db.Save(processUser).Error; err != nil {
					return err
				}
				user2pro := new(models.ProjectUser)
				if err := s.db.Where(&models.ProjectUser{ProjectID: project.ID, OauthUserID: processUser.ID}).First(&user2pro).Error; err != nil {
					return err
				}
				user2pro.RoleName = invite.RoleName
				if err := s.db.Save(&user2pro).Error; err != nil {
					return err
				}
				return nil
			}
		} else {
			return errors.New("accept只能为0或1")
		}
	}
	return err
}

func (s *Service) GetInviteListToMe(user *models.OauthUser) (*[]models.OauthInvitationCode, error) {
	invites := new([]models.OauthInvitationCode)
	if err := s.db.Where("user_code = ? or phone = ?", user.Code, user.Phone).Find(&invites).Error; err != nil {
		return nil, err
	}

	return invites, nil
}

func (s *Service) GetInviteListFromMe(user *models.OauthUser) (*[]models.OauthInvitationCode, error) {
	invites := new([]models.OauthInvitationCode)
	if err := s.db.Where("oauth_user_id = ?", user.ID).Find(&invites).Error; err != nil {
		return nil, err
	}
	return invites, nil
}
