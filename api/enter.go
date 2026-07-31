package api

import (
	"clover_server/api/article_api"
	"clover_server/api/banner_api"
	"clover_server/api/captcha_api"
	"clover_server/api/comment_api"
	"clover_server/api/data_api"
	"clover_server/api/focus_api"
	"clover_server/api/global_notification_api"
	"clover_server/api/image_api"
	"clover_server/api/log_api"
	"clover_server/api/site_api"
	"clover_server/api/user_api"
)

type Api struct {
	ImageApi              image_api.ImageApi
	SiteApi               site_api.SiteApi
	LogApi                log_api.LogApi
	BannerApi             banner_api.BannerApi
	CaptchaApi            captcha_api.CaptchaApi
	UserApi               user_api.UserApi
	ArticleApi            article_api.ArticleApi
	CommentApi            comment_api.CommentApi
	DataApi               data_api.DataApi
	FocusApi              focus_api.FocusApi
	GlobalNotificationApi global_notification_api.GlobalNotificationApi
}

var App = Api{}
