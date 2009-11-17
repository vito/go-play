package main

import (
	"flag";
	"fmt";
	"path";
	"http";
	"io";
	"os";
	"log";
	"rand";
	"strings";
	"template";
	"time";
	"./pretty";
	. "./html";
)

// Used for building a URL response for POST requests to /
const DOMAIN = "gopaste.org"


var addr = flag.String("addr", ":8000", "http service address")

var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": UrlHtmlFormatter,
	"pretty": CodePrinter,
}

var homeStr = "<!DOCTYPE html>" + Html(
	Head(
		Title("Go Paste!"),
		Link().Attrs(As{
			"rel": "stylesheet",
			"href": "/css",
			"type": "text/css",
			"media": "screen",
			"charset": "utf-8",
		}),
		Script("").Attrs(As{
			"src": "/jquery",
			"type": "text/javascript",
			"charset": "utf-8",
		}),
		Script("").Attrs(As{
			"src": "/js",
			"type": "text/javascript",
			"charset": "utf-8",
		})),
	Body(
		Div(
			Form(
				Fieldset(
					P(
						Textarea("").Attrs(As{
							"class": "paste-input",
							"rows": "30",
							"name": "code",
						}),
						Ul(
							Li("Tab key inserts tabstops."),
							Li("Mod+S to submit.")	//,
							/*Li(*/
							/*Input().Attrs(As{*/
							/*"type": "checkbox",*/
							/*"class": "paste-private",*/
							/*"name": "private"*/
							/*}),*/
							/*"Private"*/
							/*)*/
						).Attrs(As{
							"class": "paste-notes",
						})),
					Input().Attrs(As{
						"class": "paste-submit",
						"type": "submit",
						"value": "Go Paste!",
						"accesskey": "s",
					}))).Attrs(As{
				"action": "/add",
				"method": "POST",
			})).Attrs(As{
			"id": "home",
		}))).Out()
var homeTempl = template.MustParse(homeStr, fmap)

var viewStr = "<!DOCTYPE html>" + Html(
	Head(
		Title("Paste #{@|html} | Go Paste!"),
		Link().Attrs(As{
			"rel": "stylesheet",
			"href": "css",
			"type": "text/css",
			"media": "screen",
			"charset": "utf-8",
		})),
	Body(
		Div(
			A("raw").Attrs(As{
				"href": "/raw?paste={@|url+html}",
				"class": "raw",
			}),
			"{@|pretty}").Attrs(As{
			"id": "view",
		}))).Out()
var viewTempl = template.MustParse(viewStr, fmap)

func css(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.css") }

func jquery(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "jquery.js") }

func js(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.js") }

func home(c *http.Conn, req *http.Request) {
	_, path := path.Split(req.URL.Path);

	switch {
	case req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0:
		paste := savePaste(req.FormValue("code"));
		c.SetHeader("Content-type", "text/plain; charset=utf-8");
		c.Write(strings.Bytes("http://" + DOMAIN + "/" + paste + "\n"));
	case len(path) > 0:
		viewTempl.Execute(path, c);
	default:
		homeTempl.Execute(nil, c)
	}
}

func add(c *http.Conn, req *http.Request) {
	if req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0 {
		paste := savePaste(req.FormValue("code"));
		c.SetHeader("Location", "/"+paste);
		c.WriteHeader(http.StatusFound);
	} else {
		c.Write(strings.Bytes("No code submitted.\n"))
	}
}

func raw(c *http.Conn, req *http.Request) {
	if len(req.FormValue("paste")) > 0 {
		http.ServeFile(c, req, "pastes"+path.Clean("/"+req.FormValue("paste")))
	} else {
		c.Write(strings.Bytes("No paste specified.\n"))
	}
}

func view(c *http.Conn, req *http.Request) {
	if len(req.FormValue("paste")) > 0 {
		viewTempl.Execute(req.FormValue("paste"), c)
	} else {
		c.Write(strings.Bytes("No paste specified.\n"))
	}
}

func UrlHtmlFormatter(w io.Writer, v interface{}, _ string) {
	template.HTMLEscape(w, strings.Bytes(http.URLEscape(v.(string))))
}

func CodePrinter(w io.Writer, v interface{}, _ string) {
	source, ok := io.ReadFile("pastes" + path.Clean("/"+v.(string)));

	if ok != nil {
		fmt.Fprintf(w, "Could not read paste:\n\t%s\n", ok);
		return;
	}

	prettyCode := pretty.Print(v.(string), string(source));

	linesPre := Pre().Attrs(As{"class": "line-numbers"});
	codePre := Pre().Attrs(As{"class": "code-lines"});

	for i, code := range strings.Split(prettyCode, "\n", 0) {
		line := i + 1;
		linesPre.Append(
			fmt.Sprintf(
				A("%d").Attrs(As{
					"rel": "#L%d",
					"href": "#LC%d",
					"id": "LID%d",
				}).Out()+"\n",
				line, line, line, line));
		codePre.Append(
			Div(code).Attrs(As{
				"class": "line",
				"id": "LC" + fmt.Sprint(line),
			}).Out());
	}
	fmt.Fprint(
		w,
		Table(
			Tbody(
				Tr(
					Td(linesPre).Attrs(As{"width": "1%", "valign": "top"}),
					Td(codePre).Attrs(As{"valign": "top"})))).Attrs(As{
			"class": "code",
			"cellspacing": "0",
			"cellpadding": "0",
		}).Out());
}

func savePaste(source string) string {
	paste := newName();
	io.WriteFile("pastes/"+paste, strings.Bytes(source), 0644);
	return paste;
}

func newName() (name string) {
	name = randomString(5);
	_, err := os.Open("pastes/" + name, os.O_RDONLY, 0);

	if err != nil {
		return;
	}

	return newName();
}

func randomString(length int) string {
	var rng, offset, mode int;

	str := make([]int, length);
	for i := 0; i < length; i++ {
		mode = rand.Intn(3);

		if mode == 0 {
			rng, offset = 26, 65
		} else if mode == 1 {
			rng, offset = 26, 97
		} else {
			rng, offset = 10, 48
		}

		str[i] = rand.Intn(rng) + offset;
	}

	return string(str);
}

func main() {
	rand.Seed(time.Nanoseconds());

	flag.Parse();

	http.Handle("/", http.HandlerFunc(home));
	http.Handle("/add", http.HandlerFunc(add));
	http.Handle("/raw", http.HandlerFunc(raw));
	http.Handle("/view", http.HandlerFunc(view));
	http.Handle("/css", http.HandlerFunc(css));
	http.Handle("/jquery", http.HandlerFunc(jquery));
	http.Handle("/js", http.HandlerFunc(js));

	err := http.ListenAndServe(*addr, nil);
	if err != nil {
		log.Exit("ListenAndServe:", err)
	}
}
