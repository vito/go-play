package html

type As map[string]string

type Element struct {
	name       string
	contents   []string
	attributes As
}

func (self *Element) Out() string {
	attrs := ""

	for key, val := range self.attributes {
		attrs += " " + key + "=\"" + val + "\""
	}

	s := "<" + self.name + attrs

	if len(self.contents) == 0 {
		return s + " />"
	}

	s += ">"

	for _, content := range self.contents {
		s += content
	}

	s += "</" + self.name + ">"

	return s
}

func (self *Element) Append(content ...interface{}) *Element {
  for _, e := range content {
		switch e.(type) {
		case string:
			self.contents = append(self.contents, e.(string))
		default:
		  self.contents = append(self.contents, e.(*Element).Out())
		}
	}

	return self
}

func (self *Element) Attrs(attrs As) *Element {
	self.attributes = attrs
	return self
}

func New(name string, content ...interface{}) *Element {
	ele := new(Element)
	ele.name = name
	ele.contents = []string{}
	ele.attributes = nil

  ele.Append(content...)

	return ele
}

// HTML tags
func Html(content ...interface{}) *Element { return New("html", content...) }

func Head(content ...interface{}) *Element { return New("head", content...) }

func Title(content ...interface{}) *Element { return New("title", content...) }

func Base(content ...interface{}) *Element { return New("base", content...) }

func Link(content ...interface{}) *Element { return New("link", content...) }

func Meta(content ...interface{}) *Element { return New("meta", content...) }

func Style(content ...interface{}) *Element { return New("style", content...) }

func Script(content ...interface{}) *Element { return New("script", content...) }

func Noscript(content ...interface{}) *Element { return New("noscript", content...) }

func Body(content ...interface{}) *Element { return New("body", content...) }

func Section(content ...interface{}) *Element { return New("section", content...) }

func Nav(content ...interface{}) *Element { return New("nav", content...) }

func Article(content ...interface{}) *Element { return New("article", content...) }

func Aside(content ...interface{}) *Element { return New("aside", content...) }

func H1(content ...interface{}) *Element { return New("h1", content...) }

func H2(content ...interface{}) *Element { return New("h2", content...) }

func H3(content ...interface{}) *Element { return New("h3", content...) }

func H4(content ...interface{}) *Element { return New("h4", content...) }

func H5(content ...interface{}) *Element { return New("h5", content...) }

func H6(content ...interface{}) *Element { return New("h6", content...) }

func Hgroup(content ...interface{}) *Element { return New("hgroup", content...) }

func Header(content ...interface{}) *Element { return New("header", content...) }

func Footer(content ...interface{}) *Element { return New("footer", content...) }

func Address(content ...interface{}) *Element { return New("address", content...) }

func P(content ...interface{}) *Element { return New("p", content...) }

func Hr() *Element { return New("hr") }

func Br() *Element { return New("br") }

func Pre(content ...interface{}) *Element { return New("pre", content...) }

func Blockquote(content ...interface{}) *Element { return New("blockquote", content...) }

func Ol(content ...interface{}) *Element { return New("ol", content...) }

func Ul(content ...interface{}) *Element { return New("ul", content...) }

func Li(content ...interface{}) *Element { return New("li", content...) }

func Dl(content ...interface{}) *Element { return New("dl", content...) }

func Dt(content ...interface{}) *Element { return New("dt", content...) }

func Dd(content ...interface{}) *Element { return New("dd", content...) }

func A(content ...interface{}) *Element { return New("a", content...) }

func Q(content ...interface{}) *Element { return New("q", content...) }

func Cite(content ...interface{}) *Element { return New("cite", content...) }

func Em(content ...interface{}) *Element { return New("em", content...) }

func Strong(content ...interface{}) *Element { return New("strong", content...) }

func Small(content ...interface{}) *Element { return New("small", content...) }

func Mark(content ...interface{}) *Element { return New("mark", content...) }

func Dfn(content ...interface{}) *Element { return New("dfn", content...) }

