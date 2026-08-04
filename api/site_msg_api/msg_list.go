package site_msg_api

import (
	"clover_server/common"
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"clover_server/models/enum/relationship_enum"
	"clover_server/utils/jwts"

	"github.com/gin-gonic/gin"
)

type SiteMsgListRequest struct {
	common.PageInfo
	T int8 `form:"t" binding:"required,oneof=1 2 3"`
}

type SiteMsgListResponse struct {
	models.MessageModel
	Relation relationship_enum.Relation `json:"relation"`
}

func (SiteMsgApi) SiteMsgListView(c *gin.Context) {
	var req SiteMsgListRequest
	if err := c.ShouldBindQuery(&req); err != nil { res.FailWithError(err, c); return }
	claims := jwts.GetClaims(c)
	if claims == nil { res.FailWithMsg("请登录", c); return }
	typeList := getSiteMessageTypes(req.T)
	query := global.DB.Model(&models.MessageModel{}).Where("rev_user_id = ? and type in ?", claims.UserID, typeList)
	var count int64
	if err := query.Count(&count).Error; err != nil { res.FailWithError(err, c); return }
	var list []models.MessageModel
	if err := query.Order("id desc").Offset(req.GetOffset()).Limit(req.GetLimit()).Find(&list).Error; err != nil { res.FailWithError(err, c); return }
	resp := make([]SiteMsgListResponse, 0, len(list))
	for _, item := range list { resp = append(resp, SiteMsgListResponse{MessageModel: item, Relation: relationship_enum.RelationStranger}) }
	res.OkWithList(resp, int(count), c)
}

func getSiteMessageTypes(t int8) []enum.MessageType {
	switch t {
	case 1:
		return []enum.MessageType{enum.MessageCommentType, enum.MessageCommentReplyType}
	case 2:
		return []enum.MessageType{enum.MessageArticleLikeType, enum.MessageCommentLikeType, enum.MessageArticleCollectType}
	default:
		return []enum.MessageType{enum.MessageSystemNotificationType}
	}
}
