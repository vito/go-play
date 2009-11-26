package main

import (
	"fmt";
	"time";
	"./fsm";
)


type TestFSM struct {
	real		int;
	imaginary	int;
}

const (
	OK	= iota;
	ERROR;
	RESET;
	GET_STATE;
	INCREASE;
	DECREASE;
	ADD;
	SUBTRACT;
)


func (self *TestFSM) Init(s chan<- fsm.StateHandler, inst *fsm.Instance, msg *fsm.Message) {
	self.real = msg.Data[0].(int);
	self.imaginary = msg.Data[1].(int);

	s <- Real;
}

func (self *TestFSM) HandleEvent(inst *fsm.Instance, msg *fsm.Message) {
	switch msg.What {
	case RESET:
		self.real = 0;
		self.imaginary = 0;
	default:
		fmt.Printf("Got unknown event: %#v\n", msg)
	}
}

func (self *TestFSM) HandleSyncEvent(r chan<- *fsm.Message, inst *fsm.Instance, msg *fsm.Message) {
	switch msg.What {
	case GET_STATE:
		r <- fsm.M(OK, self.real, self.imaginary)
	default:
		fmt.Printf("Got unknown sync event: %#v\n", msg);
		r <- fsm.M(ERROR, "Unknown sync event.");
	}
}

func Real(srv fsm.Server, inst *fsm.Instance, msg *fsm.Message) {
	self := srv.(*TestFSM);
	switch msg.What {
	case INCREASE:
		fmt.Printf("REAL: Increasing.\n");
		self.real++;
	case DECREASE:
		fmt.Printf("REAL: Decreasing.\n");
		self.real--;
	case ADD:
		fmt.Printf("REAL: Adding %#v.\n", msg.Data[0]);
		self.real += msg.Data[0].(int);
	case SUBTRACT:
		fmt.Printf("REAL: Subtracting %#v.\n", msg.Data[0]);
		self.real -= msg.Data[0].(int);
	}
}

func Imaginary(srv fsm.Server, inst *fsm.Instance, msg *fsm.Message) {
	self := srv.(*TestFSM);
	switch msg.What {
	case INCREASE:
		fmt.Printf("IMAGINARY: Increasing.\n");
		self.imaginary++;
	case DECREASE:
		fmt.Printf("IMAGINARY: Decreasing.\n");
		self.imaginary--;
	case ADD:
		fmt.Printf("IMAGINARY: Adding %#v.\n", msg.Data[0]);
		self.imaginary += msg.Data[0].(int);
	case SUBTRACT:
		fmt.Printf("IMAGINARY: Subtracting %#v.\n", msg.Data[0]);
		self.imaginary -= msg.Data[0].(int);
	}
}

func main() {
	inst := fsm.Start(new(TestFSM), fsm.M(OK, 0, 0));

	fmt.Printf("Started.\n");

	inst.SendEvent(fsm.M(INCREASE));
	inst.SendEvent(fsm.M(ADD, 100));
	inst.SendEvent(fsm.M(SUBTRACT, -320));
	inst.SendEvent(fsm.M(DECREASE));

	inst.SetState(Imaginary);

	inst.SendEvent(fsm.M(INCREASE));
	inst.SendEvent(fsm.M(ADD, 10));
	inst.SendEvent(fsm.M(SUBTRACT, -32));
	inst.SendEvent(fsm.M(DECREASE));

	time.Sleep(100);	// Slight delay to let all those finish, just for the sake of this demo.

	fmt.Printf("Current state: %#v\n", <-inst.SendSyncEvent(fsm.M(GET_STATE)));

	inst.SendAllEvent(fsm.M(RESET));

	fmt.Printf("Current state: %#v\n", <-inst.SendSyncEvent(fsm.M(GET_STATE)));
}
