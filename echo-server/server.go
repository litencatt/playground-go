package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"strings"
	"syscall"
)

func main() {
	var address string
	var port int

	address = "0.0.0.0"
	port = 8080

	server := &Server{address: address, port: port}
	err := server.Start()
	if err != nil {
		log.Fatalf("TCP Echo Server startup failed, reason: %s", err)
	}
}

type Server struct {
	address string
	port    int
}

func (s *Server) Start() error {
	listenConfig := net.ListenConfig{Control: SetUpSocket}
	listner, err := listenConfig.Listen(context.Background(), "tcp", fmt.Sprintf("%s:%d", s.address, s.port))
	if err != nil {
		return err
	}

	for {
		conn, err := listner.Accept()
		if err != nil {
			continue
		}

		go func() {
			defer conn.Close()

			log.Printf("[%s] Accepted, client", conn.RemoteAddr())

			reader := bufio.NewReader(conn)
			message, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("[%s] Read error, reason: %s", conn.RemoteAddr(), err)
				return
			}

			log.Printf("[%s] Received message: %s", conn.RemoteAddr(), message)

			if strings.HasSuffix(message, "\r\n") {
				message = message[:len(message)-2]
			} else if strings.HasSuffix(message, "\n") {
				message = message[:len(message)-1]
			}

			writer := bufio.NewWriter(conn)
			writer.WriteString(message)
			writer.Flush()

			log.Printf("[%s] Replied, client", conn.RemoteAddr())
		}()
	}
}

func SetUpSocket(network string, address string, conn syscall.RawConn) error {
	return conn.Control(func(fd uintptr) {
		syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
	})
}
