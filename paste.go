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
	"sort";
	"strconv";
	"strings";
	"template";
	"time";
	"./pretty";
	. "./html";
)

// Used for building a URL response for POST requests to /
const DOMAIN = "gopaste.org"

// Location for pastes
const PATH = "pastes/"

// Pastes per page at /all
const PER_PAGE = 5

// Sort paste files by modification date, not name
type pasteList []*os.Dir

func (d pasteList) Len() int		{ return len(d) }
func (d pasteList) Less(i, j int) bool	{ return d[i].Mtime_ns > d[j].Mtime_ns }
func (d pasteList) Swap(i, j int)	{ d[i], d[j] = d[j], d[i] }


// Templates
var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": urlHtmlFormatter,
	"code": codePrinter,
	"code-truncated": truncatedCodePrinter,
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
							Li("Mod+S to submit."),
							Li(
								Input().Attrs(As{
									"type": "checkbox",
									"class": "paste-private",
									"name": "private",
								}),
								"Private")).Attrs(As{
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
			"{@|code}").Attrs(As{
			"id": "view",
		}))).Out()
var viewTempl = template.MustParse(viewStr, fmap)

type allEnv struct {
	prev	string;
	next	string;
	pastes	[]string;
}

var allStr = "<!DOCTYPE html>" + Html(
	Head(
		Title("All | Go Paste!"),
		Link().Attrs(As{
			"rel": "stylesheet",
			"href": "css",
			"type": "text/css",
			"media": "screen",
			"charset": "utf-8",
		})),
	Body(
		Div(
			H1("All Public Pastes"),
			"{.repeated section pastes}",
			H2(
				"Paste ",
				A("#{@}").Attrs(As{
					"href": "/{@|url+html}",
				})),
			"{@|code-truncated}",
			"{.end}",

			"{.section prev}",
			A("&larr; Prev").Attrs(As{
				"class": "page prev-page",
				"href": "{prev}",
			}),
			"{.end}",

			"{.section next}",
			A("Next &rarr;").Attrs(As{
				"class": "page next-page",
				"href": "{next}",
			}),
			"{.end}").Attrs(As{
			"id": "all",
		}))).Out()
var allTempl = template.MustParse(allStr, fmap)


// Actions
func home(c *http.Conn, req *http.Request) {
	_, path := path.Split(req.URL.Path);

	switch {
	case req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0:
		paste := savePaste(req.FormValue("code"), req.FormValue("private") != "");
		c.SetHeader("Content-type", "text/plain; charset=utf-8");
		c.Write(strings.Bytes("http://" + DOMAIN + "/" + paste + "\n"));
	case len(path) > 0:
		viewTempl.Execute(path, c)
	default:
		homeTempl.Execute(nil, c)
	}
}

func add(c *http.Conn, req *http.Request) {
	fmt.Println(req.FormValue("private"));
	if req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0 {
		paste := savePaste(req.FormValue("code"), req.FormValue("private") != "");
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

func all(c *http.Conn, req *http.Request) {
	files, ok := io.ReadDir(PATH);
	sort.Sort(pasteList(files));

	if ok != nil {
		c.Write(strings.Bytes("Error reading pastes directory.\n"));
		return;
	}

	page := 1;
	if req.FormValue("page") != "" {
		page, _ = strconv.Atoi(req.FormValue("page"))
	}

	offset := PER_PAGE * (page - 1);

	limit := len(files) - offset;
	if limit > PER_PAGE {
		limit = PER_PAGE
	}

	if limit < 0 {
		c.Write(strings.Bytes("Page too far.\n"));
		return;
	}

	pastes := make([]string, limit);
	for i, j := 0, offset; j < offset+limit; j++ {
		if strings.HasPrefix(files[j].Name, "private:") {
			limit++;
			continue;
		}

		pastes[i] = files[j].Name;
		i++;
	}

	prev := "";
	if page > 1 {
		prev = "/all?page=" + strconv.Itoa(page-1)
	}

	next := "";
	if len(files)/PER_PAGE > page {
		next = "/all?page=" + strconv.Itoa(page+1)
	}
	allTempl.Execute(allEnv{prev, next, pastes}, c);
}

func css(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.css") }

func jquery(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "jquery.js") }

func js(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.js") }


// Template filters
func urlHtmlFormatter(w io.Writer, v interface{}, _ string) {
	template.HTMLEscape(w, strings.Bytes(http.URLEscape(v.(string))))
}

func codePrinter(w io.Writer, v interface{}, _ string) {
	codeLines(w, v.(string), 0)
}

func codeLines(w io.Writer, paste string, limit int) {
	source, ok := io.ReadFile("pastes" + path.Clean("/"+paste));

	if ok != nil {
		fmt.Fprintf(w, "Could not read paste:\n\t%s\n", ok);
		return;
	}

	prettyCode, ok := pretty.Print(paste, string(source));
	if ok != nil {	// If it fails to parse, just serve it raw.
		prettyCode = string(source)
	}

	linesPre := Pre().Attrs(As{"class": "line-numbers"});
	codePre := Pre().Attrs(As{"class": "code-lines"});

	stopped := false;
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

		if limit != 0 && i == limit {
			stopped = true;
			break;
		}
	}

	if stopped {
		linesPre.Append("\n");
		codePre.Append(
			Div(
				A("\n...").Attrs(As{
					"href": "/{@|url+html}",
					"class": "go-comment",
				})).Attrs(As{
				"class": "line",
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

func truncatedCodePrinter(w io.Writer, v interface{}, _ string) {
	codeLines(w, v.(string), 10)
}


// Paste util functions
func savePaste(source string, private bool) string {
	var paste string;

	if private {
		paste = newName("private:")
	} else {
		paste = newName("")
	}

	io.WriteFile(PATH+paste, strings.Bytes(source), 0644);

	return paste;
}

func newName(prefix string) (file string) {
	file = prefix + randomString(5);
	_, err := os.Open(PATH+file, os.O_RDONLY, 0);

	if err != nil {
		return
	}

	return newName(prefix);
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


// Runtime
var addr = flag.String("addr", ":8000", "http service address")

func main() {
	rand.Seed(time.Nanoseconds());

	flag.Parse();

	http.Handle("/", http.HandlerFunc(home));
	http.Handle("/add", http.HandlerFunc(add));
	http.Handle("/raw", http.HandlerFunc(raw));
	http.Handle("/view", http.HandlerFunc(view));
	http.Handle("/all", http.HandlerFunc(all));
	http.Handle("/css", http.HandlerFunc(css));
	http.Handle("/jquery", http.HandlerFunc(jquery));
	http.Handle("/js", http.HandlerFunc(js));

	err := http.ListenAndServe(*addr, nil);
	if err != nil {
		log.Exit("ListenAndServe:", err)
	}
}
