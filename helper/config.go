package helper

import (
	"github.com/micro/go-micro/v2/config"
	"github.com/micro/go-micro/v2/config/reader"
	"os"
	"strings"
)

func InitCfg() {
	_ = config.LoadFile("config/" + os.Getenv("ENV") + ".yaml")
}

func GetCfgString(route string) string {
	arr, _len := parseRoute(route)
	if _len == 0 {
		return ""
	}
	cfg := getCfg(_len, arr)
	return cfg.String("")
}

func GetCfgInt(route string) int {
	arr, _len := parseRoute(route)
	if _len == 0 {
		return 0
	}
	cfg := getCfg(_len, arr)
	return cfg.Int(0)
}

func parseRoute(route string) ([]string, int) {
	arr := strings.Split(route, ".")
	return arr, len(arr)
}

func getCfg(len int, arr []string) reader.Value {
	var ret reader.Value
	switch len {
	case 1:
		ret = config.Get(arr[0])
	case 2:
		ret = config.Get(arr[0], arr[1])
	case 3:
		ret = config.Get(arr[0], arr[1], arr[2])
	case 4:
		ret = config.Get(arr[0], arr[1], arr[2], arr[3])
	case 5:
		ret = config.Get(arr[0], arr[1], arr[2], arr[3], arr[4])
	case 6:
		ret = config.Get(arr[0], arr[1], arr[2], arr[3], arr[4], arr[5])
	case 7:
		ret = config.Get(arr[0], arr[1], arr[2], arr[3], arr[4], arr[5], arr[6])
	case 8:
		ret = config.Get(arr[0], arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], arr[7])
	case 9:
		ret = config.Get(arr[0], arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], arr[7], arr[8])
	case 10:
		ret = config.Get(arr[0], arr[1], arr[2], arr[3], arr[4], arr[5], arr[6], arr[7], arr[8], arr[9])
	}
	return ret
}
