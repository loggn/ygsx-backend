package main

import (
	"fmt"
	"log"
	"ygsx-v2/internal/router"
)

func main() {
	r := router.NewRouter()

	fmt.Println("鸭古生鲜SaaS服务器运行中...")

	if err := r.Run(":8000"); err != nil {
		log.Fatal(err)
	}

}
