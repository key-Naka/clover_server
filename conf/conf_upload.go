package conf

// Upload 定义文件上传配置。
type Upload struct {
	Size      int      `yaml:"size" json:"size"`           // 大小限制，单位 MB
	WhiteList []string `yaml:"whiteList" json:"whiteList"` // 允许上传的后缀列表，不带点
	UploadDir string   `yaml:"uploadDir" json:"uploadDir"` // 上传目录
}
