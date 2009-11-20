package main

import (
	"fmt";
	"./server";
)


type TestServer struct {
	count int;
}

const (
	OK	= iota;
	GET;
	INCREASE;
	DECREASE;
	ADD;
	SUBTRACT;
)


func (self *TestServer) Init(r chan<- server.Value, arg server.Value) {
	self.count = arg.(int);
	r <- server.M(OK, "Started.");
}

func (self *TestServer) HandleCall(r chan<- server.Value, msg *server.Message) {
	switch msg.What {
	case GET:
		r <- self.count
	}
}

func (self *TestServer) HandleCast(msg *server.Message) {
	switch msg.What {
	case INCREASE:
		self.count++
	case DECREASE:
		self.count--
	case ADD:
		self.count += msg.Data[0].(int)
	case SUBTRACT:
		self.count -= msg.Data[0].(int)
	}
}


func main() {
	test := new(TestServer);
	result := server.Start(test, 0);

	fmt.Printf("Started; result: %#v\n", <-result);

	fmt.Printf("Call: %v\n", <-server.Call(test, server.M(GET)));

	server.Cast(test, server.M(INCREASE));
	fmt.Printf("Call: %v\n", <-server.Call(test, server.M(GET)));

	server.Cast(test, server.M(ADD, 100));
	fmt.Printf("Call: %v\n", <-server.Call(test, server.M(GET)));

	server.Cast(test, server.M(DECREASE));
	fmt.Printf("Call: %v\n", <-server.Call(test, server.M(GET)));

	server.Cast(test, server.M(SUBTRACT, 100));
	fmt.Printf("Call: %v\n", <-server.Call(test, server.M(GET)));

	fmt.Printf("All done!\n");
}
