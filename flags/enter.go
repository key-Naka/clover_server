package flags

import (
	"clover_server/flags/flag_user"
	"clover_server/global"
	"flag"
	"fmt"
	"os"
)

type Options struct {
	File    string
	DB      bool
	Version bool
	Type    string
	Sub     string
}

var FlagsOptions = new(Options)

func ParseFlags() {
	flag.StringVar(&FlagsOptions.File, "file", "settings.yaml", "配置文件")
	flag.BoolVar(&FlagsOptions.DB, "db", false, "数据库迁移")
	flag.BoolVar(&FlagsOptions.Version, "version", false, "版本")
	flag.StringVar(&FlagsOptions.Type, "type", "server", "运行类型")
	flag.StringVar(&FlagsOptions.Sub, "sub", "", "子命令")
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
	switch FlagsOptions.Type {
	case "user":
		u := flag_user.FlagUser{}
		switch FlagsOptions.Sub {
		case "create":
			u.Create()
			os.Exit(0)
		default:
			fmt.Println("未知运行类型")
		}
	case "es":
		switch FlagsOptions.Sub {
		case "index":
			EsIndex()
			os.Exit(0)
		case "sync":
			EsSync()
			os.Exit(0)
		default:
			fmt.Println("未知运行类型")
		}
	}

}
