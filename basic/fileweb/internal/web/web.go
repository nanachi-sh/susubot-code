package web

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nanachi-sh/susubot-code/basic/fileweb/internal/configs"
)

var (
	LOAD struct{}
)

func init() {
	fmt.Printf("Starting http file server at %d...\n", configs.HTTP_LISTEN_PORT)
	log.Fatalln(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", configs.HTTP_LISTEN_PORT), http.FileServer(http.Dir(configs.WebDir))))
}
