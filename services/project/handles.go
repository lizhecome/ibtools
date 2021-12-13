package project

import (
	"errors"
	"ibtools_server/drerror"
	"ibtools_server/models"

	"github.com/gin-gonic/gin"
)

// 新建项目
// @Summary 根据模板编号创建项目
// @Description 通过坐标查找附近的地点
// @ID createPrjByTemplateHandler
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param title body string true "XX公司IPO项目"
// @Param role body string true "项目组长（可以填四项：发行人员工、发行人负责人、项目组成员、项目组长）"
// @Param templatecode body string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Success 200 {string} string "{project}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 614 {string} string "{"Param is Null","message": ""}"
// @Failure 621 {string} string "{"创建项目失败": ""}"
// @Router /project/createprojectbytemplate [post]
func (s *Service) createPrjByTemplateHandler(c *gin.Context) {

	title := c.PostForm("title")
	if title == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("title"))
		c.Abort()
		return
	}

	role := c.PostForm("role")
	if role == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("role"))
		c.Abort()
		return
	}

	code := c.PostForm("templatecode")
	if code == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("templatecode"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
		return
	}

	project, e := s.createPrjByTemplate(code, title, role, *user.(*models.OauthUser))
	if e != nil {
		drerror.ResponseError(c, drerror.APIErrCreateProjectFailed, e)
		c.Abort()
		return
	} else {
		c.JSON(200, project)
	}

}

// 获取所有模板
// @Summary 获取所有模板
// @Description 获取所有模板
// @ID getAllTemplateHandler
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Success 200 {string} string "{projects}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 622 {string} string "{"获取模板失败": ""}"
// @Router /project/getalltemplates [post]
func (s *Service) getAllTemplateHandler(c *gin.Context) {
	if templates, err := s.getAllTemplate(); err != nil {
		drerror.ResponseError(c, drerror.APIErrGetTemplatesFailed, nil)
		c.Abort()
	} else {
		c.JSON(200, templates)
	}

}

// 获取项目文件的读写权限
// @Summary 获取项目文件的读写权限
// @Description 获取项目文件的读写权限
// @ID getProjectFilePermission
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param projectcode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Success 200 {string} string "{token}}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 623 {string} string "{"获取权限失败": ""}"
// @Router /project/getprojectfilepermission [post]
func (s *Service) getProjectFilePermissionHandler(c *gin.Context) {
	code := c.PostForm("projectcode")
	if code == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("projectcode"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}

	stoken := new(STSToken)
	stoken.Domain = s.cnf.Aliyun.OssAccelerateDomain
	stoken.EndPoint = s.cnf.Aliyun.OssAccelerateEndPoint
	stoken.BucketName = s.cnf.Aliyun.BucketName
	stoken.Path = s.cnf.Aliyun.ProjectRootPath + code
	var err error
	stoken.AccessKeyId, stoken.AccessKeySecret, stoken.SecurityToken, err = s.getSTSToken(user.(*models.OauthUser), code)
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrGetPermissionFailed, err)
		c.Abort()
	} else {
		c.JSON(200, stoken)
	}
}

type STSToken struct {
	AccessKeyId     string
	AccessKeySecret string
	SecurityToken   string
	Domain          string
	EndPoint        string
	BucketName      string
	Path            string
}

// 获取我的所有项目
// @Summary 获取我的所有项目
// @Description 获取我的所有项目
// @ID getMyProjectsHandler
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Success 200 {string} string "{projects}}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 624 {string} string "{"获取项目失败": ""}"
// @Router /project/getmyprojects [post]
func (s *Service) getMyProjectsHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}

	results, err := s.getMyProjets(user.(*models.OauthUser))
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrGetProjectsFailed, err)
		c.Abort()
	} else {
		c.JSON(200, results)
	}

}

// 获得一个项目的全部内容
// @Summary 获得一个项目的全部内容
// @Description 获得一个项目的全部内容
// @ID getFullProjectHandler
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param projectcode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Success 200 {string} string "{projects}}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 624 {string} string "{"获取项目失败": ""}"
// @Router /project/getfullproject [post]
func (s *Service) getFullProjectHandler(c *gin.Context) {
	code := c.PostForm("projectcode")
	if code == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("projectcode"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}

	result, err := s.getFullProject(user.(*models.OauthUser), code)
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrGetProjectsFailed, err)
		c.Abort()
	} else {
		c.JSON(200, result)
	}

}

// 创建一个项目的尽调模块
// @Summary 创建一个项目的尽调模块
// @Description 创建一个项目的尽调模块
// @ID createDDModelHandler
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param projectcode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param title body string true "历史沿革"
// @Success 200 {string} string "{200}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 625 {string} string "{"创建尽调模块失败": ""}"
// @Router /project/createDDModel [post]
func (s *Service) createDDModelHandler(c *gin.Context) {
	code := c.PostForm("projectcode")
	if code == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("projectcode"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}

	title := c.PostForm("title")
	if title == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("title"))
	}

	results, err := s.getMyProjets(user.(*models.OauthUser))
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrGetProjectsFailed, err)
		c.Abort()
	} else {
		if len(results) > 0 {
			ddModel, e := s.createDDModel(code, title, *user.(*models.OauthUser))
			if e != nil {
				drerror.ResponseError(c, drerror.APIErrCreateModelFailed, e)
				c.Abort()
			} else {
				c.JSON(200, ddModel.Code)
			}
		} else {
			c.JSON(200, "")
		}
	}

}

// 根据模板创建尽调模块
// @Summary 根据模板创建尽调模块
// @Description 根据模板创建尽调模块
// @ID createDDModelByTemplateHandler
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param modeltemplatecode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param projectcode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Success 200 {string} string "{projects}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 625 {string} string "{"创建尽调模块失败": ""}"
// @Router /project/createDDModelByTemplate [post]
func (s *Service) createDDModelByTemplateHandler(c *gin.Context) {
	projectcode := c.PostForm("projectcode")
	if projectcode == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("projectcode"))
		c.Abort()
	}

	modeltemplatecode := c.PostForm("modeltemplatecode")
	if modeltemplatecode == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("modeltemplatecode"))
		return
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		return
	}

	result, err := s.createDDModelByTemplate(projectcode, modeltemplatecode, *user.(*models.OauthUser))
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrCreateModelFailed, nil)
		return
	} else {
		c.JSON(200, result)
	}
}

// 创建条目
// @Summary 创建条目
// @Description 创建条目
// @ID createDDItemHandler
// @Tags project
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param title body string true "21-1-2"
// @Param filePointer body string true "21-1-2"
// @Success 200 {string} string "{projects}}"
// @Failure 612 {string} string "{"AccessToken is nul","message": ""}"
// @Failure 613 {string} string "{"AccessToken Auth Failed","message": ""}"
// @Failure 626 {string} string "{"创建条目失败": ""}"
// @Router /project/createDDItem [post]
func (s *Service) createDDItemHandler(c *gin.Context) {
	title := c.PostForm("title")
	if title == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("title"))
		return
	}

	filePointer := c.PostForm("filePointer")
	if filePointer == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("filePointer"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		return
	}
	//fixme ownerType ownerId
	result, err := s.createDDItem("", title, filePointer, 0, *user.(*models.OauthUser))
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrCreateModelFailed, nil)
		return
	} else {
		c.JSON(200, result)
	}
}
