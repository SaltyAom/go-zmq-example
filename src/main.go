package main

import (
	"fmt"
	"math"

	zmq "github.com/pebbe/zmq4"

	"time"
)

func sub() {
	socket, _ := zmq.NewSocket(zmq.REP)

	defer socket.Close()
	socket.Bind("tcp://0.0.0.0:5556")

	for {
		socket.Recv(0)

		socket.Send("World", 0)
	}
}

func sendMessage(message string) {
	socket, _ := zmq.NewSocket(zmq.REQ)
	socket.Connect("tcp://0.0.0.0:5556")

	socket.Send(message, 0)

	socket.Recv(0)

	socket.Close()
}

func pub(queues <-chan string) {
	for message := range queues {
		sendMessage(message)
	}
}

const limit = 10_000_000

func main() {
	go sub()

	ch := make(chan string, int(math.Pow(2, 24)))
	go pub(ch)

	// ? Make sure everything is initialize
	ch <- "Hello"

	t := time.Now()

	for i := 0; i <= limit; i++ {
		ch <- "Hello"
	}

	m := time.Since(t).Microseconds()

	fmt.Printf("Done %v process in %v ms or %v us\n", limit, m/1000, m)
	if m != 0 {
		fmt.Printf("Average %v req/s\n", int(math.Round(100000/float64(m)*float64(limit))))
	}
}
