package log_service

import (
	"clover_server/global"
	"clover_server/models"
	"clover_server/models/enum"
	"encoding/json"
	"fmt"
	"log/slog"
	"reflect"
	"strings"
	"time"
)

type RuntimeLog struct {
	level       enum.LogLevelType
	title       string
	itemList    []string
	serviceName string
	dateType    RuntimeDateType
}

func (r *RuntimeLog) Save() {
	r.SetNowTime()
	var log models.LogModel
	global.DB.Find(&log, fmt.Sprintf("service_name = ? and log_type = ? and create_time >= date_sub(now(), %s)", r.dateType.GetSqlTime()), r.serviceName, enum.RuntimeLogType)
	if log.ID != 0 {
		c := strings.Join(r.itemList, "\n")
		newContent := log.Content + "\n" + c
		global.DB.Model(&log).Updates(map[string]any{
			"content": newContent,
		})
		r.itemList = []string{}
		return
	}
	log = models.LogModel{
		LogType:     enum.RuntimeLogType,
		Title:       r.title,
		Content:     strings.Join(r.itemList, "\n"),
		Level:       r.level,
		ServiceName: r.serviceName,
	}
	err := global.DB.Create(&log).Error
	if err != nil {
		slog.Error("保存运行时日志失败", "err", err)
	}
}
func (ac *RuntimeLog) SetTitle(title string) {
	ac.title = title
}

func (ac *RuntimeLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}
func (ac *RuntimeLog) SetLink(label string, href string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item link\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\"><a href=\"%s\" target=\"_blank\">%s</a></div></div> ",
		label,
		href, href))
}
func (ac *RuntimeLog) SetImage(src string) {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_image\"><img src=\"%s\" alt=\"\"></div>", src))
}
func (ac *RuntimeLog) setItem(label string, value any, logLevelType enum.LogLevelType) {
	var v string

	t := reflect.TypeOf(value)
	switch t.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		byteData, _ := json.Marshal(value)
		v = string(byteData)
	default:
		v = fmt.Sprintf("%v", value)
	}

	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_item %s\"><div class=\"log_item_label\">%s</div><div class=\"log_item_content\">%s</div></div>",
		logLevelType,
		label, v))
}
func (ac *RuntimeLog) SetItem(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}
func (ac *RuntimeLog) SetItemInfo(label string, value any) {
	ac.setItem(label, value, enum.LogInfoLevel)
}
func (ac *RuntimeLog) SetItemWarn(label string, value any) {
	ac.setItem(label, value, enum.LogWarnLevel)
}
func (ac *RuntimeLog) SetItemError(label string, value any) {
	ac.setItem(label, value, enum.LogErrorLevel)
}
func (ac *RuntimeLog) SetError(label string, err error) {
	slog.Error(label, "err", err)
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_error\"><div class=\"line\"><div class=\"label\">%s</div><div class=\"value\">%s</div><div class=\"type\">%T</div></div><div class=\"stack\">%+v</div></div>",
		label, err, err, err))
}
func (ac *RuntimeLog) SetNowTime() {
	ac.itemList = append(ac.itemList, fmt.Sprintf("<div class=\"log_time\">%s</div>", time.Now().Format("2006-01-02 15:04:05")))
}

type RuntimeDateType int8

const (
	RuntimeDateHour      RuntimeDateType = 1
	RuntimeDateTypeDay   RuntimeDateType = 2
	RuntimeDateTypeWeek  RuntimeDateType = 3
	RuntimeDateTypeMonth RuntimeDateType = 4
)

func (r RuntimeDateType) GetSqlTime() string {
	switch r {
	case RuntimeDateHour:
		return "interval 1 HOUR"
	case RuntimeDateTypeDay:
		return "interval 1 DAY"
	case RuntimeDateTypeWeek:
		return "interval 1 WEEK"
	case RuntimeDateTypeMonth:
		return "interval 1 MONTH"
	}
	return "interval 1 DAY"
}
func NewRuntimeLog(serviceName string, dateType RuntimeDateType) *RuntimeLog {
	return &RuntimeLog{serviceName: serviceName, dateType: dateType}
}
