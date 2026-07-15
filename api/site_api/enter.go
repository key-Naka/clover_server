package site_api

import (
	"clover_server/common/res"
	"clover_server/conf"
	"clover_server/core"
	"clover_server/global"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type SiteApi struct {
}
type SiteInfoRequest struct {
	Name string `uri:"name"`
}

func (s *SiteApi) SiteinfoView(c *gin.Context) {
	var cr SiteInfoRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	if cr.Name == "site" {
		global.Config.Site.About.Version = global.Version
		res.OkWithData(global.Config.Site, c)
		return
	}
	_, ok := c.Get("claims")
	if !ok {
		res.FailWithMsg("权限不足", c)
		return
	}

	var data any
	switch cr.Name {
	case "site":
		data = global.Config.Site
	case "email":
		rep := global.Config.Email
		rep.AuthCode = "******"
		data = rep
	case "qq":
		rep := global.Config.QQ
		rep.AppKey = "******"
		data = rep
	case "qiniu":
		rep := global.Config.QiNiu
		rep.SecretKey = "******"
		data = rep
	case "ai":
		rep := global.Config.Ai
		rep.SecretKey = "******"
		data = rep
	default:
		res.FailWithMsg("不支持的配置项", c)
		return
	}
	res.OkWithData(data, c)
}
func (SiteApi) SiteInfoQQView(c *gin.Context) {
	res.OkWithData(global.Config.QQ.Url(), c)
}
func (SiteApi) SiteUpdateView(c *gin.Context) {
	var cr SiteInfoRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var rep any
	switch cr.Name {
	case "site":
		var data conf.Site
		err = c.ShouldBindJSON(&data)
		rep = data
	case "email":
		var data conf.Email
		err = c.ShouldBindJSON(&data)
		rep = data
	case "qq":
		var data conf.QQ
		err = c.ShouldBindJSON(&data)
		rep = data
	case "qiNiu":
		var data conf.QiNiu
		err = c.ShouldBindJSON(&data)
		rep = data
	case "ai":
		var data conf.Ai
		err = c.ShouldBindJSON(&data)
		rep = data
	default:
		res.FailWithMsg("不存在的配置", c)
		return
	}
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	switch s := rep.(type) {
	case conf.Site:
		err = UpdateSite(s)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		global.Config.Site = s
	case conf.Email:
		if s.AuthCode == "******" {
			s.AuthCode = global.Config.Email.AuthCode
		}
		global.Config.Email = s
	case conf.QQ:
		if s.AppKey == "******" {
			s.AppKey = global.Config.QQ.AppKey
		}
		global.Config.QQ = s
	case conf.QiNiu:
		if s.SecretKey == "******" {
			s.SecretKey = global.Config.QiNiu.SecretKey
		}
		global.Config.QiNiu = s
	case conf.Ai:
		if s.SecretKey == "******" {
			s.SecretKey = global.Config.Ai.SecretKey
		}
		global.Config.Ai = s
	}

	// 改配置文件
	core.SetConf()

	res.OkWithMsg("更新站点配置成功", c)
	return
}
func UpdateSite(site conf.Site) error {
	if site.Project.Icon == "" && site.Project.Title == "" &&
		site.Seo.Keywords == "" && site.Seo.Description == "" &&
		site.Project.WebPath == "" {
		return nil
	}

	if site.Project.WebPath == "" {
		return errors.New("请配置前端地址")
	}

	file, err := os.Open(site.Project.WebPath)
	if err != nil {
		return fmt.Errorf("%s 文件不存在", site.Project.WebPath)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			slog.Error("关闭站点文件失败", "err", closeErr, "path", site.Project.WebPath)
		}
	}()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		slog.Error("解析站点文件失败", "err", err, "path", site.Project.WebPath)
		return errors.New("文件解析失败")
	}

	if site.Project.Title != "" {
		doc.Find("title").SetText(site.Project.Title)
	}
	if site.Project.Icon != "" {
		selection := doc.Find("link[rel=\"icon\"]")
		if selection.Length() > 0 {
			selection.SetAttr("href", site.Project.Icon)
		} else {
			doc.Find("head").AppendHtml(fmt.Sprintf("<link rel=\"icon\" href=\"%s\">", site.Project.Icon))
		}
	}
	if site.Seo.Keywords != "" {
		selection := doc.Find("meta[name=\"keywords\"]")
		if selection.Length() > 0 {
			selection.SetAttr("content", site.Seo.Keywords)
		} else {
			doc.Find("head").AppendHtml(fmt.Sprintf("<meta name=\"keywords\" content=\"%s\">", site.Seo.Keywords))
		}
	}
	if site.Seo.Description != "" {
		selection := doc.Find("meta[name=\"description\"]")
		if selection.Length() > 0 {
			selection.SetAttr("content", site.Seo.Description)
		} else {
			doc.Find("head").AppendHtml(fmt.Sprintf("<meta name=\"description\" content=\"%s\">", site.Seo.Description))
		}
	}

	html, err := doc.Html()
	if err != nil {
		slog.Error("生成站点 html 失败", "err", err, "path", site.Project.WebPath)
		return errors.New("生成html失败")
	}

	err = os.WriteFile(site.Project.WebPath, []byte(html), 0666)
	if err != nil {
		slog.Error("写入站点文件失败", "err", err, "path", site.Project.WebPath)
		return errors.New("文件写入失败")
	}
	return nil
}
