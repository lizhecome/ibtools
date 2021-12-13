package drerror

import "github.com/gin-gonic/gin"

//APIError 接口错误
type APIError struct {
	Code  int    `json:"status"`
	Title string `json:"title"`
}

var (
	//APIErrPhoneNumInValidate 手机号非法
	APIErrPhoneNumInValidate = &APIError{601, "Cell phone number is not correct"}
	//APIErrEmailInValidate mail非法
	APIErrEmailInValidate = &APIError{602, "The email is invalid."}
	//APIErrEmailorPhoneNumInValidate mail非法
	APIErrEmailorPhoneNumInValidate = &APIError{603, "Email or cell phone number is not correct"}
	//APIErrPhoneValidateFailed 手机号验证失败
	APIErrPhoneValidateFailed = &APIError{604, "Phone number verification failed"}
	//APIErrCreateUserFailed 创建用户失败
	APIErrCreateUserFailed = &APIError{605, "Create User Failed"}
	//APIErrLoginFailed 登陆失败
	APIErrLoginFailed = &APIError{606, "Login Failed"}
	//APIErrSendPhoneNumValidateMessageFailed 发送手机验证码失败
	APIErrSendPhoneNumValidateMessageFailed = &APIError{607, "Send Phone Number ValidateMessage Failed"}
	//APIErrUserNotFound 找不到用户
	APIErrUserNotFound = &APIError{608, "User not Found"}
	//APIErrSetPasswordFailed 设置密码失败
	APIErrSetPasswordFailed = &APIError{609, "Set Password Failed"}
	//APIErrRefreshTokenInValidate 手机号非法
	APIErrRefreshTokenInValidate = &APIError{610, "Refresh token is not correct"}
	//APIErrFirstOrLastNameIsNull firstname或者lastname为空
	APIErrFirstOrLastNameIsNull = &APIError{611, "First or last name is null"}
	//APIErrAccessTokenIsNull AccessToken为空
	APIErrAccessTokenIsNull = &APIError{612, "AccessToken is null"}
	//APIErrAuthFailed AccessToken认证失败
	APIErrAuthFailed = &APIError{613, "AccessToken Auth Failed"}
	//APIErrParamIsNull 参数为空
	APIErrParamIsNull = &APIError{614, "Param is Null"}
	//APIErrPraseFailed 参数解析失败
	APIErrPraseFailed = &APIError{615, "Param prase failed"}

	//APIErrAccessDeny 拒绝访问
	APIErrAccessDeny = &APIError{616, "Access Deny"}

	//APIErrOther 其他错误
	APIErrOther = &APIError{617, "other error"}
	//APIErrCodeIsNull code为空
	APIErrCodeIsNull = &APIError{618, "code is null"}
	//APIErrLimitIsNull 额度已经用尽
	APIErrLimitIsNull = &APIError{619, "额度已经用尽"}
	//APIErrTokenIsNull token为空
	APIErrTokenIsNull = &APIError{620, "token is null"}
	//APIErrCreateProjectFailed 创建项目失败
	APIErrCreateProjectFailed = &APIError{621, "创建项目失败"}
	//APIErrGetTemplatesFailed 获取模板失败
	APIErrGetTemplatesFailed = &APIError{622, "获取模板失败"}
	//APIErrGetTemplatesFailed 获取权限失败
	APIErrGetPermissionFailed = &APIError{623, "获取权限失败"}
	//APIErrGetProjectsFailed 获取权限失败
	APIErrGetProjectsFailed = &APIError{624, "获取项目失败"}
	//APIErrCreateInviteCodeFailed 获取权限失败
	APIErrCreateInviteCodeFailed = &APIError{625, "生成邀请码失败"}
	//APIErrCreateModelFailed 创建项目失败
	APIErrCreateModelFailed = &APIError{626, "创建尽调模块失败"}
	//APIErrCreateInviteFailed 邀请失败
	APIErrCreateInviteFailed = &APIError{627, "邀请失败尽调模块"}
	//APIErrCreateItemFailed 创建尽调条目失败
	APIErrCreateItemFailed = &APIError{628, "创建尽调条目失败"}
	//ProcessInviteFailed 处理邀请项目失败
	ProcessInviteFailed = &APIError{629, "处理邀请项目失败"}
	//GetInviteListToMeFailed 获取邀请我的列表失败
	GetInviteListToMeFailed = &APIError{630, "获取邀请我的列表失败"}
	//GetInviteListFromMeFailed 获取邀请我的列表失败
	GetInviteListFromMeFailed = &APIError{631, "获取我发出的邀请列表失败"}
)

func ResponseError(c *gin.Context, apierr *APIError, e error) {

	errorJSON := gin.H{
		"title": apierr.Title,
	}
	if e != nil {
		errorJSON["message"] = e.Error()
	} else {
		errorJSON["message"] = ""
	}

	c.JSON(apierr.Code, errorJSON)
}
