package api

import (
	"clover_server/api/site_api"
)

type Api struct {
	SiteApi site_api.SiteApi
}

var App = Api{}
