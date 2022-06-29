package main

import (
	"fmt"
	"github.com/oaago/cloud/config"
	_ "github.com/oaago/cloud/preload"
)

func main() {
	fmt.Println("demo", config.Op)
	select {}
}
