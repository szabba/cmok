package main

import (
	"log"
	"net/http"

	"github.com/szabba/cmok"
	"github.com/szabba/cmok/fs"
)

func main() {
	storage := fs.NewStorage(".")
	handler := cmok.NewHandler(storage)

	err := http.ListenAndServe(":1157", handler)
	if err != nil {
		log.Fatal(err)
	}
}
