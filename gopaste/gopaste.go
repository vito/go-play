package main

import (
	"fmt";
	"http";
	"io";
	"os";
	"rand";
	"regexp";
	"sort";
	"strconv";
	"strings";

	"./controller";
	. "./html";
)


// Used for building a URL response for POST requests to /
const DOMAIN = "gopaste.org"

// Location for pastes
const PATH = "pastes/"

// Pastes per page at /all
const PER_PAGE = 15

// Default theme
const THEME = "twilight"

// All themes available
var THEMES = map[string]string{
	"twilight": "Twilight",
	"clean": "Clean (Pastie)",
	"slate": "Slate",
	"vibrant_ink": "Vibrant Ink",
	"white_black": "White on Black",
	"black_white": "Black on White",
}


func gopaste() *controller.Controller {
	cont := controller.New();

	cont.AddHandler(`/$`, home);
	cont.AddHandler(`/add`, add);
	cont.AddHandler(`/all/page/([0-9]+)`, allPaged);
	cont.AddHandler(`/all`, all);
	cont.AddHandler(`/view/([a-zA-Z0-9:]+)$`, view);
	cont.AddHandler(`/raw/([a-zA-Z0-9:]+)$`, raw);
	cont.AddHandler(`/css/([a-z_]+)`, cssFile);
	cont.AddHandler(`/css`, css);
	cont.AddHandler(`/jquery`, jQuery);
	cont.AddHandler(`/js`, js);
	cont.AddHandler(`/([a-zA-Z0-9:]+)$`, view);

	return cont;
}

func curTheme(req *http.Request) (theme string) {
	theme = THEME;

	if cookie, ok := req.Header["Cookie"]; ok {
		match := regexp.MustCompile(`theme=([a-z_]+)`).MatchStrings(cookie);
		theme = match[1];
	}

	return;
}

func home(c *http.Conn, req *http.Request) {
	if req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0 {
		paste := savePaste(req.FormValue("code"), req.FormValue("private") != "");
		c.SetHeader("Content-type", "text/plain; charset=utf-8");
		c.Write(strings.Bytes("http://" + DOMAIN + "/view/" + paste + "\n"));
	} else {
		theme := curTheme(req);
		homePage.Execute(map[string]string{
			"theme": theme,
			"theme_select": themeSelect(theme),
		}, c)
	}
}

func view(c *http.Conn, req *http.Request, id string) {
	theme := curTheme(req);
	viewPage.Execute(map[string]string{
		"id": id,
		"theme": theme,
		"theme_select": themeSelect(theme)
	}, c)
}

func raw(c *http.Conn, req *http.Request, id string) {
	http.ServeFile(c, req, "pastes/"+id)
}

func add(c *http.Conn, req *http.Request) {
	if req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0 {
		paste := savePaste(req.FormValue("code"), req.FormValue("private") != "");
		c.SetHeader("Location", "/view/"+paste);
		c.WriteHeader(http.StatusFound);
	} else {
		c.Write(strings.Bytes("No code submitted.\n"))
	}
}

func all(c *http.Conn, req *http.Request)	{ allPaged(c, req, 1) }

func allPaged(c *http.Conn, req *http.Request, page int) {
	files, ok := io.ReadDir(PATH);
	sort.Sort(pasteList(files));

	if ok != nil {
		c.Write(strings.Bytes("Error reading pastes directory.\n"));
		return;
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
		prev = "/all/page/" + strconv.Itoa(page-1)
	}

	next := "";
	if len(files)/PER_PAGE > page {
		next = "/all/page/" + strconv.Itoa(page+1)
	}

	codeList := make([]string, len(pastes));
	results := make(chan int);
	for i, paste := range pastes {
		go func(i int, paste string) {
			code, err := prettyPaste(paste, 10);
			if err != nil {
				code[0] = err.String()
			}

			codeList[i] =
				fmt.Sprintf(
					H2(
						"Paste ",
						A("#%s").Attrs(As{
							"href": "/view/%s",
						})).Out(),
					paste, paste) + code[0];

			results <- i;
		}(i, paste)
	}

	for i := 0; i < len(pastes); i++ {
		<-results;
	}

	theme := curTheme(req);
	allPage.Execute(allEnv{prev, next, codeList, theme, themeSelect(theme)}, c);
}


func css(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.css") }
func cssFile(c *http.Conn, req *http.Request, file string) {
	http.ServeFile(c, req, "css/" + file + ".css");
}

func jQuery(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "jquery.js") }

func js(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.js") }


// Sort paste files by modification date, not name
type pasteList []*os.Dir

func (d pasteList) Len() int		{ return len(d) }
func (d pasteList) Less(i, j int) bool	{ return d[i].Mtime_ns > d[j].Mtime_ns }
func (d pasteList) Swap(i, j int)	{ d[i], d[j] = d[j], d[i] }


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
