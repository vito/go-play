package main

import (
	"bytes"
	"fmt"
	"net/url"
	"text/template"

	. "gopaste/html"
)

type allEnv struct {
	prev         string
	next         string
	pastes       []string
	theme        string
	theme_select string
}

var fmap = template.FuncMap{
	"url_html": urlHtmlFormatter,
	"code":     codePrinter,
}

func templ(name string) *template.Template {
  t := template.New(name)
  t.Funcs(fmap)
  return t
}

func themeSelect(theme string) string {
	sel := Select().Attrs(As{
		"name": "theme",
	})

	for id, name := range THEMES {
		if id == theme {
			sel.Append(Option(name).Attrs(As{
				"value":    id,
				"selected": "selected",
			}))
		} else {
			sel.Append(Option(name).Attrs(As{
				"value": id,
			}))
		}
	}

	return sel.Out()
}

func page(title string, contents *Element) string {
	head := Head(
		Title(title),
		Link().Attrs(As{
			"rel":     "stylesheet",
			"href":    "/css",
			"type":    "text/css",
			"media":   "screen",
			"charset": "utf-8",
		}),
		Link().Attrs(As{
			"rel":     "stylesheet",
			"title":   "{{.theme}}",
			"href":    "/css/{{.theme}}",
			"type":    "text/css",
			"media":   "screen",
			"charset": "utf-8",
		}),
		Script("").Attrs(As{
			"src":     "/jquery",
			"type":    "text/javascript",
			"charset": "utf-8",
		}),
		Script("").Attrs(As{
			"src":     "/js",
			"type":    "text/javascript",
			"charset": "utf-8",
		}))

	for id, _ := range THEMES {
		head.Append(Link().Attrs(As{
			"rel":     "alternate stylesheet",
			"title":   id,
			"href":    "/css/" + id,
			"type":    "text/css",
			"media":   "screen",
			"charset": "utf-8",
		}))
	}

	return "<!DOCTYPE html>" + Html(
		head,
		Body(
			Form(
				Fieldset(
					"{{.theme_select}}")).Attrs(As{
				"class": "theme-select",
			}),
			contents)).Out()
}

var homePage = template.Must(templ("home").Parse(
	page(
		"Go Paste!",
		Div(
			Form(
				Fieldset(
					Textarea("").Attrs(As{
						"class": "paste-input",
						"rows":  "30",
						"name":  "code",
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
								"type":  "checkbox",
								"class": "paste-private",
								"name":  "private",
							}),
							"Private")).Attrs(As{
						"class": "paste-notes",
					}),
					Input().Attrs(As{
						"class":     "paste-submit",
						"type":      "submit",
						"value":     "Go Paste!",
						"accesskey": "s",
					}))).Attrs(As{
				"action": "/add",
				"method": "POST",
			})).Attrs(As{
			"id": "home",
		}))))

var viewPage = template.Must(templ("view").Parse(
	page(
		"Paste #{{.id|html}} | Go Paste!",
		Div(
			Ul(
				Li(
					A("raw").Attrs(As{
						"href":  "/raw/{{.id|url_html}}",
						"class": "raw",
					})),
				Li(
					A("home").Attrs(As{
						"href": "/",
					}))).Attrs(As{
				"class": "view-links",
			}),
			"{{.id|code}}").Attrs(As{
			"id": "view",
		}))))

var allPage = template.Must(templ("all").Parse(
	page(
		"All | Go Paste!",
		Div(
      "{{range $paste := .pastes}}",
			"{{.}}",
      "{{end}}",

			Div(
				"{{if .prev}}",
				A("&larr; Prev").Attrs(As{
					"class": "page prev-page",
					"href":  "{prev}",
				}),
				"{{end}}",

				"{{if .next}}",
				A("Next &rarr;").Attrs(As{
					"class": "page next-page",
					"href":  "{next}",
				}),
				"{{end}}").Attrs(As{
				"class": "pagination",
			})).Attrs(As{
			"id": "all",
		}))))

func urlHtmlFormatter(v string) string {
	return template.HTMLEscapeString(url.QueryEscape(v))
}

func codePrinter(v string) string {
	fmt.Printf("Printin' the ol' code...\n")
	var buffer bytes.Buffer

	code, e := prettyPaste(v, 0)

  if e != nil {
    fmt.Printf("Error pretty-printing: %s\n", e)
  }

  fmt.Printf("Writing %s\n", code)

  if len(code) > 1 {
    buffer.WriteString(`<div class="multi-paste">`)
  }

  for _, part := range code {
    buffer.WriteString(part)
  }

  if len(code) > 1 {
    buffer.WriteString(`</div>`)
  }

  return buffer.String()
}
