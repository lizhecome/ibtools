package oauth

import (
	"errors"
	"ibtools_server/drerror"
	"ibtools_server/models"
	"ibtools_server/util"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// 用户注册
// @Summary 用户注册
// @Description 通过用户名、昵称、角色及密码进行用户注册
// @ID registerHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param name body string true "李喆"
// @Param password body string true "123456"
// @Param phone body string true "5555555555"
// @Param cpmpany body string true "XXX证券公司"
// @Param title body string true "高级经理"
// @Param email body string true "abc@a.com"
// @Param code body string true "6666"
// @Param invitationcode body string true "h7jgyg6fsHe5"
// @Success 200 {string} string "{"message": "ok"}"
// @Failure 605 {string} string "{"title":"Create User Failed","message": ""}"
// @Router /oauth/register [post]
func (s *Service) registerHandler(c *gin.Context) {
	cpmpany := c.PostForm("company")
	if cpmpany == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("公司名不能为空"))
		c.Abort()
	}

	title := c.PostForm("title")
	if title == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("职位不能为空"))
		c.Abort()
	}

	email := c.PostForm("email")
	if email == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("电子邮箱不能为空"))
		c.Abort()
	}

	phone := c.PostForm("phone")
	if !util.ValidatePhone(phone) {
		drerror.ResponseError(c, drerror.APIErrPhoneNumInValidate, nil)
		return
	}

	if c.PostForm("code") == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("验证码不能为空"))
		return
	}

	if err := s.PhoneNumValidate(c.PostForm("phone"), c.PostForm("code")); err != nil {
		drerror.ResponseError(c, drerror.APIErrPhoneValidateFailed, err)
		return
	}

	if c.PostForm("name") == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("姓名不能为空"))
		return
	}

	//invitationCode, roleName, username, nickname, password string
	_, err := s.CreateUser(phone, c.PostForm("invitationcode"), "普通用户", c.PostForm("name"), c.PostForm("password"), cpmpany, title, email)
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrCreateUserFailed, err)
	} else {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	}

}

// 用户登陆
// @Summary 用户登陆
// @Description 通过用户名及密码进行用户登陆
// @ID loginHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param phone body string true "186XXXXXXXX"
// @Param password body string true "123456"
// @Success 200 {string} string "{"ExpiresAt": "6019-11-14 08:29:48.446901 +0000 UTC","accesstoken": "1e0e848b-38d3-4336-845e-5b09d4065608","message": "ok","refreshtoken": "12aa1a3f-c0d5-4934-a73b-88af2c02fa5d"}"
// @Failure 606 {string} string "{"title":"Login Failed","message": ""}"
// @Router /oauth/login [post]
func (s *Service) loginHandler(c *gin.Context) {
	phone := c.PostForm("phone")
	if phone != "" {
		user, err := s.AuthUserByPhone(phone, c.PostForm("password"))
		if err != nil {
			drerror.ResponseError(c, drerror.APIErrLoginFailed, err)
			return
		}
		at, rt, err := s.Login(user, "app")
		if err != nil {
			drerror.ResponseError(c, drerror.APIErrLoginFailed, err)
			return
		}
		c.JSON(200, gin.H{
			"accesstoken":  at.Token,
			"ExpiresAt":    at.ExpiresAt.String(),
			"refreshtoken": rt.Token,
			"user":         user,
		})
	}
}

// 刷新accesstoken
// @Summary 刷新accesstoken
// @Description 使用refreshtoken刷新accesstoken
// @ID refreshtokenHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param refresh_token body string true "12aa1a3f-c0d5-4934-a73b-88af2c02fa5d"
// @Success 200 {string} string "{"ExpiresAt": "6019-11-14 08:29:48.446901 +0000 UTC","accesstoken": "1e0e848b-38d3-4336-845e-5b09d4065608","message": "ok","refreshtoken": "12aa1a3f-c0d5-4934-a73b-88af2c02fa5d"}"
// @Failure 606 {string} string "{"title":"Login Failed","message": ""}"
// @Failure 610 {string} string "{"title":"Refresh token is not correct": ""}"
// @Router /oauth/refreshtoken [post]
func (s *Service) refreshtokenHandler(c *gin.Context) {
	theRefreshToken, err := s.GetValidRefreshToken(c.PostForm("refresh_token"))
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrRefreshTokenInValidate, err)
		return
	}

	// Log in the user
	accessToken, refreshToken, err := s.Login(
		theRefreshToken.User,
		"app",
	)
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrLoginFailed, err)
		return
	}

	// Create response
	c.JSON(200, gin.H{
		"accesstoken":  accessToken.Token,
		"ExpiresAt":    accessToken.ExpiresAt.String(),
		"refreshtoken": refreshToken.Token,
	})
}

