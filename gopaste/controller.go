package controller

import (
	"http";
	"fmt";
	"reflect";
	"regexp";
	"strconv";
	"container/vector";
)


type callback struct {
	match		string;
	funcValue	*reflect.FuncValue;
	funcType	*reflect.FuncType;
}

type Controller struct {
	callbacks *vector.Vector;
}

func New() *Controller {
	cont := new(Controller);
	cont.callbacks = vector.New(0);
	return cont;
}

func (self *Controller) AddHandler(regexp string, fun interface{}) {
	self.callbacks.Push(&callback{
		regexp,
		reflect.NewValue(fun).(*reflect.FuncValue),
		reflect.Typeof(fun).(*reflect.FuncType),
	})
}

func (self *Controller) Handler() (func(*http.Conn, *http.Request)) {
	return func(c *http.Conn, req *http.Request) { self.Handle(c, req) }
}

func (self *Controller) Handle(c *http.Conn, req *http.Request) {
	for i := 0; i < self.callbacks.Len(); i++ {
		callback := self.callbacks.At(i).(*callback);

		match := `^` + callback.match;

		regexp, ok := regexp.Compile(match);
		if ok != nil {
			fmt.Printf("Match could not compile: %#v\n", match);
			continue;
		}

		values := make([]string, 0);
		if callback.funcType.NumIn() > 0 {
			values = regexp.MatchStrings(req.URL.Path)
		}

		if len(values) == 0 || (len(values)-1+2) != callback.funcType.NumIn() {
			continue
		}

		args := make([]reflect.Value, len(values)-1+2);

		args[0] = reflect.NewValue(c);
		args[1] = reflect.NewValue(req);
		for i := 0; i < len(values)-1; i++ {
			switch callback.funcType.In(i + 2).(type) {
			case *reflect.IntType:
				asInt, ok := strconv.Atoi(values[i+1]);
				if ok == nil {
					args[i+2] = reflect.NewValue(asInt)
				} else {
					args[i+2] = reflect.NewValue(-1)
				}
			case *reflect.BoolType:
				args[i+2] = reflect.NewValue(values[i+1] == "1" || values[i+1] == "true")
			default:
				args[i+2] = reflect.NewValue(values[i+1])
			}
		}

		callback.funcValue.Call(args);

		return;
	}
}
