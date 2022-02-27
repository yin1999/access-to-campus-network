package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/ini.v1"
)

var (
	username    string
	password    string
	networkType service
	count       uint
)

func main() {
exit:
	for {
		client := newClient()
		err := client.check()
		switch err {
		case ErrNeedLogin:
			err = client.login(username, password, networkType)
			if err == nil {
				log.Print("[Info] 成功登录到校园网")
				break exit
			}
			log.Printf("[Error] 登录校园网失败，err: %s\n", err.Error())
		case nil:
			log.Print("[Info] 已连接到互联网")
			break exit
		default:
			log.Println(err.Error())
		}
		if count == 0 {
			exit(2)
		}
		count--
		log.Println("[Info] 10秒后重试...")
		time.Sleep(10 * time.Second)
	}
	exit(0)
}

func exit(code int) {
	if len(os.Args) == 1 {
		fmt.Println("5秒后自动关闭窗口...")
		time.Sleep(5 * time.Second)
	}
	os.Exit(code)
}

func init() {
	var net string
	if len(os.Args) == 1 {
		file := filepath.Join(filepath.Dir(os.Args[0]), "config.ini")
		cfg, err := ini.Load(file)
		if err != nil {
			log.Printf("[Error] 加载配置文件失败, err: %s\n", err.Error())
			exit(1)
		}
		sec := cfg.Section("")
		username = sec.Key("user").String()
		password = sec.Key("password").String()
		count = sec.Key("count").MustUint(5)
		net = sec.Key("net").MustString("out-campus")

	} else {
		flagSet := flag.NewFlagSet("hohai campus network access", flag.ExitOnError)
		flagSet.StringVar(&username, "u", "", "设置登录用户名")
		flagSet.StringVar(&password, "p", "", "设置登录密码")
		flagSet.UintVar(&count, "c", 5, "设置重试次数")
		flagSet.StringVar(&net, "net", "out-campus", `设置网络提供商，支持：[ out-campus | cmcc ]
out-campus  校园外网服务(out-campus NET)
cmcc        中国移动(CMCC NET)`)
		flagSet.Parse(os.Args[1:])
	}
	if username == "" || password == "" {
		log.Printf("[Error] 用户名或密码为空")
		exit(1)
	}
	switch net {
	case "out-campus":
		networkType = outCampus
	case "cmcc":
		networkType = cmcc
	default:
		log.Printf("[Error] -net %s 不支持的网络提供商\n", net)
		exit(1)
	}
}
