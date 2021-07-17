package main

import (
	"log"
	"math"

	zmq "github.com/pebbe/zmq4"
)

type Queue struct {
	message  string
	response chan string
}

func sub() {
	socket, err := zmq.NewSocket(zmq.REP)

	if err != nil {
		log.Fatal(err)
		return
	}

	defer socket.Close()
	socket.Bind("tcp://0.0.0.0:5556")

	for {
		message, _ := socket.Recv(0)
		println("Req:", message)

		socket.Send("World", 0)
	}
}

func pub(queues <-chan Queue) {
	for queue := range queues {
		sendMessage(queue)
	}
}

func sendMessage(queue Queue) {
	socket, err := zmq.NewSocket(zmq.REQ)

	if err != nil {
		log.Fatal(err)
		return
	}

	socket.Connect("tcp://0.0.0.0:5556")
	socket.Send(queue.message, 0)

	message, _ := socket.Recv(0)

	queue.response <- message

	socket.Close()
}

func main() {
	go sub()

	ch := make(chan Queue, int(math.Pow(2, 24)))
	go pub(ch)

	// ? response waifu
	responseChan := make(chan string)

	ch <- Queue{
		message:  "Hello",
		response: responseChan,
	}

	println("Res:", <-responseChan)
}
