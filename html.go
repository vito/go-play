package html

import ("reflect");


type A map[string]string;

type Element struct {
    name string;
    contents []string;
    attributes map[string]string;
}


func (self *Element) Out() string {
    attrs := "";

    for key, val := range self.attributes {
        attrs += " " + key + "=\"" + val + "\"";
    }

    s := "<" + self.name + attrs;

    if len(self.contents) == 0 {
        return s + " />";
    }

    s += ">";

    for idx, val := range self.contents {
        if idx > 0 {
            s += " ";
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
            ele.contents[i] = v.Field(i).Interface().(string);
        default:
            ele.contents[i] = v.Field(i).Interface().(*Element).Out();
        }
    }

    ele.attributes = nil;

    return ele;
}


// HTML tags
func Html(content ...) *Element {
    return New("html", content);
}

func Head(content ...) *Element {
	return New("head", content);
}

func Title(content ...) *Element {
	return New("title", content);
}

func Body(content ...) *Element {
	return New("body", content);
}

func Img(content ...) *Element {
	return New("img", content);
}

func Form(content ...) *Element {
	return New("form", content);
}

func Input(content ...) *Element {
	return New("input", content);
}

func Div(content ...) *Element {
    return New("div", content);
}

func Br() *Element {
    return New("br");
}


