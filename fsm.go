package fsm

import "./server"


type Value interface{}

type StateHandler func(Server, *Instance, *Message)

type Instance struct {
	server	*Server;
	state	StateHandler;
}

type Message struct {
	What	int;
	Data	[]Value;
}

type Server interface {
	Init(chan<- StateHandler, *Instance, *Message);
	HandleEvent(*Instance, *Message);
	HandleSyncEvent(chan<- *Message, *Instance, *Message);
}


func Start(srv Server, msg *Message) *Instance {
	inst := new(Instance);
	inst.server = &srv;

	s := make(chan StateHandler);
	go srv.Init(s, inst, msg);
	inst.state = <-s;

	return inst;
}

func M(what int, data ...) *Message {
	msg := Message(*server.M(what, data));
	return &msg;
}


func (inst *Instance) SetState(state func(Server, *Instance, *Message)) {
	inst.state = state
}

func (inst *Instance) SendEvent(msg *Message)	{ go inst.state(*inst.server, inst, msg) }

func (inst *Instance) SendSyncEvent(msg *Message) chan *Message {
	r := make(chan *Message);
	go inst.server.HandleSyncEvent(r, inst, msg);
	return r;
}

func (inst *Instance) SendAllEvent(msg *Message) {
	go inst.server.HandleEvent(inst, msg)
}
