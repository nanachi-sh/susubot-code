package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/internal/configs"
)

func Serve() {
	addr := fmt.Sprintf("0.0.0.0:%d", configs.HTTP_LISTEN_PORT)
	fmt.Printf("Starting http file server at %s...\n", addr)
	log.Fatalln(http.ListenAndServe(addr, http.FileServer(http.Dir(configs.WebDir))))
}
