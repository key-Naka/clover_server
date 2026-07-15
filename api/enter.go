package api

import (
	"clover_server/api/image_api"
	"clover_server/api/log_api"
	"clover_server/api/site_api"
)

type Api struct {
	ImageApi image_api.ImageApi
	SiteApi  site_api.SiteApi
	LogApi   log_api.LogApi
}

var App = Api{}
