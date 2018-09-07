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
	storageDir string
)

func main() {
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

	err := http.ListenAndServe(":1157", handler)
	if err != nil {
		log.Fatal(err)
	}
}
