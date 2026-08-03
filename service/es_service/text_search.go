package es_service

import (
	"bytes"
	"clover_server/global"
	"clover_server/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gorm.io/plugin/dbresolver"
)

// SearchTextsRequest 文本搜索请求。
type SearchTextsRequest struct {
	Keyword string
	Page    int
	Limit   int
}

// TextSearchItem 文本搜索结果。
type TextSearchItem struct {
	ArticleID uint   `json:"articleID"`
	Head      string `json:"head"`
	Body      string `json:"body"`
	Flag      string `json:"flag"`
}

// SearchTextsResponse 文本搜索响应。
type SearchTextsResponse struct {
	List       []TextSearchItem `json:"list"`
	Count      int64            `json:"count"`
	Source     string           `json:"source"`
	Degraded   bool             `json:"degraded"`
	DegradeMsg string           `json:"degradeMsg,omitempty"`
}

// SearchTexts 文本搜索，优先走 ES，ES 不可用时降级到 MySQL。
func SearchTexts(req SearchTextsRequest) (*SearchTextsResponse, error) {
	req = normalizeTextSearchRequest(req)
	if global.DB == nil {
		return nil, fmt.Errorf("db 未初始化")
	}
	if global.ES == nil {
		return searchTextsByMySQL(req, "es 未初始化")
	}
	response, err := searchTextsByES(req)
	if err == nil {
		return response, nil
	}
	return searchTextsByMySQL(req, err.Error())
}

func normalizeTextSearchRequest(req SearchTextsRequest) SearchTextsRequest {
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

func searchTextsByMySQL(req SearchTextsRequest, degradeMsg string) (*SearchTextsResponse, error) {
	db := global.DB.Clauses(dbresolver.Write).Model(&models.TextModel{})
	if strings.TrimSpace(req.Keyword) != "" {
		like := "%" + req.Keyword + "%"
		db = db.Where("head LIKE ? OR body LIKE ?", like, like)
	}
	var count int64
	if err := db.Count(&count).Error; err != nil {
		return nil, fmt.Errorf("mysql 文本搜索统计失败: %w", err)
	}
	var list []models.TextModel
	if err := db.Order("id desc").Offset((req.Page-1)*req.Limit).Limit(req.Limit).Find(&list).Error; err != nil {
		return nil, fmt.Errorf("mysql 文本搜索列表失败: %w", err)
	}
	items := make([]TextSearchItem, 0, len(list))
	for _, item := range list {
		items = append(items, TextSearchItem{ArticleID: item.ArticleID, Head: item.Head, Body: item.Body, Flag: item.Head})
	}
	return &SearchTextsResponse{List: items, Count: count, Source: "mysql", Degraded: true, DegradeMsg: degradeMsg}, nil
}

func searchTextsByES(req SearchTextsRequest) (*SearchTextsResponse, error) {
	body, err := buildTextESSearchBody(req)
	if err != nil {
		return nil, err
	}
	response, err := global.ES.Search(
		global.ES.Search.WithContext(context.Background()),
		global.ES.Search.WithIndex(TextIndex),
		global.ES.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return nil, fmt.Errorf("es 文本搜索失败: %w", err)
	}
	defer response.Body.Close()
	if response.IsError() {
		return nil, fmt.Errorf("es 文本搜索失败: status=%s", response.Status())
	}
	var result esTextSearchResponse
	if err = json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("解析 es 文本搜索结果失败: %w", err)
	}
	items := make([]TextSearchItem, 0, len(result.Hits.Hits))
	for _, hit := range result.Hits.Hits {
		item := TextSearchItem{ArticleID: hit.Source.ArticleID, Head: hit.Source.Head, Body: hit.Source.Body, Flag: hit.Source.Head}
		if len(hit.Highlight.Head) > 0 {
			item.Head = hit.Highlight.Head[0]
		}
		if len(hit.Highlight.Body) > 0 {
			item.Body = hit.Highlight.Body[0]
		}
		items = append(items, item)
	}
	return &SearchTextsResponse{List: items, Count: result.Hits.Total.Value, Source: "elasticsearch", Degraded: false}, nil
}

func buildTextESSearchBody(req SearchTextsRequest) ([]byte, error) {
	mustList := make([]map[string]any, 0, 1)
	if strings.TrimSpace(req.Keyword) != "" {
		mustList = append(mustList, map[string]any{
			"multi_match": map[string]any{
				"query":  req.Keyword,
				"fields": []string{"head", "body"},
			},
		})
	} else {
		mustList = append(mustList, map[string]any{"match_all": map[string]any{}})
	}
	body := map[string]any{
		"from": (req.Page - 1) * req.Limit,
		"size": req.Limit,
		"query": map[string]any{
			"bool": map[string]any{
				"must": mustList,
			},
		},
		"highlight": map[string]any{
			"fields": map[string]any{
				"head": map[string]any{},
				"body": map[string]any{},
			},
		},
	}
	byteData, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("构建 es 文本搜索条件失败: %w", err)
	}
	return byteData, nil
}

type esTextSearchResponse struct {
	Hits struct {
		Total struct {
			Value int64 `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source struct {
				ArticleID uint   `json:"article_id"`
				Head      string `json:"head"`
				Body      string `json:"body"`
			} `json:"_source"`
			Highlight struct {
				Head []string `json:"head"`
				Body []string `json:"body"`
			} `json:"highlight"`
		} `json:"hits"`
	} `json:"hits"`
}
