package image_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/service/qiniu_service"

	"github.com/gin-gonic/gin"
)

type QiNiuGenTokenResponse struct {
	Key    string `json:"key"`
	Token  string `json:"token"`
	Region string `json:"region"`
}

// QiNiuGenToken 获取七牛云直传凭证。
// @Summary 获取七牛云上传凭证
// @Tags 图片
// @Produce json
// @Security BearerAuth
// @Success 200 {object} res.QiNiuTokenResponse "七牛 bucket、上传 token 与区域"
// @Failure 400 {object} res.ErrorResponse "七牛云未启用"
// @Failure 500 {object} res.ErrorResponse "生成上传凭证失败"
// @Router /images/qiniu [post]
func (i ImageApi) QiNiuGenToken(c *gin.Context) {
	glo := global.Config.QiNiu
	if !glo.Enable {
		res.FailWithMsg("七牛云未开启", c)
		c.Abort()
		return
	}
	token, err := qiniu_service.GenToken()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	req := &QiNiuGenTokenResponse{
		Key:    glo.Bucket,
		Token:  token,
		Region: glo.Region,
	}
	res.OkWithData(req, c)
}
