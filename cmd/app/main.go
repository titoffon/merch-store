package main

import (
	"fmt"

	"github.com/titoffon/merch-store/internal/config"
	"github.com/titoffon/merch-store/internal/server"
)

func main(){

	cfg := config.LoadConfig()
	
	err := server.Run(cfg)
	if err != nil {
		fmt.Println()
		return
	}
}