// 生成邀请码
// @Summary 生成邀请码
// @Description 生成邀请码
// @ID createInviteCodeHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param projectcode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param role body string true "项目组长"
// @Success 200 {string} string "invitecode"
// @Failure 614 {string} string "{"title":"Param is Null","message": ""}"
// @Failure 625 {string} string "{"title":"生成邀请码失败": ""}"
// @Router /oauth/createinvitecode [post]
func (s *Service) createInviteCodeHandler(c *gin.Context) {
	code := c.PostForm("projectcode")
	if code == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("projectcode"))
		c.Abort()
	}

	role := c.PostForm("role")
	if role == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("role"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}

	result, err := s.CreateInvite(user.(*models.OauthUser), code, uuid.NewString(), "", "", role)
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrCreateInviteCodeFailed, err)
		c.Abort()
	} else {
		c.JSON(200, result)
	}
}

// 通过手机号邀请
// @Summary 通过手机号邀请
// @Description 通过手机号邀请
// @ID createInviteByPhone
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param projectcode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param phone body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param role body string true "项目组长"
// @Success 200 {string} string "invitecode"
// @Failure 614 {string} string "{"title":"Param is Null","message": ""}"
// @Failure 625 {string} string "{"title":"生成邀请码失败": ""}"
// @Router /oauth/createinvitebyphone [post]
func (s *Service) createInviteByPhoneHandler(c *gin.Context) {
	projectCode := c.PostForm("projectcode")
	if projectCode == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("projectcode"))
		c.Abort()
	}

	phone := c.PostForm("phone")
	if phone == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("phone"))
		c.Abort()
	}

	role := c.PostForm("role")
	if role == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("role"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}

	result, err := s.CreateInvite(user.(*models.OauthUser), projectCode, uuid.NewString(), "", phone, role)
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrCreateInviteCodeFailed, err)
		c.Abort()
	} else {
		c.JSON(200, result)
	}
}

// 通过用户code邀请
// @Summary 通过用户code邀请
// @Description 通过用户code邀请
// @ID createInviteByUserCodeHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param projectcode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param usercode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param role body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Success 200 {string} string "invitecode"
// @Failure 614 {string} string "{"title":"Param is Null","message": ""}"
// @Failure 625 {string} string "{"title":"生成邀请码失败": ""}"
// @Router /oauth/createinvitebyusercode [post]
func (s *Service) createInviteByUserCodeHandler(c *gin.Context) {
	projectCode := c.PostForm("projectcode")
	if projectCode == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("projectcode"))
		c.Abort()
	}

	usercode := c.PostForm("usercode")
	if usercode == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("usercode"))
		c.Abort()
	}

	role := c.PostForm("role")
	if role == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("role"))
		c.Abort()
	}

	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}

	result, err := s.CreateInvite(user.(*models.OauthUser), projectCode, uuid.NewString(), usercode, "", role)
	if err != nil {
		drerror.ResponseError(c, drerror.APIErrCreateInviteCodeFailed, err)
		c.Abort()
	} else {
		c.JSON(200, result)
	}
}

// 发送短信验证码
// @Summary 发送短信验证码
// @Description 根据手机号发送验证码短信
// @ID sendverificationcodeHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param phone body string true "5555555555"
// @Success 200 {string} string "{"message": "ok"}"
// @Failure 607 {string} string "{"title":"Send Phone Number ValidateMessage Failed","message": ""}"
// @Router /oauth/sendverificationcode [post]
func (s *Service) sendverificationcodeHandler(c *gin.Context) {
	phone := c.PostForm("phone")
	if err := s.SendPhoneNumValidateMessage(phone); err != nil {
		drerror.ResponseError(c, drerror.APIErrSendPhoneNumValidateMessageFailed, err)
	} else {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	}
}

// 检查电话号码是否可用
// @Summary 检查电话号码是否可用
// @Description 检查电话号码是否可用
// @ID checkphonenumberavailableHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param phone body string true "5555555555"
// @Success 200 {string} string "{"message": "ok"}"
// @Failure 601 {string} string "{"title":"Cell phone number is not correct","message": ""}"
// @Router /oauth/checkphonenumberavailable [post]
func (s *Service) checkphonenumberavailableHandler(c *gin.Context) {
	phone := c.PostForm("phone")
	if !util.ValidatePhone(phone) {
		drerror.ResponseError(c, drerror.APIErrPhoneNumInValidate, nil)
		return
	}
	if s.UserPhoneExists(phone) {
		c.JSON(200, gin.H{
			"message": false,
		})
	} else {
		c.JSON(200, gin.H{
			"message": true,
		})
	}
}

