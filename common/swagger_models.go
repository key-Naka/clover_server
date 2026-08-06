package common

// APIResponse 是 Swagger 文档使用的统一响应结构。
type APIResponse struct {
	Code int    `json:"code" example:"0"`
	Data any    `json:"data"`
	Msg  string `json:"msg" example:"成功"`
}

// ArticleListItemDoc 是文章列表项目的文档模型。
type ArticleListItemDoc struct {
	ID            uint     `json:"id"`
	Title         string   `json:"title"`
	Abstract      string   `json:"abstract"`
	CategoryID    *uint    `json:"categoryID"`
	CategoryTitle string   `json:"categoryTitle"`
	TagList       []string `json:"tagList"`
	Cover         string   `json:"cover"`
	UserID        uint     `json:"userID"`
	UserNickname  string   `json:"userNickname"`
	UserAvatar    string   `json:"userAvatar"`
	LookCount     int      `json:"lookCount"`
	DiggCount     int      `json:"diggCount"`
	CommentCount  int      `json:"commentCount"`
	CollectCount  int      `json:"collectCount"`
	OpenComment   bool     `json:"openComment"`
	Status        int8     `json:"status"`
}

// ArticleListDataDoc 是文章列表的分页数据。
type ArticleListDataDoc struct {
	List  []ArticleListItemDoc `json:"list"`
	Count int                  `json:"count"`
}

// UploadResponse 是图片上传成功后的数据结构。
type UploadResponse struct {
	Path     string `json:"path"`
	WebPath  string `json:"web_path"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
}
