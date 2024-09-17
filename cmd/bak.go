package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func handleClient(conn net.Conn) {
	defer conn.Close()
	// Create a buffer to read data into
	// buffer := make([]byte, 1024)
	recReader := bufio.NewReader(conn)
	for {
		// Read data from the client
		recData, err := recReader.ReadBytes('\n')
		// n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		// Process and use the data (here, we'll just print it)
		fmt.Printf("Received: %s", recData)
	}
}

func connectTo2(addr string) {
	// Connect to the server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	c := make(chan string)

	// Read and process data from the server
	go readFrom("client", conn, c)

	// Send data to the server
	reader := bufio.NewReader(os.Stdin)
	for {
		data, err := reader.ReadBytes('\n')
		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		select {
		case e := <-c: // read from c - nothing there, so block
			fmt.Println("connection dropped", e)
			return
		case <-time.After(500 * time.Millisecond):
			fmt.Println("timeout")
		}

		_, err = conn.Write(data)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
}

func readFrom(whoami string, conn net.Conn, c chan string) {
	defer conn.Close()

	respReader := bufio.NewReader(conn)
	for {
		// Read data from connection, line-by-line
		respData, err := respReader.ReadBytes('\n')
		// n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error:", err)
			c <- err.Error()
			return
		}

		// Process and use the data (here, we'll just print it)
		fmt.Printf("%s received: %s", whoami, respData)
	}
}

func ChanReadBytes(c chan byte, delim byte, f func(part []byte)) {
	buf := make([]byte, 1024)

	for b := range c {
		buf = append(buf, b)

		if b == delim {
			f(buf)
			clear(buf)
		}
	}

	if len(buf) > 0 {
		f(buf)
	}
}

func ChanWriteBytes(c chan byte, b []byte) {
	for _, b := range b {
		c <- b
	}

	close(c)
}
