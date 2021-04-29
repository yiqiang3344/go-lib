package helper

import (
	"fmt"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/config/source/file"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
)

var Cfg config.Config

var prefixs = []string{"cfg"}

func InitCfg() {
	err := config.Load(
		etcd.NewSource(
			etcd.WithAddress("127.0.0.1:2379"),
			etcd.WithPrefix("/common/config"),
			etcd.StripPrefix(true),
		),
		etcd.NewSource(
			etcd.WithAddress("127.0.0.1:2379"),
			etcd.WithPrefix("/xyf-robot-srv/config"),
			etcd.StripPrefix(true),
		),
		file.NewSource(file.WithPath("config/"+os.Getenv("ENV")+".yaml")),
	)
	if err != nil {
		log.Fatalf("log init failed:", err.Error())
	}

	go hotUpdate()
}

func hotUpdate() {
	for {
		//启动热更新
		w, err := config.Watch("cfg")
		if err != nil {
			ErrorLog("配置热更新失败:"+err.Error(), RunFuncName())
			return
		}
		v, err := w.Next()
		if err != nil {
			ErrorLog("配置热更新下一个失败:"+err.Error(), RunFuncName())
			return
		}
		before := fmt.Sprint(config.Get(prefixs...).StringMap(map[string]string{}))
		err = config.Sync()
		if err != nil {
			ErrorLog("配置热更新同步失败:"+err.Error(), RunFuncName())
			return
		}
		DebugLog(
			"配置热更新成功",
			RunFuncName(),
			zap.String("beforeCfgData", before),
			zap.String("afterCfgData", fmt.Sprint(v.StringMap(map[string]string{}))),
		)
	}
}

func GetCfgString(route string) string {
	return getCfg(route).String("")
}

func GetCfgInt(route string) int {
	return getCfg(route).Int(0)
}

func parseRoute(route string) []string {
	arr := prefixs[0:]
	for _, v := range strings.Split(route, ".") {
		arr = append(arr, v)
	}
	return arr
}

func getCfg(route string) reader.Value {
	arr := parseRoute(route)
	return config.Get(arr...)
}
