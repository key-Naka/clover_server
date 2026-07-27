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
