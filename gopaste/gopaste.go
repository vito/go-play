package gopaste

import (
	"fmt";
	"http";
	"io";
	"os";
	"rand";
	"sort";
	"strconv";
	"strings";
	"./gopaste_view";
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


func Home(c *http.Conn, req *http.Request) {
	if req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0 {
		paste := savePaste(req.FormValue("code"), req.FormValue("private") != "");
		c.SetHeader("Content-type", "text/plain; charset=utf-8");
		c.Write(strings.Bytes("http://" + DOMAIN + "/view/" + paste + "\n"));
	} else {
		gopaste_view.Home.Execute(nil, c)
	}
}

func View(c *http.Conn, _ *http.Request, id string) {
	gopaste_view.View.Execute(id, c)
}

func Raw(c *http.Conn, req *http.Request, id string) {
	http.ServeFile(c, req, "pastes/"+id)
}

func Add(c *http.Conn, req *http.Request) {
	if req.Method == "POST" && len(strings.TrimSpace(req.FormValue("code"))) > 0 {
		paste := savePaste(req.FormValue("code"), req.FormValue("private") != "");
		c.SetHeader("Location", "/view/"+paste);
		c.WriteHeader(http.StatusFound);
	} else {
		c.Write(strings.Bytes("No code submitted.\n"))
	}
}

func All(c *http.Conn, req *http.Request)	{ AllPaged(c, req, 1) }

func AllPaged(c *http.Conn, req *http.Request, page int) {
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

	gopaste_view.All.Execute(gopaste_view.AllEnv{prev, next, pastes}, c);
}


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


func Test(c *http.Conn, req *http.Request, arg string, page int) {
	fmt.Printf("Got test request: %s, %d\n", arg, page)
}

func Css(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.css") }

func JQuery(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "jquery.js") }

func Js(c *http.Conn, req *http.Request)	{ http.ServeFile(c, req, "paste.js") }
