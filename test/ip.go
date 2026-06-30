// testdata/10.获取地理位置.go
package main

import (
	"fmt"
	"time"

	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

func main() {
	ip2region()
}
func ip2region() {
	var dbPath = "init/ip2region.xdb"
	searcher, err := xdb.NewWithFileOnly(xdb.IPv4, dbPath)
	if err != nil {
		fmt.Printf("failed to create searcher: %s\n", err.Error())
		return
	}

	defer searcher.Close()
	var ip = "175.0.201.207"
	var tStart = time.Now()
	region, err := searcher.Search(ip)
	if err != nil {
		fmt.Printf("failed to SearchIP(%s): %s\n", ip, err)
		return
	}
	fmt.Printf("{region: %s, took: %s}\n\n", region, time.Since(tStart))
	// 备注：并发使用，每个 goroutine 需要创建一个独立的 searcher 对象。
}
