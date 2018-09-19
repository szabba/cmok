package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/szabba/cmok"
	"github.com/szabba/cmok/config"
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

	config, err := configuration.Parse(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	authSvc := userlist.NewAuthService(config.AuthConfig)
	accessPolicy := userlist.NewAccessPolicy(config.AccessConfig)

	storage := fs.NewStorage(storageDir)

	// TODO: move muxing out of handler
	handler := cmok.NewHandler("", authSvc, accessPolicy, storage)

	log.Printf("listening on %q", addr)
	// TODO: non-global, properly configured server
	err = http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal(err)
	}
}
