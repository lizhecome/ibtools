package project

import (
	"encoding/json"
	"errors"
	"ibtools_server/models"
	policy "ibtools_server/util/policy"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (

	// ErrPrjTemplateNotFound ...
	ErrPrjTemplateNotFound = errors.New("项目模板无法找到")
	// ErrPrjTemplateNotFound ...

	ErrModelTemplateNotFound = errors.New("尽调模块模板无法找到")
	// ErrPrjNotFound ...
	ErrPrjNotFound = errors.New("项目无法找到")
)

//UserNameExists 通过模板创建项目
func (s *Service) createPrjByTemplate(templateCode, title, role string, user models.OauthUser) (*models.Project, error) {
	template, err := s.getTemplateByID(templateCode)
	if err != nil {
		return nil, err
	}
	jsonPessoal, _ := json.Marshal(template)

	var project models.Project
	json.Unmarshal(jsonPessoal, &project)
	//project.Code = uuid.New().String()
	project.Tilte = title
	project.Users = make([]models.OauthUser, 0)
	project.Users = append(project.Users, user)
	project.IsTemplate = 0
	project.Code = uuid.Nil
	project.OauthUserID = user.ID
	if err := s.db.Create(&project).Error; err != nil {
		return nil, err
	}

	user2pro := new(models.ProjectUser)
	if err := s.db.Where(&models.ProjectUser{ProjectID: project.ID, OauthUserID: user.ID}).First(&user2pro).Error; err != nil {
		return nil, err
	}
	user2pro.RoleName = role
	if err := s.db.Save(&user2pro).Error; err != nil {
		return nil, err
	}

	if bucket, err := s.ossClient.Bucket("ibtools"); err != nil {
		return nil, err
	} else {
		e := bucket.PutObject("projects/"+project.Code.String()+"/create.txt", strings.NewReader(project.Tilte))
		print(e)
	}

	return &project, nil
}

func (s *Service) getTemplateByID(templateCode string) (*models.Project, error) {
	project := new(models.Project)
	err := s.db.Preload("DDModelList.DDItems").Preload("DDModelList.Questions").Preload("DDModelList.Questions.AdditionDDItem").Preload(clause.Associations).Where("is_template = 1 and code = ?", templateCode).
		First(project).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPrjTemplateNotFound
	}

	return project, nil
}

func (s *Service) getAllTemplate() (*[]models.Project, error) {
	var projects []models.Project
	s.db.Where("is_template = 1").Find(&projects)

	return &projects, nil
}

func (s *Service) getSTSToken(user *models.OauthUser, projectCode string) (string, string, string, error) {
	request := sts.CreateAssumeRoleRequest()
	request.Scheme = "https"

	request.RoleArn = s.cnf.Aliyun.OssRoleARN
	request.RoleSessionName = user.UserName
	request.Method = "POST"
	request.DurationSeconds = (requests.Integer)(s.cnf.Aliyun.OssPermissionDurationSeconds)
	request.Policy, _ = policy.GetOssPolicy(projectCode)
	//TODO 处理Error
	response, err := s.stsClient.AssumeRole(request)
	if err != nil {
		return "", "", "", err
	} else {
		return response.Credentials.AccessKeyId, response.Credentials.AccessKeySecret, response.Credentials.SecurityToken, nil
	}
}

func (s *Service) getMyProjets(user *models.OauthUser) ([]*models.Project, error) {
	results := make([]*models.Project, 0)
	s.db.Model(&user).Association("Projects").Find(&results)
	return results, nil
}

func (s *Service) getFullProject(user *models.OauthUser, projectCode string) (*models.Project, error) {
	results := make([]*models.Project, 0)
	s.db.Model(&user).Where("code = ?", projectCode).Association("Projects").Find(&results)
	if len(results) == 0 {
		return nil, ErrPrjNotFound
	}

	result := new(models.Project)
	result.ID = results[0].ID
	s.db.Preload("DDModelList.DDItems").Preload("DDModelList.Questions").Preload("DDModelList.Questions.AdditionDDItem").Preload(clause.Associations).Find(result)

	return result, nil
}

// 创建尽调模块
func (s *Service) createDDModel(projectCode, title string, user models.OauthUser) (*models.DDModel, error) {
	var ddModel models.DDModel
	ddModel.Code = uuid.Nil

	results := make([]*models.Project, 0)
	s.db.Model(&user).Where("code = ?", projectCode).Association("Projects").Find(&results)
	if len(results) == 0 {
		return nil, ErrPrjNotFound
	}

	var dm models.DDModel
	s.db.Model(&models.DDModel{}).Where("project_id = ?", results[0].ID).Order("\"order\" DESC").Limit(1).Find(&dm)

	ddModel.Title = title
	ddModel.ProjectID = results[0].ID
	ddModel.Order = dm.Order + 1
	ddModel.Code = uuid.Nil
	ddModel.DDItems = make([]models.DDItem, 0)
	ddModel.Questions = make([]models.ReviewQuestion, 0)
	if err := s.db.Create(&ddModel).Error; err != nil {
		return nil, err
	}

	return &ddModel, nil
}

func (s *Service) createDDModelByTemplate(projectcode, modelTemplatecode string, user models.OauthUser) (*models.DDModel, error) {

	project, err := s.getFullProject(&user, projectcode)
	if err != nil {
		return nil, err
	}

	modelTemplate, err := s.getModelTemplateByID(modelTemplatecode)
	if err != nil {
		return nil, err
	}
	jsonPessoal, _ := json.Marshal(modelTemplate)

	var model models.DDModel
	json.Unmarshal(jsonPessoal, &model)
	model.Code = uuid.Nil
	max := 0
	for i := 0; i < len(project.DDModelList); i++ {
		if project.DDModelList[i].Order > 0 {
			max = project.DDModelList[i].Order
		}
	}
	model.Order = max + 1
	s.db.Model(project).Association("DDModelList").Append(&model)

	return &model, nil
}

func (s *Service) getModelTemplateByID(templateCode string) (*models.DDModel, error) {
	model := new(models.DDModel)
	err := s.db.Preload("DDModelList.DDItems").Preload("DDModelList.Questions").Preload("DDModelList.Questions.AdditionDDItem").Preload(clause.Associations).Where("code = ?", templateCode).
		First(model).Error

	// Not found
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrModelTemplateNotFound
	}

	return model, nil
}

//创建条目
func (s *Service) createDDItem(ownerType, title, filePointer string, ownerID uint, user models.OauthUser) (*models.DDItem, error) {
	var ddItem models.DDItem
	ddItem.Code = uuid.Nil

	var di models.DDItem
	s.db.Model(&models.DDItem{}).Where("owner_id = ? and owner_type = ?", ownerID, ownerType).Order("\"order\" DESC").Limit(1).Find(&di)

	ddItem.Title = title
	ddItem.Order = di.Order + 1
	ddItem.Code = uuid.Nil
	// fixme
	ddItem.ReviewMethod = ""
	ddItem.FilePointer = ""
	ddItem.Status = ""
	ddItem.Events = make([]models.DDEvent, 0)
	ddItem.Comments = make([]models.Comment, 0)
	ddItem.CollectFiles = make([]models.DDFile, 0)
	if err := s.db.Create(&ddItem).Error; err != nil {
		return nil, err
	}

	return &ddItem, nil
}
