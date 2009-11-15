package main

import (
    "flag";
    "http";
    "io";
    "log";
    "strings";
    "template";
    . "./html";
)

var addr = flag.String("addr", ":1718", "http service address") // Q=17, R=18

var fmap = template.FormatterMap{
    "html": template.HTMLFormatter,
    "url+html": UrlHtmlFormatter,
}

var templateStr = Html(
    Head(Title("QR Link Generator")),
    Body(
        "{.section @}",
        Img().Attrs(A { "src": "http://chart.apis.google.com/chart?chs=300x300&cht=qr&choe=UTF-8&chl={@|url+html}" }),
        Br(),
        Br(),
        "{.end}",
        Form(
            Input().Attrs(A { "maxLength": "1024",
                              "size": "70",
                              "name": "s",
                              "value": "{@|html}",
                              "title": "Text to QR Encode" }),
            Input().Attrs(A { "type": "submit",
                              "value": "Show QR",
                              "name": "qr" })
        ).Attrs(A { "action": "/",
                    "name": "f",
                    "method": "GET" })
    )
).Out();
var templ = template.MustParse(templateStr, fmap)

func main() {
    flag.Parse();
    http.Handle("/", http.HandlerFunc(QR));
    err := http.ListenAndServe(*addr, nil);
    if err != nil {
        log.Exit("ListenAndServe:", err);
    }
}

func QR(c *http.Conn, req *http.Request) {
    templ.Execute(req.FormValue("s"), c);
}

func UrlHtmlFormatter(w io.Writer, v interface{}, fmt string) {
    template.HTMLEscape(w, strings.Bytes(http.URLEscape(v.(string))));
}


