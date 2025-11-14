package main

import (
	"gee"
	"log"
	"net/http"
)

func main() {
	engine := gee.New()
	log.Fatal(http.ListenAndServe(":9999", engine))
}
