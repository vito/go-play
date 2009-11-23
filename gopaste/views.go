package main

import (
	"http";
	"io";
	"strings";
	"template";

	. "./html";
)


type allEnv struct {
	prev	string;
	next	string;
	pastes	[]string;
}

var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": urlHtmlFormatter,
	"code": codePrinter,
}


func page(title string, contents *Element) string {
	return "<!DOCTYPE html>" + Html(
		Head(
			Title(title),
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
		Body(contents)).Out()
}

var homePage = template.MustParse(
	page(
		"Go Paste!",
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
									"href": "http://github.com/vito/go-play/tree/master/gopaste",
								})),
							Li(""),
							Li("Tab key inserts tabstops."),
							Li("Seperate content with ---."),
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
		})),
	fmap)

var viewPage = template.MustParse(
	page(
		"Paste #{@|html} | Go Paste!",
		Div(
			A("raw").Attrs(As{
				"href": "/raw/{@|url+html}",
				"class": "raw",
			}),
			"{@|code}").Attrs(As{
			"id": "view",
		})),
	fmap)

var allPage = template.MustParse(
	page(
		"All | Go Paste!",
		Div(
			"{.repeated section pastes}",
			"{@}",
			"{.end}",

			Div(
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
				"class": "pagination",
			})).Attrs(As{
			"id": "all",
		})),
	fmap)


func urlHtmlFormatter(w io.Writer, v interface{}, _ string) {
	template.HTMLEscape(w, strings.Bytes(http.URLEscape(v.(string))))
}

func codePrinter(w io.Writer, v interface{}, _ string) {
	code, _ := prettyPaste(v.(string), 0);

	io.WriteString(w, code);
}
