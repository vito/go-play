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
    "time";
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
        Link().Attrs(As{
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
                Textarea("").Attrs(As{
                    "cols": "100",
                    "rows": "15",
                    "name": "code"
                }),
                Br(),
                Input().Attrs(As{
                    "type": "submit",
                    "value": "Go Paste!"
                })).Attrs(As{
                "action": "/add",
                "name": "f",
                "method": "POST",
            })).Attrs(As{
            "id": "home"
            }))).Out()
var homeTempl = template.MustParse(homeStr, fmap)

var viewStr = Html(
    Head(
        Title("Pasted!"),
        Link().Attrs(As{
            "rel": "stylesheet",
            "href": "css",
            "type": "text/css",
            "media": "screen",
            "charset": "utf-8"
        })
    ),
    Body(
        Div(
            H1(
                "Paste: <a href=\"/view?paste={@|url+html}\">#{@|html}</a>"
            ),
            "{@|pretty}"
        ).Attrs(As{
            "id": "view"
        }))).Out()
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
    http.ServeFile(c, req, "paste.css");
}

func home(c *http.Conn, req *http.Request)	{
    homeTempl.Execute(nil, c)
}

func add(c *http.Conn, req *http.Request)	{
    if req.Method == "POST" {
        paste := savePaste(req.FormValue("code"));
        c.SetHeader("Location", "/view?paste=" + paste);
        c.WriteHeader(http.StatusFound);
    }
}

func view(c *http.Conn, req *http.Request)	{
    // Set the method to GET so redirects from /add will parse the URL.
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

    prettyCode := pretty.Print(v.(string), string(source));

    linesPre := Pre().Attrs(As{"class": "line_numbers"});
    codePre := Pre();

    for i, code := range strings.Split(prettyCode, "\n", 0) {
        line := i + 1;
        linesPre.Append(
            fmt.Sprintf(
                A("%d").Attrs(As{
                    "rel": "#L%d",
                    "href": "#L%d",
                    "id": "LID%d"
                }).Out() + "\n",
                line, line, line, line
            )
        );
        codePre.Append(
            Div(code).Attrs(As{
                "class": "line",
                "id": "LC" + fmt.Sprint(line)
            }).Out()
        );
    }
    fmt.Fprintf(
        w,
        Table(
            Tbody(
                Tr(
                    Td(linesPre).Attrs(As{"valign": "top"}),
                    Td(codePre).Attrs(As{"valign": "top"})
                )
            )
        ).Attrs(As{
            "class": "code",
            "cellspacing": "0",
            "cellpadding": "0"
        }).Out()
    );
}

func savePaste(source string) string {
    paste := randomString(32);
    io.WriteFile("pastes/" + paste, strings.Bytes(source), 0644);
    return paste;
}

func randomString(length int) string {
    rand.Seed(time.Nanoseconds());

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

