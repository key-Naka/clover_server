package main

import (
	"clover_server/core"
	"fmt"
)

func main() {
	// 传入配置文件路径，并处理可能返回的错误
	cfg, err := core.ReadConf("settings.yaml")
	if err != nil {
		panic(err)
	}

	fmt.Println("成功加载配置，IP:", cfg.System.IP)
}
