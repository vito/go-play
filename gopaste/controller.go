package controller

import (
	"http";
	"fmt";
	"reflect";
	"regexp";
	"strconv";
)


type Func struct {
	Value	*reflect.FuncValue;
	Type	*reflect.FuncType;
}

type Controller struct {
	callbacks map[string]Func;
}

func New(callbacks map[string]interface{}) *Controller {
	cont := new(Controller);
	cont.callbacks = make(map[string]Func);

	for regexp, fun := range callbacks {
		cont.SetHandler(regexp, fun)
	}

	return cont;
}

func (self *Controller) SetHandler(regexp string, callback interface{}) {
	self.callbacks[regexp] = Func{
		reflect.NewValue(callback).(*reflect.FuncValue),
		reflect.Typeof(callback).(*reflect.FuncType),
	}
}

func (self *Controller) Handler() (func (*http.Conn, *http.Request)) {
	return func(c *http.Conn, req *http.Request) {
		self.Handle(c, req);
	}
}

func (self *Controller) Handle(c *http.Conn, req *http.Request) {
	for match, callback := range self.callbacks {
		match = `^` + match;

		regexp, ok := regexp.Compile(match);
		if ok != nil {
			fmt.Printf("Match could not compile: %#v\n", match);
			continue;
		}

		values := []string{};
		if callback.Type.NumIn() > 0 {
			values = regexp.MatchStrings(req.URL.Path)
		}

		if len(values) == 0 || (len(values)-1+2) != callback.Type.NumIn() {
			continue
		}

		args := make([]reflect.Value, len(values)-1+2);

		args[0] = reflect.NewValue(c);
		args[1] = reflect.NewValue(req);
		for i := 0; i < len(values)-1; i++ {
			switch callback.Type.In(i + 2).String() {
			case "int":
				asInt, ok := strconv.Atoi(values[i+1]);
				if ok == nil {
					args[i+2] = reflect.NewValue(asInt)
				} else {
					args[i+2] = reflect.NewValue(-1)
				}
			default:
				args[i+2] = reflect.NewValue(values[i+1])
			}
		}

		callback.Value.Call(args);

		break;
	}
}
