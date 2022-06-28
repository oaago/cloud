package main

import (
	"fmt"
	"github.com/oaago/component/config"
	_ "github.com/oaago/component/preload"
)

func main() {
	fmt.Println("demo", config.Op)
	select {}
}
