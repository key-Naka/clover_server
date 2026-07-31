package article_api

import (
	"clover_server/global"
	"clover_server/models"
	"clover_server/service/es_service"
	"strings"
	"unicode/utf8"
)

func buildArticleAbstract(content string) string {
	text := strings.TrimSpace(content)
	if text == "" {
		return ""
	}
	if utf8.RuneCountInString(text) <= 100 {
		return text
	}
	runes := []rune(text)
	return string(runes[:100])
}

func queryCategoryTitle(categoryID *uint) string {
	if categoryID == nil {
		return ""
	}
	var category models.CategoryModel
	if err := global.DB.Take(&category, *categoryID).Error; err != nil {
		return ""
	}
	return category.Title
}

func syncArticleAfterWrite(articleID uint) error {
	if global.ES == nil {
		return nil
	}
	return es_service.SyncArticleDocument(articleID)
}
