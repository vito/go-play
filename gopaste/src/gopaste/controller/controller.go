package controller

import (
	"container/list"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
)

type callback struct {
	match     string
	funcValue reflect.Value
	funcType  reflect.Type
}

type Controller struct {
	callbacks *list.List
}

func New() *Controller {
	cont := new(Controller)
	cont.callbacks = new(list.List)
	return cont
}

func (self *Controller) AddHandler(regexp string, fun interface{}) {
	fmt.Printf("Adding handler for %s\n", regexp)

	self.callbacks.PushBack(&callback{
		regexp,
		reflect.ValueOf(fun),
		reflect.TypeOf(fun),
	})
}

func (self *Controller) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(c http.ResponseWriter, req *http.Request) { self.Handle(c, req) })
}

func (self *Controller) Handle(c http.ResponseWriter, req *http.Request) {
  fmt.Printf("Handling %v\n", req)

	for d := self.callbacks.Front(); d != nil; d = d.Next() {
		callback := d.Value.(*callback)

		match := `^` + callback.match

		regexp, ok := regexp.Compile(match)
		if ok != nil {
			fmt.Printf("Match could not compile: %#v\n", match)
			continue
		}

		values := make([]string, 0)
		if callback.funcType.NumIn() > 0 {
			values = regexp.FindStringSubmatch(req.URL.Path)
		}

    fmt.Printf("Match: %s, %s, %v\n", match, req.URL.Path, values)

    fmt.Printf(
      "Value count: %v, required args: %v\n",
      len(values) - 1 + 2, callback.funcType.NumIn())

		if len(values) == 0 || (len(values)-1+2) != callback.funcType.NumIn() {
		  fmt.Printf("Skipping!\n")
			continue
		}

		args := make([]reflect.Value, len(values)-1+2)

		args[0] = reflect.ValueOf(c)
		args[1] = reflect.ValueOf(req)
		for i := 0; i < len(values)-1; i++ {
			switch callback.funcType.In(i + 2).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				asInt, ok := strconv.Atoi(values[i+1])
				if ok == nil {
					args[i+2] = reflect.ValueOf(asInt)
				} else {
					args[i+2] = reflect.ValueOf(-1)
				}
			case reflect.Bool:
				args[i+2] = reflect.ValueOf(values[i+1] == "1" || values[i+1] == "true")
			default:
				args[i+2] = reflect.ValueOf(values[i+1])
			}
		}

		fmt.Printf("Invoking callback with %v\n", args)

		callback.funcValue.Call(args)

		return
	}
}