func Abbr(content ...interface{}) *Element { return New("abbr", content...) }

func Time(content ...interface{}) *Element { return New("time", content...) }

func Progress(content ...interface{}) *Element { return New("progress", content...) }

func Meter(content ...interface{}) *Element { return New("meter", content...) }

func Code(content ...interface{}) *Element { return New("code", content...) }

func Var(content ...interface{}) *Element { return New("var", content...) }

func Samp(content ...interface{}) *Element { return New("samp", content...) }

func Kbd(content ...interface{}) *Element { return New("kbd", content...) }

func Sub(content ...interface{}) *Element { return New("sub", content...) }

func Sup(content ...interface{}) *Element { return New("sup", content...) }

func Span(content ...interface{}) *Element { return New("span", content...) }

func I(content ...interface{}) *Element { return New("i", content...) }

func B(content ...interface{}) *Element { return New("b", content...) }

func Bdo(content ...interface{}) *Element { return New("bdo", content...) }

func Ruby(content ...interface{}) *Element { return New("ruby", content...) }

func Rt(content ...interface{}) *Element { return New("rt", content...) }

func Rp(content ...interface{}) *Element { return New("rp", content...) }

func Ins(content ...interface{}) *Element { return New("ins", content...) }

func Del(content ...interface{}) *Element { return New("del", content...) }

func Figure(content ...interface{}) *Element { return New("figure", content...) }

func Img(content ...interface{}) *Element { return New("img", content...) }

func Iframe(content ...interface{}) *Element { return New("iframe", content...) }

func Embed(content ...interface{}) *Element { return New("embed", content...) }

func Object(content ...interface{}) *Element { return New("object", content...) }

func Param(content ...interface{}) *Element { return New("param", content...) }

func Video(content ...interface{}) *Element { return New("video", content...) }

func Audio(content ...interface{}) *Element { return New("audio", content...) }

func Source(content ...interface{}) *Element { return New("source", content...) }

func Canvas(content ...interface{}) *Element { return New("canvas", content...) }

func Map(content ...interface{}) *Element { return New("map", content...) }

func Area(content ...interface{}) *Element { return New("area", content...) }

func Table(content ...interface{}) *Element { return New("table", content...) }

func Caption(content ...interface{}) *Element { return New("caption", content...) }

func Colgroup(content ...interface{}) *Element { return New("colgroup", content...) }

func Col(content ...interface{}) *Element { return New("col", content...) }

func Tbody(content ...interface{}) *Element { return New("tbody", content...) }

func Thead(content ...interface{}) *Element { return New("thead", content...) }

func Tfoot(content ...interface{}) *Element { return New("tfoot", content...) }

func Tr(content ...interface{}) *Element { return New("tr", content...) }

func Td(content ...interface{}) *Element { return New("td", content...) }

func Th(content ...interface{}) *Element { return New("th", content...) }

func Form(content ...interface{}) *Element { return New("form", content...) }

func Fieldset(content ...interface{}) *Element { return New("fieldset", content...) }

func Label(content ...interface{}) *Element { return New("label", content...) }

func Input(content ...interface{}) *Element { return New("input", content...) }

func Button(content ...interface{}) *Element { return New("button", content...) }

func Select(content ...interface{}) *Element { return New("select", content...) }

func Datalist(content ...interface{}) *Element { return New("datalist", content...) }

func Optgroup(content ...interface{}) *Element { return New("optgroup", content...) }

func Option(content ...interface{}) *Element { return New("option", content...) }

func Textarea(content ...interface{}) *Element { return New("textarea", content...) }

func Keygen(content ...interface{}) *Element { return New("keygen", content...) }

func Output(content ...interface{}) *Element { return New("output", content...) }

func Details(content ...interface{}) *Element { return New("details", content...) }

func Command(content ...interface{}) *Element { return New("command", content...) }

func Menu(content ...interface{}) *Element { return New("menu", content...) }

func Legend(content ...interface{}) *Element { return New("legend", content...) }

func Div(content ...interface{}) *Element { return New("div", content...) }
