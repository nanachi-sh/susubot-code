package main

import (
	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/cmd/api"
	"github.com/nanachi-sh/susubot-code/basic/handler/LLOneBot/log"
)

var logger = log.Get()

func main() {
	go func() {
		if err := api.GRPCServe(); err != nil {
			logger.Fatalln(err)
		}
	}()
	select {}
}
