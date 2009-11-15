package main

import (
	"flag";
    "fmt";
    "path";
	"http";
	"io";
	"log";
    "rand";
	"strings";
	"template";
    "./pretty";
	. "./html";
)

var addr = flag.String("addr", ":8000", "http service address")

var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
    "pretty": CodePrinter
}

var homeStr = Html(
	Head(
        Title("Go Paste!"),
        Link().Attrs(A {
            "rel": "stylesheet",
            "href": "css",
            "type": "text/css",
            "media": "screen",
            "charset": "utf-8"
        })
    ),
	Body(
        Div(
            H1("Go Paste!"),
            Form(
                Textarea("").Attrs(A{
                    "cols": "100",
                    "rows": "15",
                    "name": "code"
                }),
                Br(),
                Input().Attrs(A{
                    "type": "submit",
                    "value": "Go Paste!"
                })).Attrs(A{
                "action": "/add",
                "name": "f",
                "method": "POST",
            })).Attrs(A{
            "id": "home"
            }))).Out()
var homeTempl = template.MustParse(homeStr, fmap)

var viewStr = Html(
    Head(
        Title("Pasted!"),
        Link().Attrs(A {
            "rel": "stylesheet",
            "href": "css",
            "type": "text/css",
            "media": "screen",
            "charset": "utf-8"
        })
    ),
    Body(
        H1(
            "Paste: <a href=\"/view?paste={@|url+html}\">#{@|html}</a>"
        ),
        Pre("{@|pretty}")
    )).Out()
var viewTempl = template.MustParse(viewStr, fmap)

func main() {
	flag.Parse();
    http.Handle("/", http.HandlerFunc(home));
	http.Handle("/add", http.HandlerFunc(add));
	http.Handle("/view", http.HandlerFunc(view));
	http.Handle("/css", http.HandlerFunc(css));
	err := http.ListenAndServe(*addr, nil);
	if err != nil {
		log.Exit("ListenAndServe:", err)
	}
}

func css(c *http.Conn, req *http.Request) {
    css, _ := io.ReadFile("paste.css");
    c.Write(css);
}
func home(c *http.Conn, req *http.Request)	{ homeTempl.Execute(nil, c) }
func add(c *http.Conn, req *http.Request)	{
    if req.Method == "POST" {
        paste := savePaste(req.FormValue("code"));
        viewTempl.Execute(paste, c);
    }
}
func view(c *http.Conn, req *http.Request)	{
    if len(req.FormValue("paste")) > 0 {
        viewTempl.Execute(req.FormValue("paste"), c);
    } else {
        c.Write(strings.Bytes("No paste specified.\n"));
    }
}

func UrlHtmlFormatter(w io.Writer, v interface{}, _ string) {
	template.HTMLEscape(w, strings.Bytes(http.URLEscape(v.(string))))
}

func CodePrinter(w io.Writer, v interface {}, _ string) {
    source, ok := io.ReadFile("pastes" + path.Clean("/" + v.(string)));

    if ok != nil {
        fmt.Fprintf(w, "Could not read paste:\n\t%s\n", ok);
        return;
    }

    pretty.Print(w, "", string(source));
}

func savePaste(source string) string {
    paste := randomString(32);
    io.WriteFile("pastes/" + paste, strings.Bytes(source), 0644);
    return paste;
}

func randomString(length int) string {
    var rng, offset int;

    str := make([]int, length);
    for i := 0; i < length; i++ {
        if i % 3 == 0 {
            rng, offset = 26, 65;
        } else if i % 2 == 0 {
            rng, offset = 26, 97;
        } else {
            rng, offset = 10, 48;
        }

        str[i] = rand.Intn(rng) + offset;
    }

    return string(str);
}
