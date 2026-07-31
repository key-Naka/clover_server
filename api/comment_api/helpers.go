package comment_api

import (
	"clover_server/global"
	"clover_server/models"
	"clover_server/service/es_service"
)

func loadArticle(articleID uint) (*models.ArticleModel, error) {
	var article models.ArticleModel
	if err := global.DB.Take(&article, articleID).Error; err != nil {
		return nil, err
	}
	return &article, nil
}

func syncArticleCommentToES(articleID uint) error {
	if global.ES == nil {
		return nil
	}
	return es_service.SyncArticleDocument(articleID)
}

func findRootParentID(comment models.CommentModel) *uint {
	if comment.ParentID == nil {
		return nil
	}
	current := comment
	for current.ParentID != nil {
		var parent models.CommentModel
		if err := global.DB.Take(&parent, *current.ParentID).Error; err != nil {
			return comment.ParentID
		}
		if parent.ParentID == nil {
			return &parent.ID
		}
		current = parent
	}
	return comment.ParentID
}

func collectCommentIDs(rootID uint) ([]uint, error) {
	var list []models.CommentModel
	if err := global.DB.Where("article_id = ?", rootID).Find(&list).Error; err != nil {
		return nil, err
	}
	childMap := map[uint][]uint{}
	for _, item := range list {
		if item.ParentID != nil {
			childMap[*item.ParentID] = append(childMap[*item.ParentID], item.ID)
		}
	}
	result := []uint{rootID}
	queue := []uint{rootID}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		children := childMap[current]
		result = append(result, children...)
		queue = append(queue, children...)
	}
	return result, nil
}
