package html

import (
	"reflect";
)


type A map[string]string

type Element struct {
	name		string;
	contents	[]string;
	attributes	map[string]string;
}


func (self *Element) Out() string {
	attrs := "";

	for key, val := range self.attributes {
		attrs += " " + key + "=\"" + val + "\""
	}

	s := "<" + self.name + attrs;

	if self.contents == nil {
		return s + " />"
	}

	s += ">";

	for idx, val := range self.contents {
		if idx > 0 {
			s += " "
		}

		s += val;
	}

	s += "</" + self.name + ">";

	return s;
}

func (self *Element) Attrs(attrs A) *Element {
	self.attributes = attrs;
	return self;
}


func New(name string, content ...) *Element {
	ele := new(Element);
	ele.name = name;

	v := reflect.NewValue(content).(*reflect.StructValue);

	ele.contents = make([]string, v.NumField());
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Interface().(type) {
		case string:
			ele.contents[i] = v.Field(i).Interface().(string)
		default:
			ele.contents[i] = v.Field(i).Interface().(*Element).Out()
		}
	}

	ele.attributes = nil;

	return ele;
}


// HTML tags
func Html(content ...) *Element	{ return New("html", content) }

func Head(content ...) *Element	{ return New("head", content) }

func Title(content ...) *Element	{ return New("title", content) }

func Link(content ...) *Element	{ return New("link", content) }

func Body(content ...) *Element	{ return New("body", content) }

func Img(content ...) *Element	{ return New("img", content) }

func Form(content ...) *Element	{ return New("form", content) }

func Input(content ...) *Element	{ return New("input", content) }

func Textarea(content ...) *Element	{ return New("textarea", content) }

func Div(content ...) *Element	{ return New("div", content) }

func Pre(content ...) *Element	{ return New("pre", content) }

func H1(content ...) *Element	{ return New("h1", content) }

func H2(content ...) *Element	{ return New("h2", content) }

func H3(content ...) *Element	{ return New("h3", content) }

func H4(content ...) *Element	{ return New("h4", content) }

func H5(content ...) *Element	{ return New("h5", content) }

func H6(content ...) *Element	{ return New("h6", content) }

func Br() *Element	{ return New("br") }
