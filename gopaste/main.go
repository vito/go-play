package main

import (
	"flag";
	"http";
	"log";
	"rand";
	"runtime";
	"time";
)

// Controller
var cont = gopaste()

var addr = flag.String("addr", ":8000", "http service address")

func main() {
	runtime.GOMAXPROCS(8);

	rand.Seed(time.Nanoseconds());

	flag.Parse();

	http.Handle("/", http.HandlerFunc(cont.Handler()));

	err := http.ListenAndServe(*addr, nil);
	if err != nil {
		log.Exit("ListenAndServe:", err)
	}
}
