package flags

import (
	"clover_server/global"
	"flag"
	"fmt"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
}

var FlagsOptions = new(Options)

func ParseFlags() {
	flag.StringVar(&FlagsOptions.File, "file", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagsOptions.DB, "db", false, "数据库迁移")
	flag.BoolVar(&FlagsOptions.Version, "version", false, "版本")
	flag.Parse()
}

func Run() {
	if FlagsOptions.Version {
		fmt.Println(global.Version)
		os.Exit(0)
	}

	if FlagsOptions.DB {
		FlagDB()
		os.Exit(0)
	}
}
