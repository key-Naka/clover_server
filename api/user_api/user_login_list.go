package user_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"
	"time"

	"github.com/gin-gonic/gin"
)

type UserLoginListRequest struct {
	common.PageInfo
	UserID    uint   `json:"userID"`
	Ip        string `json:"ip"`
	Addr      string `json:"addr"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	Type      int8   `json:"type"`
}
type UserLoginListResponse struct {
	models.UserLoginModel
	UserNickname string `json:"userNickname,omitempty"`
	UserAvatar   string `json:"userAvatar,omitempty"`
}

func (UserApi) UserLoginListView(c *gin.Context) {
	var cr UserLoginListRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithMsg("参数错误", c)
		return
	}
	claims := jwts.GetClaims(c)
	if cr.Type == 1 {
		cr.UserID = claims.UserID
	}
	var query = global.DB.Where("")
	if cr.StartTime != 0 {
		_, err = time.Parse("2006-01-02 15:04:05", time.Unix(cr.StartTime, 0).Format("2006-01-02 15:04:05"))
		if err != nil {
			res.FailWithMsg("开始时间格式错误", c)
			return
		}
		query.Where("created_at >= ?", cr.StartTime)
	}
	if cr.EndTime != 0 {
		_, err = time.Parse("2006-01-02 15:04:05", time.Unix(cr.EndTime, 0).Format("2006-01-02 15:04:05"))
		if err != nil {
			res.FailWithMsg("结束时间格式错误", c)
			return
		}
		query.Where("created_at <= ?", cr.EndTime)
	}
	var preloads []string
	if cr.Type == 2 {
		preloads = []string{"UserModel"}
	}

	_list, count, _ := common.ListQuery(models.UserLoginModel{
		UserID: cr.UserID,
		Ip:     cr.Ip,
		Addr:   cr.Addr,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Where:    query,
		Preloads: preloads,
	})
	var list = make([]UserLoginListResponse, 0)
	for _, model := range _list {
		list = append(list, UserLoginListResponse{
			UserLoginModel: model,
			UserNickname:   model.UserModel.Nickname,
			UserAvatar:     model.UserModel.Avatar,
		})
	}
	res.OkWithList(list, count, c)
}
