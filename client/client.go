package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

func sendCommand(c net.Conn, cmd string) (string, bool) {

	reader := bufio.NewReader(c)
	if reader == nil {
		log.Fatal("Could not open reader")
	}
	log.Println("Sending data.")
	_, err := (c).Write([]byte(cmd + "\n"))
	if err != nil {
		log.Fatalf("error writing to pipe %s\n", err)
	}
	log.Println("Reading data.")
	s, err := reader.ReadString('\n')
	if err != nil {
		if err == io.EOF {
			log.Println("Read EOF on connection")
			return "", true
		}
		log.Fatalf("Unknown error on read: %s\n", err)
	}
	log.Printf("Received data: \"%s\"", strings.TrimSpace(s))

	return s, false
}

func resetCounter() {

	c, err := net.Dial("tcp", "127.0.0.1:2000")
	if err != nil {
		log.Fatalf("Could not connect: %s\n", err)
	}
	defer c.Close()
	fmt.Println("Connection established.")
	if _, stop := sendCommand(c, "reset"); stop {
		log.Println("stopped")
	}
	if _, stop := sendCommand(c, "STOP"); stop {
		log.Println("stopped")
	}
}

func getChannelCount() {

	c, err := net.Dial("tcp", "127.0.0.1:2000")
	if err != nil {
		log.Fatalf("Could not connect: %s\n", err)
	}
	defer c.Close()
	fmt.Println("Connection established.")
	if _, stop := sendCommand(c, "getChannelCount total"); stop {
		log.Println("stopped")
	}
	if _, stop := sendCommand(c, "STOP"); stop {
		log.Println("stopped")
	}
}

func connect(i int) {
	c, err := net.Dial("tcp", "127.0.0.1:2000")
	if err != nil {
		log.Fatalf("Could not connect: %s\n", err)
	}
	defer c.Close()
	fmt.Println("Connection established.")

	if _, stop := sendCommand(c, "echo hi"); stop {
		log.Println("stopped")
	}
	if _, stop := sendCommand(c, "wait 1"); stop {
		log.Println("stopped")
	}
	if _, stop := sendCommand(c, "add"); stop {
		log.Println("stopped")
	}
	if _, stop := sendCommand(c, fmt.Sprintf("addChannel %d", i)); stop {
		log.Println("stopped")
	}
	if _, stop := sendCommand(c, "STOP"); stop {
		log.Println("stopped")
	}

}

func main() {
	resetCounter()
	var wg sync.WaitGroup
	n := 1000
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(i int) {
			defer wg.Done()
			connect(i)
		}(i)
	}
	wg.Wait()
	getChannelCount()
}
