package main

import (
	"fmt"
	"github.com/SinKingCloud/sinking-cloud/cloud-server/cloud-common/bootstrap"
)

func main() {
	bootstrap.Load()
	fmt.Println("succ")
}
