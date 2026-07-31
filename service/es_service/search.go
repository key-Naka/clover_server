package es_service

import (
	"bytes"
	"clover_server/global"
	"clover_server/models"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"gorm.io/plugin/dbresolver"
)

// SearchArticlesRequest 查询文章请求。
type SearchArticlesRequest struct {
	Keyword    string
	CategoryID *uint
	Status     *int8
	Page       int
	Limit      int
}

// SearchArticlesResponse 查询文章响应。
type SearchArticlesResponse struct {
	List       []ArticleDocument `json:"list"`
	Count      int64             `json:"count"`
	Source     string            `json:"source"`
	Degraded   bool              `json:"degraded"`
	DegradeMsg string            `json:"degradeMsg,omitempty"`
}

// SearchArticles 文章搜索，优先走 ES，ES 不可用时降级到 MySQL。
func SearchArticles(req SearchArticlesRequest) (*SearchArticlesResponse, error) {
	normalizedReq := normalizeSearchRequest(req)
	if global.DB == nil {
		return nil, fmt.Errorf("db 未初始化")
	}

	if global.ES == nil {
		return searchArticlesByMySQL(normalizedReq, "es 未初始化")
	}

	response, err := searchArticlesByES(normalizedReq)
	if err == nil {
		return response, nil
	}

	slog.Warn("ES 查询失败，开始降级到 MySQL", slog.String("error", err.Error()), slog.String("keyword", normalizedReq.Keyword))
	return searchArticlesByMySQL(normalizedReq, err.Error())
}

func normalizeSearchRequest(req SearchArticlesRequest) SearchArticlesRequest {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Limit <= 0 {
		req.Limit = 10
	}
	if req.Limit > 100 {
		req.Limit = 100
	}
	return req
}

func searchArticlesByES(req SearchArticlesRequest) (*SearchArticlesResponse, error) {
	body, err := buildESSearchBody(req)
	if err != nil {
		return nil, err
	}

	response, err := global.ES.Search(
		global.ES.Search.WithContext(context.Background()),
		global.ES.Search.WithIndex(ArticleIndex),
		global.ES.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("es 查询失败: %w", err)
	}
	defer response.Body.Close()

	if response.IsError() {
		return nil, fmt.Errorf("es 查询失败: status=%s", response.Status())
	}

	var result esSearchResponse
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 es 查询结果失败: %w", err)
	}

	list := make([]ArticleDocument, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		list = append(list, hit.Source)
	}

	return &SearchArticlesResponse{
		List:     list,
		Count:    result.Hits.Total.Value,
		Source:   "elasticsearch",
		Degraded: false,
	}, nil
}

func searchArticlesByMySQL(req SearchArticlesRequest, degradeMsg string) (*SearchArticlesResponse, error) {
	db := global.DB.Clauses(dbresolver.Write).Model(&models.ArticleModel{})
	if req.Keyword != "" {
		likeKeyword := "%" + req.Keyword + "%"
		db = db.Where("title LIKE ? OR abstract LIKE ? OR content LIKE ?", likeKeyword, likeKeyword, likeKeyword)
	}
	if req.CategoryID != nil {
		db = db.Where("category_id = ?", *req.CategoryID)
	}
	if req.Status != nil {
		db = db.Where("status = ?", *req.Status)
	}

	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("mysql 降级查询统计失败: %w", err)
	}

	var articleList []models.ArticleModel
	if err := db.Order("id desc").Offset((req.Page-1)*req.Limit).Limit(req.Limit).Find(&articleList).Error; err != nil {
		return nil, fmt.Errorf("mysql 降级查询列表失败: %w", err)
	}

	list := make([]ArticleDocument, 0, len(articleList))
	for _, article := range articleList {
		list = append(list, ArticleDocument{
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
			Comments:     nil,
		})
	}

	return &SearchArticlesResponse{
		List:       list,
		Count:      count,
		Source:     "mysql",
		Degraded:   true,
		DegradeMsg: degradeMsg,
	}, nil
}

func buildESSearchBody(req SearchArticlesRequest) ([]byte, error) {
	mustList := make([]map[string]any, 0, 3)
	filterList := make([]map[string]any, 0, 2)

	if strings.TrimSpace(req.Keyword) != "" {
		mustList = append(mustList, map[string]any{
			"multi_match": map[string]any{
				"query":  req.Keyword,
				"fields": []string{"title", "abstract", "content", "comments.content"},
			},
		})
	} else {
		mustList = append(mustList, map[string]any{"match_all": map[string]any{}})
	}

	if req.CategoryID != nil {
		filterList = append(filterList, map[string]any{
			"term": map[string]any{"category_id": *req.CategoryID},
		})
	}
	if req.Status != nil {
		filterList = append(filterList, map[string]any{
			"term": map[string]any{"status": *req.Status},
		})
	}

	query := map[string]any{
		"bool": map[string]any{
			"must":   mustList,
			"filter": filterList,
		},
	}

	body := map[string]any{
		"from":  (req.Page - 1) * req.Limit,
		"size":  req.Limit,
		"query": query,
		"sort": []map[string]any{
			{"_score": map[string]any{"order": "desc"}},
			{"id": map[string]any{"order": "desc"}},
		},
	}

	byteData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("构建 es 查询条件失败: %w", err)
	}
	return byteData, nil
}

type esSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source ArticleDocument `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
