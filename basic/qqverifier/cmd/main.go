package main

import (
	"github.com/nanachi-sh/susubot-code/basic/qqverifier/cmd/api"
	"github.com/nanachi-sh/susubot-code/basic/qqverifier/log"
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
