package user_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/service/log_service"
	"clover_server/service/user_service"
	"clover_server/utils/jwts"
	"clover_server/utils/pwd"

	"github.com/gin-gonic/gin"
)

type PwdLoginRequest struct {
	Val      string `json:"val" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PwdLoginSwaggerRequest 是密码登录的 Swagger 请求模型，包含验证码中间件读取的字段。
type PwdLoginSwaggerRequest struct {
	Val      string `json:"val" binding:"required"`
	Password string `json:"password" binding:"required"`
	ID       string `json:"id"`
	Code     string `json:"code"`
}

// PwdLoginApi 使用用户名或邮箱和密码登录。
// @Summary 密码登录
// @Description 公开接口。`val` 为用户名或邮箱，`password` 为账户密码；启用图形验证码时，验证码中间件额外校验 `id`（验证码 ID）和 `code`（验证码文本）。
// @Tags 用户
// @Accept json
// @Produce json
// @Param request body PwdLoginSwaggerRequest true "用户名或邮箱、密码及可选图形验证码"
// @Success 200 {object} res.TokenResponse "登录成功，data 为 JWT"
// @Failure 400 {object} res.ErrorResponse "参数校验、图形验证码或账号密码错误"
// @Failure 500 {object} res.ErrorResponse "服务异常"
// @Router /user/login [post]
func (UserApi) PwdLoginApi(c *gin.Context) {
	var cr PwdLoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	if !global.Config.Site.Login.UsernamePwdLogin {
		res.FailWithMsg("站点未启用密码登陆", c)
		return
	}
	var user models.UserModel
	err = global.DB.Take(&user, "(username = ? or email = ?) and password  <> ''", cr.Val, cr.Val).Error
	if err != nil {
		log_service.NewLoginFail(c, enum.UserPwdLoginType, "用户不存在", cr.Val, cr.Password)
		res.FailWithMsg("用户不存在或密码错误", c)
		return
	}
	if !pwd.CompareHashAndPassword(user.Password, cr.Password) {
		log_service.NewLoginFail(c, enum.UserPwdLoginType, "密码错误", cr.Val, cr.Password)
		res.FailWithMsg("用户不存在或密码错误", c)
		return
	}
	token, _ := jwts.GetToken(jwts.Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     enum.RoleType(user.Role),
	})
	user_service.NewUserService(user).UserLogin(c)

	res.OkWithData(token, c)
}
