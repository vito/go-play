package main

import (
	"container/vector";
	"fmt";
	"go/ast";
	"go/token";
	"go/parser";
	"go/printer";
	"io";
	"os";
	"path";
	"regexp";
	"strings";

	. "./html";
)


type HTMLStyler struct {
	comment		*ast.Comment;
	comment_text	[]string;
	comment_offset	int;
	prev		interface{};
}

type collector struct {
	contents *vector.StringVector;
}

func (self *collector) Write(p []byte) (n int, err os.Error) {
	self.contents.Push(string(p));
	return len(p), nil;
}


func (self *HTMLStyler) LineTag(line int) ([]byte, printer.HTMLTag) {
	return []byte{}, printer.HTMLTag{}
}

func (self *HTMLStyler) Comment(comment *ast.Comment, line []byte) ([]byte, printer.HTMLTag) {
	if self.comment == comment {
		self.comment_offset++;
		if self.comment_text[self.comment_offset] == "" {
			self.comment_offset++
		}
	} else {
		self.comment = comment;
		self.comment_text = strings.Split(string(comment.Text), "\r\n", 0);
		self.comment_offset = 0;
	}

	self.prev = comment;

	return strings.Bytes(self.comment_text[self.comment_offset]), printer.HTMLTag{
		Start: "<span class=\"go-comment\">",
		End: "</span>",
	};
}

func (self *HTMLStyler) BasicLit(x *ast.BasicLit) ([]byte, printer.HTMLTag) {
	kind := "other";
	switch x.Kind {
	case token.INT:
		kind = "int"
	case token.FLOAT:
		kind = "float"
	case token.CHAR:
		kind = "char"
	case token.STRING:
		kind = "string"
	}

	if x.Value[0] == '`' {
		kind = "string go-raw-string"
	}

	self.prev = x;

	return x.Value, printer.HTMLTag{
		Start: "<span class=\"go-basiclit go-" + kind + "\">",
		End: "</span>",
	};
}

func (self *HTMLStyler) Ident(id *ast.Ident) ([]byte, printer.HTMLTag) {
	classes := "go-local";
	if id.IsExported() {
		classes = "go-exported"
	}

	switch id.String() {
	case "bool", "uint8", "uint16", "uint32", "uint64", "int8", "int16", "int32", "int64", "float32", "float64", "byte", "uint", "int", "float", "uintptr", "string":
		classes += " go-prim-ident"
	default:
		if tok, ok := self.prev.(token.Token); ok && tok.String() == "func" || tok.String() == ")" {
			classes += " go-func-ident"
		}
	}

	self.prev = id;

	return strings.Bytes(id.String()), printer.HTMLTag{
		Start: "<span class=\"go-ident " + classes + "\">",
		End: "</span>",
	};
}

func (self *HTMLStyler) Token(tok token.Token) ([]byte, printer.HTMLTag) {
	extra := "";

	if tok.IsKeyword() {
		extra += " go-keyword"
	}

	if tok.IsLiteral() {
		extra += " go-literal"
	}

	if tok.IsOperator() {
		extra += " go-operator"
	}

	self.prev = tok;

	return strings.Bytes(tok.String()), printer.HTMLTag{
		Start: "<span class=\"go-token" + extra + "\">",
		End: "</span>",
	};
}

func Print(filename string, source interface{}) (pretty string, ok os.Error) {
	fileAst, ok := parser.ParseFile(filename, source, 4);

	// Make common corrections for snippet pastes
	if ok != nil && source != nil {
		src := source.(string);

		if m, _ := regexp.MatchString(`^package`, src); !m {
			src = "package main\n\n" + src
		}

		if fileAst, ok = parser.ParseFile(filename, src, 4); ok != nil {
			return
		}
	}

	coll := new(collector);
	coll.contents = new(vector.StringVector);

	(&printer.Config{
		Mode: 5,
		Tabwidth: 4,
		Styler: new(HTMLStyler),
	}).Fprint(coll, fileAst);

	pretty = strings.Join(coll.contents.Data(), "");

	return;
}

func prettyPaste(id string, limit int) (code []string, err os.Error) {
	source, err := io.ReadFile("pastes" + path.Clean("/"+id));
	if err != nil {
		return
	}

	multi := strings.Split(string(source), "\n---", 0);

	allCode := make([]string, len(multi));
	results := make(chan int);
	for i := 0; i < len(multi); i++ {
		go func(i int) {
			allCode[i], _ = prettySource(id, multi[i], limit);
			results <- i;
		}(i)
	}

	for i := 0; i < len(multi); i++ {
		<-results
	}

	code = allCode;

	return;
}

func prettySource(filename string, source string, limit int) (code string, err os.Error) {
	prettyCode, ok := Print(filename, source);
	if ok != nil {	// If it fails to parse, just serve it raw.
		prettyCode = source
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
                "\n",
				A("...").Attrs(As{
					"href": fmt.Sprintf("/view/%s", filename),
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
