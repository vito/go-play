package main

import (
	"flag";
	"http";
	"log";
	"rand";
	"time";
	"./controller";
	"./gopaste";
)

// Controller
var cont = controller.New(map[string]interface{}{
	`/$`: gopaste.Home,
	`/add`: gopaste.Add,
	`/all`: gopaste.All,
	`/all/page/([0-9]+)`: gopaste.AllPaged,
	`/view/([a-zA-Z0-9:]+)$`: gopaste.View,
	`/raw/([a-zA-Z0-9:]+)$`: gopaste.Raw,
	`/css`: gopaste.Css,
	`/jquery`: gopaste.JQuery,
	`/js`: gopaste.Js,
})


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