// 忘记密码
// @Summary 忘记密码
// @Description 忘记密码
// @ID forgetpasswordHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param phone body string true "5555555555"
// @Param code body string true "6666"
// @Param password body string true "111111"
// @Success 200 {string} string "{"message": "ok"}"
// @Failure 601 {string} string "{"title":"Cell phone number is not correct","message": ""}"
// @Failure 608 {string} string "{"title":"User not Found","message": ""}"
// @Failure 609 {string} string "{"title":"Set Password Failed": ""}"
// @Router /oauth/forgetpassword [post]
func (s *Service) forgetpasswordHandler(c *gin.Context) {
	if err := s.PhoneNumValidate(c.PostForm("phone"), c.PostForm("code")); err != nil {
		drerror.ResponseError(c, drerror.APIErrPhoneValidateFailed, err)
	} else {
		user, err := s.FindUserByPhone(util.FormatPhoneNum(c.PostForm("phone")))
		if err != nil {
			drerror.ResponseError(c, drerror.APIErrUserNotFound, err)
			return
		}
		if err := s.SetPassword(user, c.PostForm("password")); err != nil {
			drerror.ResponseError(c, drerror.APIErrSetPasswordFailed, err)
		} else {
			c.JSON(200, gin.H{
				"message": "ok",
			})
		}
	}
}

// 处理邀请项目
// @Summary 处理邀请项目
// @Description 处理邀请项目
// @ID processInviteHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Param invitecode body string true "3cb6f508-ece8-43ae-b025-bdefd9b807a9"
// @Param accept body string true "0:拒绝，1:接受"
// @Success 200 {string} string "{"message": "ok"}"
// @Failure 614 {string} string "{"title":"Param is Null","message": ""}"
// @Failure 629 {string} string "{"title":"处理邀请项目失败": ""}"
// @Router /oauth/processInvite [post]
func (s *Service) processInviteHandler(c *gin.Context) {
	invitecode := c.PostForm("invitecode")
	if invitecode == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("invitecode"))
		c.Abort()
	}
	accept := c.PostForm("accept")
	if accept == "" {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, errors.New("accept"))
		c.Abort()
	}
	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}
	if e := s.ProcessInvite(user.(*models.OauthUser), invitecode, accept); e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			c.JSON(200, gin.H{
				"message": "找不到记录，请确认邀请是否已处理或已过期",
			})
		} else {
			drerror.ResponseError(c, drerror.ProcessInviteFailed, e)
			c.Abort()
		}
	} else {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	}
}

// 获取邀请我的列表
// @Summary 获取邀请我的列表
// @Description 获取邀请我的列表
// @ID getInviteListToMeHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Success 200 {string} string "list"
// @Failure 614 {string} string "{"title":"Param is Null","message": ""}"
// @Failure 630 {string} string "{"title":"获取邀请我的列表失败": ""}"
// @Router /oauth/getInviteListToMe [post]
func (s *Service) getInviteListToMeHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}
	if results, err := s.GetInviteListToMe(user.(*models.OauthUser)); err != nil {
		drerror.ResponseError(c, drerror.GetInviteListToMeFailed, nil)
		c.Abort()
	} else {
		c.JSON(200, results)
	}
}

// 获取我发出的邀请列表
// @Summary 获取我发出的邀请列表
// @Description 获取我发出的邀请列表
// @ID getInviteListFromMeHandler
// @Tags oauth
// @Accept  mpfd
// @Produce  json
// @Param at header string true "d29788b2-a482-4f1d-9434-a84a9cfbc01d"
// @Success 200 {string} string "list"
// @Failure 614 {string} string "{"title":"Param is Null","message": ""}"
// @Failure 631 {string} string "{"title":"获取我发出的邀请列表失败": ""}"
// @Router /oauth/getInviteListFromMe [post]
func (s *Service) getInviteListFromMeHandler(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		drerror.ResponseError(c, drerror.APIErrParamIsNull, nil)
		c.Abort()
	}
	if results, err := s.GetInviteListFromMe(user.(*models.OauthUser)); err != nil {
		drerror.ResponseError(c, drerror.GetInviteListToMeFailed, nil)
		c.Abort()
	} else {
		c.JSON(200, results)
	}
}
