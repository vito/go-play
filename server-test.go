package main

import (
    "fmt";
    "./server";
)


type TestServer struct {
    count int;
}

const (
    OK = iota;
    GET;
    INCREASE;
    DECREASE;
    ADD;
    SUBTRACT;
)


func (self *TestServer) Init(inst *server.Instance, arg server.Value) {
    self.count = arg.(int);
    inst.Respond(server.M(OK, "Started."));
}

func (self *TestServer) HandleCall(response chan<- server.Value, inst *server.Instance, msg *server.Message) {
    switch msg.What {
        case GET:
            response <- self.count;
    }
}

func (self *TestServer) HandleCast(inst *server.Instance, msg *server.Message) {
    switch msg.What {
        case INCREASE:
            self.count++;
        case DECREASE:
            self.count--;
        case ADD:
            self.count += msg.Data[0].(int);
        case SUBTRACT:
            self.count -= msg.Data[0].(int);
    }
}


func main() {
    inst, result := server.Start(new(TestServer), 0);

    fmt.Printf("Started; result: %#v\n", result);

    fmt.Printf("Call: %v\n", <-inst.Call(server.M(GET)));

    inst.Cast(server.M(INCREASE));
    fmt.Printf("Call: %v\n", <-inst.Call(server.M(GET)));

    inst.Cast(server.M(ADD, 100));
    fmt.Printf("Call: %v\n", <-inst.Call(server.M(GET)));

    inst.Cast(server.M(DECREASE));
    fmt.Printf("Call: %v\n", <-inst.Call(server.M(GET)));

    inst.Cast(server.M(SUBTRACT, 100));
    fmt.Printf("Call: %v\n", <-inst.Call(server.M(GET)));

    fmt.Printf("All done!\n");
}
