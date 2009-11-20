package server

import "reflect"


type Value interface{}

type Message struct {
	What	int;
	Data	[]Value;
}

type Server interface {
    Init(chan<- Value, Value);
    HandleCall(chan<- Value, *Message);
    HandleCast(*Message);
}


func Start(srv Server, arg Value) <-chan Value {
	r := make(chan Value);
    go srv.Init(r, arg);

    return r;
}

func Call(srv Server, msg *Message) <-chan Value {
	c := make(chan Value);
    go srv.HandleCall(c, msg);
    return c;
}

func Cast(srv Server, msg *Message) {
    go srv.HandleCast(msg);
}


func M(what int, data ...) *Message {
	msg := new(Message);
	msg.What = what;

	v := reflect.NewValue(data).(*reflect.StructValue);

	msg.Data = make([]Value, v.NumField());
	for i := 0; i < v.NumField(); i++ {
		msg.Data[i] = v.Field(i).Interface()
	}

	return msg;
}

