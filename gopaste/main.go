package main

import (
	"flag";
	"http";
	"log";
	"rand";
	"time";
	"./gopaste";
)

// Controller
var cont = gopaste.New()

var addr = flag.String("addr", ":8000", "http service address")

func main() {
	rand.Seed(time.Nanoseconds());

	flag.Parse();

	http.Handle("/", http.HandlerFunc(cont.Handler()));

	err := http.ListenAndServe(*addr, nil);
	if err != nil {
		log.Exit("ListenAndServe:", err)
	}
}
