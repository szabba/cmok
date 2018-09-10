package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/szabba/cmok"
	"github.com/szabba/cmok/auth"
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

	authCfg := userlist.AuthConfig{
		Users: map[auth.User]userlist.UserConfig{
			"ci":  {Password: "pass"},
			"dev": {Password: "pass"},
		},
	}

	accessConfig := userlist.AccessConfig{
		Permissions: map[auth.User]userlist.Permissions{
			"ci":  userlist.All(),
			"dev": userlist.Read(),
		},
	}

	authSvc := userlist.NewAuthService(authCfg)
	accessPolicy := userlist.NewAccessPolicy(accessConfig)

	storage := fs.NewStorage(storageDir)

	handler := cmok.NewHandler("", authSvc, accessPolicy, storage)

	log.Printf("listening on %q", addr)
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal(err)
	}
}
