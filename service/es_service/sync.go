package es_service

import (
	"clover_server/global"
	"clover_server/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	"gorm.io/plugin/dbresolver"
)

const (
	ArticleIndex = "article_index"
	TextIndex    = "text_index"
)

// ArticleDocument ES 中的文章文档结构。
type ArticleDocument struct {
	ID           uint              `json:"id"`
	CreatedAt    any               `json:"created_at"`
	UpdatedAt    any               `json:"updated_at"`
	Title        string            `json:"title"`
	Abstract     string            `json:"abstract"`
	Content      string            `json:"content"`
	CategoryID   *uint             `json:"category_id"`
	TagList      []string          `json:"tag_list"`
	Cover        string            `json:"cover"`
	UserID       uint              `json:"user_id"`
	LookCount    int               `json:"look_count"`
	DiggCount    int               `json:"digg_count"`
	CommentCount int               `json:"comment_count"`
	CollectCount int               `json:"collect_count"`
	OpenComment  bool              `json:"open_comment"`
	Status       int8              `json:"status"`
	Comments     []CommentDocument `json:"comments"`
}

// CommentDocument ES 中的评论结构。
type CommentDocument struct {
	ID           uint  `json:"id"`
	CreatedAt    any   `json:"created_at"`
	UpdatedAt    any   `json:"updated_at"`
	Content      string `json:"content"`
	UserID       uint   `json:"user_id"`
	ArticleID    uint   `json:"article_id"`
	ParentID     *uint  `json:"parent_id"`
	RootParentID *uint  `json:"root_parent_id"`
	DiggCount    int    `json:"digg_count"`
}

// SyncAllArticleDocuments 全量同步文章及其评论到 ES。
func SyncAllArticleDocuments() error {
	if global.ES == nil {
		return fmt.Errorf("es 未初始化")
	}

	var articleList []models.ArticleModel
	if err := global.DB.Clauses(dbresolver.Write).Find(&articleList).Error; err != nil {
		return fmt.Errorf("查询文章失败: %w", err)
	}

	for _, article := range articleList {
		if err := SyncArticleDocument(article.ID); err != nil {
			return err
		}
	}
	return nil
}

// SyncArticleDocument 增量同步单篇文章及其评论到 ES。
func SyncArticleDocument(articleID uint) error {
	if global.ES == nil {
		return fmt.Errorf("es 未初始化")
	}

	document, err := buildArticleDocument(articleID)
	if err != nil {
		return err
	}

	byteData, err := json.Marshal(document)
	if err != nil {
		return fmt.Errorf("序列化文章文档失败: %w", err)
	}

	response, err := global.ES.Index(
		ArticleIndex,
		strings.NewReader(string(byteData)),
		global.ES.Index.WithDocumentID(strconv.FormatUint(uint64(articleID), 10)),
		global.ES.Index.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("写入文章索引失败: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		return fmt.Errorf("写入文章索引失败: status=%s", response.Status())
	}

	slog.Info("文章同步到 ES 成功", slog.Uint64("articleID", uint64(articleID)))
	return nil
}

// DeleteArticleDocument 删除 ES 中的文章文档。
func DeleteArticleDocument(articleID uint) error {
	if global.ES == nil {
		return fmt.Errorf("es 未初始化")
	}

	response, err := global.ES.Delete(
		ArticleIndex,
		strconv.FormatUint(uint64(articleID), 10),
		global.ES.Delete.WithContext(context.Background()),
	)
	if err != nil {
		return fmt.Errorf("删除文章索引失败: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() && response.StatusCode != 404 {
		return fmt.Errorf("删除文章索引失败: status=%s", response.Status())
	}

	slog.Info("文章索引删除成功", slog.Uint64("articleID", uint64(articleID)))
	return nil
}

func buildArticleDocument(articleID uint) (*ArticleDocument, error) {
	var article models.ArticleModel
	if err := global.DB.Clauses(dbresolver.Write).First(&article, articleID).Error; err != nil {
		return nil, fmt.Errorf("查询文章失败: %w", err)
	}

	var commentList []models.CommentModel
	if err := global.DB.Clauses(dbresolver.Write).
		Where("article_id = ?", articleID).
		Order("id asc").
		Find(&commentList).Error; err != nil {
		return nil, fmt.Errorf("查询文章评论失败: %w", err)
	}

	comments := make([]CommentDocument, 0, len(commentList))
	for _, comment := range commentList {
		comments = append(comments, CommentDocument{
			ID:           comment.ID,
			CreatedAt:    comment.CreatedAt,
			UpdatedAt:    comment.UpdatedAt,
			Content:      comment.Content,
			UserID:       comment.UserID,
			ArticleID:    comment.ArticleID,
			ParentID:     comment.ParentID,
			RootParentID: comment.RootParentID,
			DiggCount:    comment.DiggCount,
		})
	}

	return &ArticleDocument{
		ID:           article.ID,
		CreatedAt:    article.CreatedAt,
		UpdatedAt:    article.UpdatedAt,
		Title:        article.Title,
		Abstract:     article.Abstract,
		Content:      article.Content,
		CategoryID:   article.CategoryID,
		TagList:      article.TagList,
		Cover:        article.Cover,
		UserID:       article.UserID,
		LookCount:    article.LookCount,
		DiggCount:    article.DiggCount,
		CommentCount: article.CommentCount,
		CollectCount: article.CollectCount,
		OpenComment:  article.OpenComment,
		Status:       article.Status,
		Comments:     comments,
	}, nil
}
