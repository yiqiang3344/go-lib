package helper

import (
	"github.com/micro/go-micro/v2/config"
	"os"
)

func InitCfg() {
	_ = config.LoadFile("config/" + os.Getenv("ENV") + ".yaml")
}
