package server

import "reflect";


type Value interface {};

type Instance struct {
    server *Server;
    channel chan Value;
}

type Message struct {
    What int;
    Data []Value;
}

type Server interface {
    Init(*Instance, Value);
    HandleCall(*Instance, *Message);
    HandleCast(*Instance, *Message);
}


func Start(srv Server, arg Value) (*Instance, Value) {
    inst := new(Instance);
    inst.server = &srv;
    inst.channel = make(chan Value);

    go srv.Init(inst, arg);

    return inst, <-inst.channel;
}

func M(what int, data ...) *Message {
    msg := new(Message);
    msg.What = what;

    v := reflect.NewValue(data).(*reflect.StructValue);

    msg.Data = make([]Value, v.NumField());
    for i := 0; i < v.NumField(); i++ {
        msg.Data[i] = v.Field(i).Interface();
    }

    return msg;
}


func (inst *Instance) Respond(val Value) {
    inst.channel <- val;
}

func (inst *Instance) Call(msg *Message) Value {
    go inst.server.HandleCall(inst, msg);
    return <-inst.channel;
}

func (inst *Instance) Cast(msg *Message) {
    go inst.server.HandleCast(inst, msg);
}
