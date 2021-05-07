package config

import (
	"fmt"
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/reader"
	"github.com/micro/go-micro/v2/config/source/etcd"
	"github.com/micro/go-micro/v2/config/source/file"
	cLog "github.com/yiqiang3344/go-lib/utils/log"
	"github.com/yiqiang3344/go-lib/utils/trace"
	"go.uber.org/zap"
	"log"
	"os"
	"strings"
	"time"
)

var Cfg config.Config

var prefixs = []string{"cfg"}

func InitCfg() {
	local, err := config.NewConfig(
		config.WithSource(file.NewSource(file.WithPath("config/common.yaml"))),
		config.WithSource(file.NewSource(file.WithPath("config/"+os.Getenv("ENV")+".yaml"))),
	)
	if err != nil {
		log.Fatalf("local log init failed:", err.Error())
	}
	//本地配置配置中心的地址信息
	etcdAddress := GetCfg(local, "etcdCfgCenter.address").String("")
	//多数据源，越下面优先级越高
	err = config.Load(
		etcd.NewSource(
			etcd.WithAddress(etcdAddress),
			etcd.WithPrefix("/common/config"),
			etcd.StripPrefix(true),
		),
		etcd.NewSource(
			etcd.WithAddress(etcdAddress),
			etcd.WithPrefix("/"+GetCfg(local, "project").String("")+"/config"),
			etcd.StripPrefix(true),
		),
		file.NewSource(file.WithPath("config/common.yaml")),
		file.NewSource(file.WithPath("config/"+os.Getenv("ENV")+".yaml")),
	)
	_ = local.Close()
	if err != nil {
		log.Fatalf("log init failed:", err.Error())
	}

	//热更新配置
	go hotUpdate()
}

func hotUpdate() {
	time.Sleep(5 * time.Second) //暂停5秒之后再执行，等日志初始化完毕，否则可能会写入日志失败
	for {
		//启动热更新
		w, err := config.Watch("cfg")
		if err != nil {
			cLog.ErrorLog("配置热更新失败:"+err.Error(), trace.RunFuncName())
			return
		}
		v, err := w.Next()
		if err != nil {
			cLog.ErrorLog("配置热更新下一个失败:"+err.Error(), trace.RunFuncName())
			return
		}
		before := fmt.Sprint(config.Get(prefixs...).StringMap(map[string]string{}))
		err = config.Sync()
		if err != nil {
			cLog.ErrorLog("配置热更新同步失败:"+err.Error(), trace.RunFuncName())
			return
		}
		cLog.DebugLog(
			"配置热更新成功",
			trace.RunFuncName(),
			zap.String("beforeCfgData", before),
			zap.String("afterCfgData", fmt.Sprint(v.StringMap(map[string]string{}))),
		)
	}
}

func GetCfgString(route string) string {
	return GetCfg(config.DefaultConfig, route).String("")
}

func GetCfgInt(route string) int {
	return GetCfg(config.DefaultConfig, route).Int(0)
}

func GetCfgStringMap(route string) map[string]string {
	return GetCfg(config.DefaultConfig, route).StringMap(map[string]string{})
}

func parseRoute(route string) []string {
	arr := prefixs[0:]
	for _, v := range strings.Split(route, ".") {
		arr = append(arr, v)
	}
	return arr
}

func GetCfg(cfg config.Config, route string) reader.Value {
	arr := parseRoute(route)
	return cfg.Get(arr...)
}
