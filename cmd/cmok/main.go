package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/szabba/cmok"
	"github.com/szabba/cmok/fs"
)

var (
	storageDir string
)

func main() {
	flag.StringVar(&storageDir, "storage-dir", ".", "the directory to use for storage")

	flag.Parse()

	storage := fs.NewStorage(storageDir)
	handler := cmok.NewHandler(storage)

	err := http.ListenAndServe(":1157", handler)
	if err != nil {
		log.Fatal(err)
	}
}
