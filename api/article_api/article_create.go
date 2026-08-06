package article_api

import (
	"clover_server/common/res"
	"clover_server/global"
	"clover_server/models"
	"clover_server/utils/jwts"
	"strings"

	"github.com/gin-gonic/gin"
)

type ArticleCreateRequest struct {
	Title       string   `json:"title" binding:"required"`
	Abstract    string   `json:"abstract"`
	Content     string   `json:"content" binding:"required"`
	CategoryID  *uint    `json:"categoryID"`
	TagList     []string `json:"tagList"`
	Cover       string   `json:"cover"`
	OpenComment bool     `json:"openComment"`
	Status      int8     `json:"status" binding:"required,oneof=1 2 3 4"`
}

// ArticleCreateView 创建文章。
// @Summary 创建文章
// @Description `title`、`content`、`status` 必填。status：1 草稿、2 审核中、3 已发布、4 已拒绝。
// @Tags 文章
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ArticleCreateRequest true "文章内容"
// @Success 200 {object} res.ArticleResponse "创建后的完整文章"
// @Failure 400 {object} res.ErrorResponse "参数、分类或认证失败"
// @Failure 500 {object} res.ErrorResponse "服务异常"
// @Router /article [post]
func (ArticleApi) ArticleCreateView(c *gin.Context) {
	var req ArticleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		res.FailWithError(err, c)
		return
	}

	claims := jwts.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("请登录", c)
		return
	}

	if req.CategoryID != nil {
		var category models.CategoryModel
		if err := global.DB.Take(&category, "id = ? and user_id = ?", *req.CategoryID, claims.UserID).Error; err != nil {
			res.FailWithMsg("文章分类不存在", c)
			return
		}
	}

	abstract := strings.TrimSpace(req.Abstract)
	if abstract == "" {
		abstract = buildArticleAbstract(req.Content)
	}

	article := models.ArticleModel{
		Title:       req.Title,
		Abstract:    abstract,
		Content:     req.Content,
		CategoryID:  req.CategoryID,
		TagList:     req.TagList,
		Cover:       req.Cover,
		UserID:      claims.UserID,
		OpenComment: req.OpenComment,
		Status:      req.Status,
	}
	if err := global.DB.Create(&article).Error; err != nil {
		res.FailWithError(err, c)
		return
	}

	if err := syncArticleAfterWrite(article.ID); err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	res.OkWithData(article, c)
}
