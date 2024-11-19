package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

type HTTPRequest struct {
	Method      string
	Target      string
	HTTPVersion string
	Headers     map[string]string
}

func main() {
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
	log.Println("Accepted connection from:", c.RemoteAddr())

	data := make([]byte, 1024)
	_, err := c.Read(data)
	if err != nil {
		log.Printf("Error reading from %s connection: %s\n", c.RemoteAddr(), err.Error())
		return
	}

	req, err := parseRequest(data)
	if err != nil {
		log.Printf("Error parsing request from %connection: %s\n", c.RemoteAddr(), err.Error())
	}

	res := craftResponse(req.Target)

	_, err = c.Write(res)
	if err != nil {
		log.Printf("Error writing to %s connection: %s\n", c.RemoteAddr(), err.Error())
		return
	}
}

func parseRequest(d []byte) (*HTTPRequest, error) {
	request := new(HTTPRequest)

	reqLine, _, hasCRLF := bytes.Cut(d, []byte("\r\n"))
	if !hasCRLF {
		return request, fmt.Errorf("Received malformed request.")
	}

	splitReqLine := bytes.Split(reqLine, []byte(" "))
	if len(splitReqLine) != 3 {
		return request, fmt.Errorf("Received malformed request line.")
	}

	request.Method = string(splitReqLine[0])
	request.Target = string(splitReqLine[1])
	request.HTTPVersion = string(splitReqLine[2])

	return request, nil
}

func craftResponse(t string) []byte {
	if t == "/" {
		return []byte("HTTP/1.1 200 OK\r\n\r\n")
	}

	s, found := strings.CutPrefix(t, "/echo/")
	if found {
		res := "HTTP/1.1 200 OK\r\n"
		res += "Content-Type: text/plain\r\n"
		res += "Content-Length: "
		res += strconv.Itoa(len(s))
		res += "\r\n\r\n"
		res += s
		return []byte(res)
	}

	return []byte("HTTP/1.1 404 Not Found\r\n\r\n")
}
