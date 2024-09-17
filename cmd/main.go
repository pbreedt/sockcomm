package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {

	role := os.Args[1]
	if role == "" {
		fmt.Printf("Must specify role")
		os.Exit(1)
	}

	port := os.Args[2]
	if port == "" {
		fmt.Printf("Must specify port")
		os.Exit(1)
	}

	if role == "server" {
		go http.ListenAndServe("localhost:6060", nil)
		serveOn(fmt.Sprintf(":%s", port))
	} else {
		connectTo(fmt.Sprintf(":%s", port))
	}

	wg.Wait()
}

func serveOn(addr string) {

	// Listen for incoming connections
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Server is listening on port", listener.Addr().(*net.TCPAddr).Port)

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		_, err = conn.Write([]byte("welcome\n"))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		wg.Add(2)
		go readFromConn(conn, false)
		go pingOnConn(conn)
	}
}

func pingOnConn(conn net.Conn) {
	defer wg.Done()

	pinger := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-pinger.C:
			conn.Write([]byte("ping\n"))
		}
	}
}

func connectTo(svrAdr string) {
	// connect to server
	conn := connectToServer(svrAdr)

	// read from stdin
	// write to server
	wg.Add(2)
	go readFromWriteTo(os.Stdin, conn)

	// read from server
	go readFromConn(conn, true)
}

func connectToServer(addr string) net.Conn {
	// Connect to the server
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	return conn
}

func readFromWriteTo(reader io.Reader, writer io.Writer) {
	defer wg.Done()

	br := bufio.NewReader(reader)
	for {
		data, err := br.ReadBytes('\n')
		if err != nil {
			fmt.Println("Read error:", err)
			return
		}
		_, err = writer.Write(data)
		if err != nil {
			fmt.Println("Write error:", err)
			return
		}
	}
}

func readFromConn(conn net.Conn, exitOnBreak bool) {
	defer conn.Close()
	defer wg.Done()

	respReader := bufio.NewReader(conn)
	for {
		// Read data line-by-line from connection
		respData, err := respReader.ReadBytes('\n')
		if err != nil {
			fmt.Println("ERR: ", err)
			if exitOnBreak {
				os.Exit(1)
			} else {
				return
			}
		}

		// Use data
		// fmt.Printf("REC: %s", respData)
		go useData(respData)
	}
}

func useData(data []byte) {
	fmt.Printf("RECIEVED DATA: %s", data)
}
