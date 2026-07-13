// common/list_query.go
package common

import (
	"clover_server/global"
	"fmt"

	"gorm.io/gorm"
)

type PageInfo struct {
	Limit int    `form:"limit"`
	Page  int    `form:"page"`
	Key   string `form:"key"`
	Order string `form:"order"`
}

func (p PageInfo) GetPage() int {
	if p.Page > 20 || p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p PageInfo) GetLimit() int {
	if p.Limit <= 0 || p.Limit > 100 {
		return 10
	}
	return p.Limit
}
func (p PageInfo) GetOffset() int {
	return (p.GetPage() - 1) * p.GetLimit()
}

type Options struct {
	PageInfo     PageInfo
	Likes        []string
	Preloads     []string
	Where        *gorm.DB
	Debug        bool
	DefaultOrder string
}

func ListQuery[T any](model T, option Options) (list []T, count int, err error) {

	query := global.DB.Model(model).Where(model)

	if option.Debug {
		query = query.Debug()
	}

	if len(option.Likes) > 0 && option.PageInfo.Key != "" {
		likes := global.DB.Where(fmt.Sprintf("%s like ?", option.Likes[0]), fmt.Sprintf("%%%s%%", option.PageInfo.Key))
		for _, column := range option.Likes[1:] {
			likes = likes.Or(
				fmt.Sprintf("%s like ?", column),
				fmt.Sprintf("%%%s%%", option.PageInfo.Key),
			)
		}
		query = query.Where(likes)
	}

	if option.Where != nil {
		query = query.Where(option.Where)
	}

	for _, preload := range option.Preloads {
		query = query.Preload(preload)
	}

	var _c int64
	query.Count(&_c)
	count = int(_c)
	limit := option.PageInfo.GetLimit()
	offset := option.PageInfo.GetOffset()

	if option.PageInfo.Order != "" {
		query = query.Order(option.PageInfo.Order)
	} else {
		if option.DefaultOrder != "" {
			query = query.Order(option.DefaultOrder)
		}
	}

	err = query.Offset(offset).Limit(limit).Find(&list).Error
	return
}
