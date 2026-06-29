package core

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var searcher *xdb.Searcher

// InitIPDB 初始化 IP 数据库
func InitIPDB(ip string) *xdb.Searcher {
	var dbPath = "init/ip2region.xdb"

	cBuff, err := xdb.LoadContentFromFile(dbPath)
	if err != nil {
		slog.Error("ip地址数据库文件加载错误", "err", err, "path", dbPath)
		return nil
	}

	_searcher, err := xdb.NewWithBuffer(xdb.IPv4, cBuff)
	if err != nil {
		slog.Error("基于内存的 searcher 创建错误", "err", err)
		return nil
	}
	searcher = _searcher
	slog.Info("ip2region 初始化成功")
	return searcher
}

// SearchAddr 查询 IP 地址
func SearchAddr(ip string) (string, error) {
	region, err := searcher.Search(ip)
	if err != nil {
		slog.Error("ip地址查询错误", "err", err)
		return "", err
	}

	_addrList := strings.Split(region, "|")
	if len(_addrList) != 5 {
		slog.Warn("异常ip地址数据", "region", region)
		return "", errors.New("invalid region format")
	}
	country := _addrList[0]
	province := _addrList[2]
	city := _addrList[3]

	if province != "0" && city != "0" {
		return fmt.Sprintf("%s|%s", province, city), nil
	}
	if province != "0" {
		return province, nil
	}

	if country != "0" && country != "" {
		return country, nil
	}

	return "未知位置", nil
}
