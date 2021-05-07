package build

import (
	"github.com/micro/go-micro/v2"
	"github.com/yiqiang3344/go-lib/utils/config"
	cNet "github.com/yiqiang3344/go-lib/utils/net"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func CheckSignReload(project string, service micro.Service) {
	ip, _ := cNet.GetLocalIP()
	//监听USR1信号
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGUSR1)
	for {
		sig := <-ch
		log.Printf("signal: %v", sig)

		switch sig {
		case syscall.SIGUSR1:
			//重新拉起服务
			shell := "./main --server=" + project + " &"
			command := exec.Command("bash", "-c", shell)
			if err := command.Run(); err != nil {
				log.Println(err)
				return
			}
			log.Println("start other process[" + strconv.Itoa(command.Process.Pid) + "] success: " + shell)

			//注销服务
			shell = "micro --registry=etcd " +
				"--registry_address=" +
				config.GetCfgString("etcd.address") +
				" deregister service '{" +
				"\"name\": \"" + service.Server().Options().Name + "\"," +
				"\"version\": \"" + service.Server().Options().Version + "\"," +
				"\"nodes\": [{" +
				"\"id\": \"" + service.Server().Options().Name + "-" + service.Server().Options().Id + "\"," +
				"\"address\": \"" + ip + "\"," +
				"\"port\": " + strings.Split(service.Server().Options().Address, "]:")[1] +
				"}]}'"
			command = exec.Command("bash", "-c", shell)
			if err := command.Run(); err != nil {
				log.Println(err)
				return
			}
			log.Println("Deregister success: " + shell)

			//60秒之后结束此进程
			time.Sleep(60 * time.Second)

			shell = "kill -15 " + strconv.Itoa(os.Getpid())
			command = exec.Command("bash", "-c", shell)
			if err := command.Run(); err != nil {
				log.Println(err)
				return
			}
			log.Println("kill pid success: " + shell)
			return
		}
	}
}
