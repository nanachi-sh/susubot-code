package main

import "github.com/nanachi-sh/susubot-code/plugin/randomfortune/cmd/api"

func main() {
	go func() {
		if err := api.GRPCServe(); err != nil {
			panic(err)
		}
	}()
	select {}
}
