package res

import (
	"clover_server/utils/validate"

	"github.com/gin-gonic/gin"
)

type Code int

const (
	SuccessCode     Code = 0    // 成功
	FailValidCode   Code = 1001 // 校验失败
	FailServiceCode Code = 1002 // 服务异常
)

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "校验失败"
	case FailServiceCode:
		return "服务异常"
	}
	return ""
}

var empty = map[string]any{}

type Response struct {
	Code Code   `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

// SuccessResponse 是无数据成功响应的 Swagger 文档模型。
type SuccessResponse struct {
	Code Code              `json:"code" example:"0"`
	Data map[string]string `json:"data"`
	Msg  string            `json:"msg" example:"成功"`
}

// MessageResponse 是仅返回提示信息的 Swagger 文档模型。
type MessageResponse struct {
	Code Code              `json:"code" example:"0"`
	Data map[string]string `json:"data"`
	Msg  string            `json:"msg" example:"操作成功"`
}

// ErrorResponse 是失败响应的 Swagger 文档模型。
type ErrorResponse struct {
	Code Code              `json:"code" example:"1001"`
	Data map[string]string `json:"data"`
	Msg  string            `json:"msg" example:"参数校验失败"`
}

// TokenResponse 是登录成功响应的 Swagger 文档模型。
type TokenResponse struct {
	Code Code   `json:"code" example:"0"`
	Data string `json:"data" example:"eyJhbGciOiJIUzI1NiIs..."`
	Msg  string `json:"msg" example:"成功"`
}

// EmailCodeData 是邮件验证码标识。
type EmailCodeData struct {
	EmailID string `json:"emailID"`
}

// EmailCodeResponse 是发送邮件验证码成功响应。
type EmailCodeResponse struct {
	Code Code          `json:"code" example:"0"`
	Data EmailCodeData `json:"data"`
	Msg  string        `json:"msg" example:"成功"`
}

// ArticleData 是文章详情及创建成功时返回的文章数据。
type ArticleData struct {
	ID           uint     `json:"id"`
	Title        string   `json:"title"`
	Abstract     string   `json:"abstract"`
	Content      string   `json:"content"`
	CategoryID   *uint    `json:"categoryID"`
	TagList      []string `json:"tagList"`
	Cover        string   `json:"cover"`
	UserID       uint     `json:"userID"`
	LookCount    int      `json:"lookCount"`
	DiggCount    int      `json:"diggCount"`
	CommentCount int      `json:"commentCount"`
	CollectCount int      `json:"collectCount"`
	OpenComment  bool     `json:"openComment"`
	Status       int8     `json:"status"`
}

// ArticleListItem 是文章列表项目。
type ArticleListItem struct {
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

// ArticleListData 是文章分页列表数据。
type ArticleListData struct {
	List  []ArticleListItem `json:"list"`
	Count int               `json:"count"`
}

// ArticleListResponse 是文章列表成功响应。
type ArticleListResponse struct {
	Code Code            `json:"code" example:"0"`
	Data ArticleListData `json:"data"`
	Msg  string          `json:"msg" example:"成功"`
}

// ArticleDetailData 是文章详情数据。
type ArticleDetailData struct {
	ArticleData
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	UserAvatar    string `json:"userAvatar"`
	CategoryTitle string `json:"categoryTitle"`
	IsDigg        bool   `json:"isDigg"`
	IsCollect     bool   `json:"isCollect"`
}

// ArticleDetailResponse 是文章详情成功响应。
type ArticleDetailResponse struct {
	Code Code              `json:"code" example:"0"`
	Data ArticleDetailData `json:"data"`
	Msg  string            `json:"msg" example:"成功"`
}

// ArticleResponse 是创建文章成功响应。
type ArticleResponse struct {
	Code Code        `json:"code" example:"0"`
	Data ArticleData `json:"data"`
	Msg  string      `json:"msg" example:"成功"`
}

// ImageUploadData 是图片上传成功后的数据。
type ImageUploadData struct {
	Path     string `json:"path"`
	WebPath  string `json:"web_path"`
	Filename string `json:"filename"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
}

// ImageUploadResponse 是图片上传成功响应。
type ImageUploadResponse struct {
	Code Code            `json:"code" example:"0"`
	Data ImageUploadData `json:"data"`
	Msg  string          `json:"msg" example:"成功"`
}

// QiNiuTokenData 是七牛直传凭证。
type QiNiuTokenData struct {
	Key    string `json:"key"`
	Token  string `json:"token"`
	Region string `json:"region"`
}

// QiNiuTokenResponse 是获取七牛直传凭证成功响应。
type QiNiuTokenResponse struct {
	Code Code           `json:"code" example:"0"`
	Data QiNiuTokenData `json:"data"`
	Msg  string         `json:"msg" example:"成功"`
}

// ImageListItem 是图片列表项目。
type ImageListItem struct {
	ID       uint   `json:"id"`
	Filename string `json:"filename"`
	Path     string `json:"path"`
	WebPath  string `json:"web_path"`
	Size     int64  `json:"size"`
	Hash     string `json:"hash"`
}

// ImageListData 是图片分页列表数据。
type ImageListData struct {
	List  []ImageListItem `json:"list"`
	Count int             `json:"count"`
}

// ImageListResponse 是图片列表成功响应。
type ImageListResponse struct {
	Code Code          `json:"code" example:"0"`
	Data ImageListData `json:"data"`
	Msg  string        `json:"msg" example:"成功"`
}

func (r Response) Json(c *gin.Context) {
	c.JSON(200, r)
}

func Ok(data any, msg string, c *gin.Context) {
	Response{SuccessCode, data, msg}.Json(c)
}

func OkWithData(data any, c *gin.Context) {
	Response{SuccessCode, data, "成功"}.Json(c)
}

func OkWithList(list any, count int, c *gin.Context) {
	Response{SuccessCode, map[string]any{
		"list":  list,
		"count": count,
	}, "成功"}.Json(c)
}

func OkWithMsg(msg string, c *gin.Context) {
	Response{SuccessCode, empty, msg}.Json(c)
}

func FailWithMsg(msg string, c *gin.Context) {
	Response{FailValidCode, empty, msg}.Json(c)
}

func FailWithData(data any, msg string, c *gin.Context) {
	Response{FailServiceCode, data, msg}.Json(c)
}

func FailWithCode(code Code, c *gin.Context) {
	Response{code, empty, code.String()}.Json(c)
}

func FailWithError(err error, c *gin.Context) {
	data, msg := validate.ValidateError(err)
	FailWithData(data, msg, c)
}
