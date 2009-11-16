package html

import (
    "container/vector";
	"reflect";
)


type As map[string]string

type Element struct {
	name		string;
	contents	*vector.StringVector;
	attributes	As;
}


func (self *Element) Out() string {
	attrs := "";

	for key, val := range self.attributes {
		attrs += " " + key + "=\"" + val + "\""
	}

	s := "<" + self.name + attrs;

	if self.contents.Len() == 0 {
		return s + " />"
	}

	s += ">";

	for _, content := range self.contents.Data() {
		s += content;
	}

	s += "</" + self.name + ">";

	return s;
}

func (self *Element) Append(content ...) *Element {
	v := reflect.NewValue(content).(*reflect.StructValue);

	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Interface().(type) {
		case string:
			self.contents.Push(v.Field(i).Interface().(string));
		default:
			self.contents.Push(v.Field(i).Interface().(*Element).Out());
		}
	}

	return self;
}

func (self *Element) Attrs(attrs As) *Element {
	self.attributes = attrs;
	return self;
}


func New(name string, content ...) *Element {
	ele := new(Element);
	ele.name = name;
    ele.contents = vector.NewStringVector(0);
    ele.attributes = nil;

    ele.Append(content);

    return ele;
}


// HTML tags
func Html(content ...) *Element	{ return New("html", content) }

func Head(content ...) *Element	{ return New("head", content) }

func Title(content ...) *Element	{ return New("title", content) }

func Base(content ...) *Element	{ return New("base", content) }

func Link(content ...) *Element	{ return New("link", content) }

func Meta(content ...) *Element	{ return New("meta", content) }

func Style(content ...) *Element	{ return New("style", content) }

func Script(content ...) *Element	{ return New("script", content) }

func Noscript(content ...) *Element	{ return New("noscript", content) }

func Body(content ...) *Element	{ return New("body", content) }

func Section(content ...) *Element	{ return New("section", content) }

func Nav(content ...) *Element	{ return New("nav", content) }

func Article(content ...) *Element	{ return New("article", content) }

func Aside(content ...) *Element	{ return New("aside", content) }

func H1(content ...) *Element	{ return New("h1", content) }

func H2(content ...) *Element	{ return New("h2", content) }

func H3(content ...) *Element	{ return New("h3", content) }

func H4(content ...) *Element	{ return New("h4", content) }

func H5(content ...) *Element	{ return New("h5", content) }

func H6(content ...) *Element	{ return New("h6", content) }

func Hgroup(content ...) *Element	{ return New("hgroup", content) }

func Header(content ...) *Element	{ return New("header", content) }

func Footer(content ...) *Element	{ return New("footer", content) }

func Address(content ...) *Element	{ return New("address", content) }

func P(content ...) *Element	{ return New("p", content) }

func Hr() *Element	{ return New("hr") }

func Br() *Element	{ return New("br") }

func Pre(content ...) *Element	{ return New("pre", content) }

func Blockquote(content ...) *Element	{ return New("blockquote", content) }

func Ol(content ...) *Element	{ return New("ol", content) }

func Ul(content ...) *Element	{ return New("ul", content) }

func Li(content ...) *Element	{ return New("li", content) }

func Dl(content ...) *Element	{ return New("dl", content) }

func Dt(content ...) *Element	{ return New("dt", content) }

func Dd(content ...) *Element	{ return New("dd", content) }

func A(content ...) *Element	{ return New("a", content) }

func Q(content ...) *Element	{ return New("q", content) }

func Cite(content ...) *Element	{ return New("cite", content) }

func Em(content ...) *Element	{ return New("em", content) }

func Strong(content ...) *Element	{ return New("strong", content) }

func Small(content ...) *Element	{ return New("small", content) }

func Mark(content ...) *Element	{ return New("mark", content) }

func Dfn(content ...) *Element	{ return New("dfn", content) }

func Abbr(content ...) *Element	{ return New("abbr", content) }

func Time(content ...) *Element	{ return New("time", content) }

func Progress(content ...) *Element	{ return New("progress", content) }

func Meter(content ...) *Element	{ return New("meter", content) }

func Code(content ...) *Element	{ return New("code", content) }

func Var(content ...) *Element	{ return New("var", content) }

func Samp(content ...) *Element	{ return New("samp", content) }

func Kbd(content ...) *Element	{ return New("kbd", content) }

func Sub(content ...) *Element	{ return New("sub", content) }

func Sup(content ...) *Element	{ return New("sup", content) }

func Span(content ...) *Element	{ return New("span", content) }

func I(content ...) *Element	{ return New("i", content) }

func B(content ...) *Element	{ return New("b", content) }

func Bdo(content ...) *Element	{ return New("bdo", content) }

func Ruby(content ...) *Element	{ return New("ruby", content) }

func Rt(content ...) *Element	{ return New("rt", content) }

func Rp(content ...) *Element	{ return New("rp", content) }

func Ins(content ...) *Element	{ return New("ins", content) }

func Del(content ...) *Element	{ return New("del", content) }

func Figure(content ...) *Element	{ return New("figure", content) }

func Img(content ...) *Element	{ return New("img", content) }

func Iframe(content ...) *Element	{ return New("iframe", content) }

func Embed(content ...) *Element	{ return New("embed", content) }

func Object(content ...) *Element	{ return New("object", content) }

func Param(content ...) *Element	{ return New("param", content) }

func Video(content ...) *Element	{ return New("video", content) }

func Audio(content ...) *Element	{ return New("audio", content) }

func Source(content ...) *Element	{ return New("source", content) }

func Canvas(content ...) *Element	{ return New("canvas", content) }

func Map(content ...) *Element	{ return New("map", content) }

func Area(content ...) *Element	{ return New("area", content) }

func Table(content ...) *Element	{ return New("table", content) }

func Caption(content ...) *Element	{ return New("caption", content) }

func Colgroup(content ...) *Element	{ return New("colgroup", content) }

func Col(content ...) *Element	{ return New("col", content) }

func Tbody(content ...) *Element	{ return New("tbody", content) }

func Thead(content ...) *Element	{ return New("thead", content) }

func Tfoot(content ...) *Element	{ return New("tfoot", content) }

func Tr(content ...) *Element	{ return New("tr", content) }

func Td(content ...) *Element	{ return New("td", content) }

func Th(content ...) *Element	{ return New("th", content) }

func Form(content ...) *Element	{ return New("form", content) }

func Fieldset(content ...) *Element	{ return New("fieldset", content) }

func Label(content ...) *Element	{ return New("label", content) }

func Input(content ...) *Element	{ return New("input", content) }

func Button(content ...) *Element	{ return New("button", content) }

func Select(content ...) *Element	{ return New("select", content) }

func Datalist(content ...) *Element	{ return New("datalist", content) }

func Optgroup(content ...) *Element	{ return New("optgroup", content) }

func Option(content ...) *Element	{ return New("option", content) }

func Textarea(content ...) *Element	{ return New("textarea", content) }

func Keygen(content ...) *Element	{ return New("keygen", content) }

func Output(content ...) *Element	{ return New("output", content) }

func Details(content ...) *Element	{ return New("details", content) }

func Command(content ...) *Element	{ return New("command", content) }

func Menu(content ...) *Element	{ return New("menu", content) }

func Legend(content ...) *Element	{ return New("legend", content) }

func Div(content ...) *Element	{ return New("div", content) }

