package pretty

import (
	"fmt";
    "go/ast";
    "go/token";
    "go/parser";
    "go/printer";
    "io";
    "strings";
)


type HTMLStyler struct {}


func (self *HTMLStyler) LineTag(line int) ([]byte, printer.HTMLTag) {
    return strings.Bytes(fmt.Sprintf("%04d", line)), printer.HTMLTag{
		Start: "<a rel='#L" + fmt.Sprint(line) + "' href='#L" + fmt.Sprint(line) + "' id='L" + fmt.Sprint(line) + "' class='line-number'>",
		End: "</a>"
	};
}

func (self *HTMLStyler) Comment(comment *ast.Comment, line []byte) ([]byte, printer.HTMLTag) {
    return comment.Text, printer.HTMLTag{
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

func Print(w io.Writer, filename string, source interface {}) {
    fileAst, _ := parser.ParseFile(filename, source, 4);

	// Assume they forgot the package declaration
    if len(fileAst.Decls) == 0 && source != nil {
		fileAst, _ = parser.ParseFile(filename, "package main\n\n" + source.(string), 4);
	}

	(&printer.Config{
		Mode: 5,
		Tabwidth: 4,
		Styler: new(HTMLStyler)
	}).Fprint(w, fileAst);
}
