package uuid

import (
	"fmt"
	"github.com/sony/sonyflake"
	"log"
	"time"
)

func GenUuId() string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		log.Fatalf("flake.NextID() failed with %s\n", err)
	}
	return time.Now().Format("20060102150405") + fmt.Sprintf("%x", id)
}
