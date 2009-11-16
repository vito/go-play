package pretty

import (
    "go/ast";
    "go/token";
    "go/parser";
    "go/printer";
	"os";
    "strings";
)


type HTMLStyler struct {
	comment *ast.Comment;
	comment_text []string;
	comment_offset int;
}

type collector struct {
	contents string;
}

func (self *collector) Write(p []byte) (n int, err os.Error) {
	self.contents += string(p);
	return len(p), nil;
}


func (self *HTMLStyler) LineTag(line int) ([]byte, printer.HTMLTag) {
	return []byte{}, printer.HTMLTag{};
}

func (self *HTMLStyler) Comment(comment *ast.Comment, line []byte) ([]byte, printer.HTMLTag) {
	if self.comment == comment {
		self.comment_offset += 1;
	} else {
		self.comment = comment;
		self.comment_text = strings.Split(string(comment.Text), "\n", 0);
		self.comment_offset = 0;
	}

    return strings.Bytes(self.comment_text[self.comment_offset]), printer.HTMLTag{
        Start: "<span class=\"go-comment\">",
        End: "</span>"
    };
}

func (self *HTMLStyler) BasicLit(x *ast.BasicLit) ([]byte, printer.HTMLTag) {
	kind := "other";
	switch x.Kind {
	case token.INT:
		kind = "int";
	case token.FLOAT:
		kind = "float";
	case token.CHAR:
		kind = "char";
	case token.STRING:
		kind = "string";
	}

    return x.Value, printer.HTMLTag{
        Start: "<span class=\"go-basiclit go-" + kind + "\">",
        End: "</span>"
    };
}

func (self *HTMLStyler) Ident(id *ast.Ident) ([]byte, printer.HTMLTag) {
	exported := "local";
	if id.IsExported() {
		exported = "exported";
	}

    return strings.Bytes(id.String()), printer.HTMLTag{
        Start: "<span class=\"go-ident go-" + exported + "\">",
        End: "</span>"
    };
}

func (self *HTMLStyler) Token(tok token.Token) ([]byte, printer.HTMLTag) {
	extra := "";

	if tok.IsKeyword() {
		extra += " go-keyword";
	}

	if tok.IsLiteral() {
		extra += " go-literal";
	}

	if tok.IsOperator() {
		extra += " go-operator";
	}

    return strings.Bytes(tok.String()), printer.HTMLTag{
        Start: "<span class=\"go-token" + extra + "\">",
        End: "</span>"
    };
}

func Print(filename string, source interface {}) string {
    fileAst, ok := parser.ParseFile(filename, source, 4);

	// Assume they forgot the package declaration
    if ok != nil && source != nil {
		fileAst, _ = parser.ParseFile(filename, "package main\n\n" + source.(string), 4);
	}

	coll := new(collector);
	(&printer.Config{
		Mode: 5,
		Tabwidth: 4,
		Styler: new(HTMLStyler)
	}).Fprint(coll, fileAst);

	return coll.contents;
}
