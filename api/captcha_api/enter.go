package captcha_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
)

type CaptchaApi struct {
}

type CaptchaResponse struct {
	CaptchaID string `json:"captchaID"`
	Captcha   string `json:"captcha"`
}

func (CaptchaApi) CaptchaView(c *gin.Context) {
	captchaConfig := base64Captcha.DriverString{
		Height:          60,
		Width:           200,
		NoiseCount:      1,
		ShowLineOptions: 2 | 4,
		Length:          4,
		Source:          "1234567890",
	}

	driver := captchaConfig.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, global.CaptchaStore)
	captchaID, captchaBase64, _, err := captcha.Generate()
	if err != nil {
		slog.Error("生成图片验证码失败", "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "err", err)
		res.FailWithMsg("图片验证码生成失败", c)
		return
	}

	res.OkWithData(CaptchaResponse{
		CaptchaID: captchaID,
		Captcha:   captchaBase64,
	}, c)
}
