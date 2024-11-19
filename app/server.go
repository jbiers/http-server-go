package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Println("Error failed to bind to port 4221")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()
	log.Println("Accepted connection from: ", c.RemoteAddr())

	// request will be interpreted here

	res := []byte("HTTP/1.1 200 OK\r\n\r\n")

	_, err := c.Write(res)
	if err != nil {
		log.Printf("Error writing to %s connection: %s\n", c.RemoteAddr(), "err")
		return
	}
}
