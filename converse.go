package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	var host = flag.String("h", "localhost", "host to connect to")
	var port = flag.String("p", "8080", "port to communicate with")
	flag.Parse()

	stop := make(chan struct{})
	go spinner(stop) // give user feedback
	conn, err := net.Dial("tcp", *host+":"+*port)
	stop <- struct{}{}
	if err != nil {
		log.Fatalf("client: %v", err)
	}
	done := make(chan struct{})
	go func() { // receive from connection
		if _, err := io.Copy(os.Stdout, conn); err != nil {
			log.Fatalf("client: receive: %v", err)
		}
		done <- struct{}{}
	}()
	if _, err := io.Copy(conn, os.Stdin); err != nil { // send to connection
		log.Fatalf("client: send: %v", err)
	}
	if tconn, ok := conn.(*net.TCPConn); ok {
		tconn.CloseWrite() // only close write half so program continues to print final reads
	}
	<-done // wait for read goroutine to finish
}

func spinner(stop <-chan struct{}) {
	for {
		select {
		case <-stop:
			fmt.Println()
			return
		default:
			for _, r := range `-\|/` {
				fmt.Printf("\rConnecting...%c", r)
				time.Sleep(80 * time.Millisecond)
			}
		}
	}
}
