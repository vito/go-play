package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"
)

// Controller
var cont = gopaste()

var addr = flag.String("addr", ":8000", "http service address")

func main() {
	runtime.GOMAXPROCS(8)

  rand.Seed(time.Now().Unix())

	flag.Parse()

	http.Handle("/", cont.HandlerFunc())

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
