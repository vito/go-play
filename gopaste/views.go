 package main

import (
	"http";
	"fmt";
	"io";
	"os";
	"path";
	"strings";
	"template";
	"./pretty";

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
	"code-truncated": truncatedCodePrinter,
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
				"href": "/raw?paste={@|url+html}",
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
	code, _ := codeLines(v.(string), 0);
	io.WriteString(w, code);
}

func truncatedCodePrinter(w io.Writer, v interface{}, _ string) {
	code, _ := codeLines(v.(string), 10);
	io.WriteString(w, code);
}

func codeLines(paste string, limit int) (code string, err os.Error) {
	source, ok := io.ReadFile("pastes" + path.Clean("/"+paste));

	if ok != nil {
		err = os.NewError(fmt.Sprintf("io.ReadFile: %s", ok));
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

	code = Table(
		Tbody(
			Tr(
				Td(linesPre).Attrs(As{"width": "1%", "valign": "top"}),
				Td(codePre).Attrs(As{"valign": "top"})))).Attrs(As{
		"class": "code",
		"cellspacing": "0",
		"cellpadding": "0",
	}).Out();

	return;
}
