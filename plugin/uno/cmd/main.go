package main

import (
	"github.com/nanachi-sh/susubot-code/plugin/uno/cmd/api"
	"github.com/nanachi-sh/susubot-code/plugin/uno/log"
)

var logger = log.Get()

func main() {
	go func() {
		if err := api.GRPCServe(); err != nil {
			logger.Fatalln(err)
		}
	}()
	go func() {
		if err := api.HTTPServe(); err != nil {
			logger.Fatalln(err)
		}
	}()
	select {}
}
