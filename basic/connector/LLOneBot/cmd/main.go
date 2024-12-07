package main

import (
	"github.com/nanachi-sh/susubot-code/basic/connector/LLOneBot/cmd/api"
	"github.com/nanachi-sh/susubot-code/basic/connector/LLOneBot/log"
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
