package main

import "github.com/nanachi-sh/susubot-code/plugin/randomanimal/LLOneBot/cmd/api"

func main() {
	go func() {
		if err := api.GRPCServe(); err != nil {
			panic(err)
		}
	}()
	select {}
}