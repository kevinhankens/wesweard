package main

import (
	"flag"
	"fmt"
	"github.com/grafov/bcast"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var group = bcast.NewGroup() // create broadcast group
//var members = make([]*bcast.Member, 5)

func main() {
	var portBcast string
	flag.StringVar(&portBcast, "port-bcast", "5555", "The port to use for broadcasting.")
	var portRecv string
	flag.StringVar(&portRecv, "port-recv", "4444", "The port to use for receiving.")
	flag.Parse()
	fmt.Printf("Listening for receiving clients on port %s\n", portRecv)
	fmt.Printf("Listening for broadcasting clients on port %s\n", portBcast)

	go group.Broadcasting(0) // accepts messages and broadcast it to all members
	//go randomBroadcast()
	go listenerBroadcasting(portBcast)
	go listenerReceiving(portRecv)

	//group.Send("Group message 1")
	//members[0].Send("Group message 2")
	//fmt.Println(members[0].Recv())

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs,
		syscall.SIGABRT,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	<-done
	fmt.Println("exiting")
}

func listenerBroadcasting(portRecv string) {
	listener, err := net.Listen("tcp", ":"+portRecv)
	if err != nil {
		panic(fmt.Sprintf("Could not establish a TCP socket. %s", err))
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not accept connection", err)
			return
		}
		go handleBroadcasting(conn)
	}

}

func handleBroadcasting(conn net.Conn) {
	//member := group.Join() // joined member1 from one routine

	readLength := 4096
	buf := make([]byte, readLength)
	length, _ := io.ReadFull(conn, buf)

	group.Send(string(buf[:length]))
	conn.Close()
}

func listenerReceiving(portBcast string) {
	listener, err := net.Listen("tcp", ":"+portBcast)
	if err != nil {
		panic(fmt.Sprintf("Could not establish a TCP socket. %s", err))
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Could not accept connection", err)
			return
		}
		go handleReceiving(conn)
	}

}

func randomBroadcast() {
	for {
		time.Sleep(5000 * time.Millisecond)
		group.Send("Group message!")
	}
}

func handleReceiving(conn net.Conn) {
	member := group.Join() // joined member1 from one routine

	defer member.Close()
	defer conn.Close()
	defer fmt.Println("Connection closed")

	message := ""
	fmt.Println("Connection opened")
	for {
		select {
		case tmp := <-member.In:
			fmt.Println("Write.")
			message = tmp.(string)
			conn.Write([]byte(message))
			conn.Write([]byte{0})
		case <-time.After(60 * time.Second):
			fmt.Println("Timeout.")
			length, err := conn.Write([]byte{0})
			fmt.Printf("Read %v %v\n", length, err)
			if err != nil {
				return
			}
		}
	}
}
