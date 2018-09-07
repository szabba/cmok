package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/szabba/cmok"
	"github.com/szabba/cmok/fs"
	"github.com/szabba/cmok/userlist"
)

var (
	addr       string
	storageDir string
)

func main() {
	flag.StringVar(&addr, "addr", ":1157", "the address to listen on")
	flag.StringVar(&storageDir, "storage-dir", ".", "the directory to use for storage")

	flag.Parse()

	authCfg := userlist.Config{
		Users: map[cmok.User]userlist.UserConfig{
			"uploader": {Password: "download"},
		},
	}
	authSvc := userlist.NewAuthService(authCfg)
	storage := fs.NewStorage(storageDir)
	handler := cmok.NewHandler(authSvc, storage)

	log.Printf("listening on %q", addr)
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal(err)
	}
}
