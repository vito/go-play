package gopaste_view

import (
	"http";
	"fmt";
	"io";
	"path";
	"strings";
	"template";
	. "./html";
	"./pretty";
)


type AllEnv struct {
	Prev	string;
	Next	string;
	Pastes	[]string;
}

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
							Li(
								A("All Pastes").Attrs(As{
									"href": "/all",
								})),
							Li(
								A("Source Code").Attrs(As{
									"href": "http://github.com/vito/go-play",
								})),
							Li(""),
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
var Home = template.MustParse(homeStr, fmap)

var viewStr = "<!DOCTYPE html>" + Html(
	Head(
		Title("Paste #{@|html} | Go Paste!"),
		Link().Attrs(As{
			"rel": "stylesheet",
			"href": "/css",
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
var View = template.MustParse(viewStr, fmap)

var allStr = "<!DOCTYPE html>" + Html(
	Head(
		Title("All | Go Paste!"),
		Link().Attrs(As{
			"rel": "stylesheet",
			"href": "/css",
			"type": "text/css",
			"media": "screen",
			"charset": "utf-8",
		})),
	Body(
		Div(
			"{.repeated section Pastes}",
			H2(
				"Paste ",
				A("#{@}").Attrs(As{
					"href": "/{@|url+html}",
				})),
			"{@|code-truncated}",
			"{.end}",

			Div(
				"{.section Prev}",
				A("&larr; Prev").Attrs(As{
					"class": "page prev-page",
					"href": "{Prev}",
				}),
				"{.end}",

				"{.section Next}",
				A("Next &rarr;").Attrs(As{
					"class": "page next-page",
					"href": "{Next}",
				}),
				"{.end}").Attrs(As{
				"class": "pagination",
			})).Attrs(As{
			"id": "all",
		}))).Out()
var All = template.MustParse(allStr, fmap)


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
