package main

import (
	"http";
	"io";
	"strings";
	"template";

	. "./html";
)


type allEnv struct {
	prev		string;
	next		string;
	pastes		[]string;
	theme		string;
	theme_select	string;
}

var fmap = template.FormatterMap{
	"html": template.HTMLFormatter,
	"url+html": urlHtmlFormatter,
	"code": codePrinter,
}


func themeSelect(theme string) string {
	sel := Select().Attrs(As{
		"name": "theme",
	});

	for id, name := range THEMES {
		if id == theme {
			sel.Append(Option(name).Attrs(As{
				"value": id,
				"selected": "selected",
			}))
		} else {
			sel.Append(Option(name).Attrs(As{
				"value": id,
			}))
		}
	}

	return sel.Out();
}

func page(title string, contents *Element) string {
	head := Head(
		Title(title),
		Link().Attrs(As{
			"rel": "stylesheet",
			"href": "/css",
			"type": "text/css",
			"media": "screen",
			"charset": "utf-8",
		}),
		Link().Attrs(As{
			"rel": "stylesheet",
			"title": "{theme}",
			"href": "/css/{theme}",
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
		}));

	for id, _ := range THEMES {
		head.Append(Link().Attrs(As{
			"rel": "alternate stylesheet",
			"title": id,
			"href": "/css/" + id,
			"type": "text/css",
			"media": "screen",
			"charset": "utf-8",
		}))
	}

	return "<!DOCTYPE html>" + Html(
		head,
		Body(
			Form(
				Fieldset(
					"{theme_select}")).Attrs(As{
				"class": "theme-select",
			}),
			contents)).Out();
}

var homePage = template.MustParse(
	page(
		"Go Paste!",
		Div(
			Form(
				Fieldset(
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
					}),
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
		"Paste #{id|html} | Go Paste!",
		Div(
			Ul(
				Li(
					A("raw").Attrs(As{
						"href": "/raw/{id|url+html}",
						"class": "raw",
					})),
				Li(
					A("home").Attrs(As{
						"href": "/",
					}))).Attrs(As{
				"class": "view-links",
			}),
			"{id|code}").Attrs(As{
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

	if len(code) > 1 {
		io.WriteString(w, `<div class="multi-paste">`)
	}

	for _, part := range code {
		io.WriteString(w, part)
	}

	if len(code) > 1 {
		io.WriteString(w, `</div>`)
	}
}
