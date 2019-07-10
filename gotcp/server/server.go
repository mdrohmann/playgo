package main

import (
	// "io"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

type commandHandler func(string) string

type server struct {
	commands map[string]commandHandler
	count    int32

	countChan    chan int
	getCountChan chan int
	resetChan    chan int
}

func newServer() *server {
	s := server{commands: make(map[string]commandHandler)}
	s.count = 0

	s.countChan = make(chan int)
	s.resetChan = make(chan int)
	s.getCountChan = make(chan int)

	go func() {
		var channelCount int
		channelCount = 0
	outerLoop:
		for {
			select {
			case v := <-s.resetChan:
				if v == -1 {
					log.Print("break out of outer loop")
					break outerLoop
				}
				channelCount = v
			case addend := <-s.countChan:
				channelCount += addend
				s.getCountChan <- channelCount
			}
		}
		log.Println("Leaving channel count go-routine.")
	}()

	s.commands["echo"] = func(t string) string {
		return t + "\n"
	}
	s.commands["wait"] = func(t string) string {
		timeOut, err := strconv.Atoi(t)
		if err != nil {
			return fmt.Sprintf("Failed to parse argument %s, expected number\n", t)
		}
		time.Sleep(time.Duration(timeOut) * time.Second)
		return fmt.Sprintf("Slept for %d seconds", timeOut) + "\n"
	}
	s.commands["add"] = func(_ string) string {
		atomic.AddInt32(&(s.count), 1)
		return fmt.Sprintf("New count is %d\n", s.count)
	}
	s.commands["addChannel"] = func(t string) string {
		s.countChan <- 1
		return fmt.Sprintf("New channel count for %s is %d\n", t, <-s.getCountChan)
	}
	s.commands["getChannelCount"] = func(t string) string {
		s.countChan <- 0
		return fmt.Sprintf("Channel count for %s is %d\n", t, <-s.getCountChan)
	}
	s.commands["reset"] = func(_ string) string {
		atomic.StoreInt32(&(s.count), 0)
		s.resetChan <- 0
		return fmt.Sprintf("Reset counter to %d\n", s.count)
	}
	s.commands["STOP"] = func(t string) string {
		return "Closing connection\n"
	}
	s.commands["CLOSE"] = func(t string) string {
		return "Shutting down server\n"
	}
	return &s
}

func (s *server) close() {
	log.Print("closing server struct")
	s.resetChan <- -1

}

func (s *server) tcpServerListen() {

	log.Print("hello world")

	l, err := net.Listen("tcp", ":2000")
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	closeServer := make(chan bool, 1)

	for {
		// wait for a connection

		select {
		case cl := <-closeServer:
			log.Print("Closing server.")
			if cl {
				return
			}
		default:
			log.Print("Accepting new connections.")
		}

		conn, err := l.Accept()
		if err != nil {
			log.Print("Error on accepting the connection: ", err)
			continue
		}
		go func(c net.Conn) {
			log.Print("Received some data")
			defer c.Close()
			for {
				buf := make([]byte, 1024)

				n, err := c.Read(buf)
				if err != nil {
					log.Printf("Error reading bytes after receiving %d bytes: %s", n, err)
					break
				}
				trimmed := strings.TrimSpace(string(buf[0:n]))
				args := strings.SplitN(trimmed, " ", 2)
				var response = "unknown command\n"
				if len(args) > 0 {
					if handler, hasHandler := s.commands[args[0]]; hasHandler {
						hArg := ""
						if len(args) == 2 {
							hArg = args[1]
						}
						response = handler(hArg)
					}
				}
				log.Printf("Read %d bytes: %q", n, trimmed)
				_, err = c.Write([]byte(response))
				if err != nil {
					log.Print("error writing bytes", err)
				}
				if trimmed == "STOP" {
					break
				}
				if trimmed == "CLOSE" {
					closeServer <- true
					l.Close()
					break
				}
			}
		}(conn)
	}
}

func main() {

	s := newServer()
	defer s.close()
	s.tcpServerListen()

}
