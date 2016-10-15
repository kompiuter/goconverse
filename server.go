package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"

	c "github.com/kompiuter/goconverse/client"
)

var (
	entering = make(chan c.Client) // clients connecting
	leaving  = make(chan c.Client) // clients disconnecting
	messages = make(chan string)   // all messages from clients

	serverName = flag.String("n", "Converse", "server name")
	verbose    = flag.Bool("v", false, "Enable or disable verbose logging")

	clients = c.New()
)

func main() {
	var port = flag.String("p", "8080", "port to listen on")
	flag.Parse()

	ln, err := net.Listen("tcp", "localhost:"+*port)
	if err != nil {
		log.Fatalf("server: %v", err)
	}

	go hub()
	for {
		conn, err := ln.Accept()
		serverLog("Incoming connection", 1)
		if err != nil {
			serverLog(fmt.Sprintf("server: %v", err), 0)
			continue
		}
		go handleClient(conn)
	}

}

func serverLog(msg string, level int) {
	if level == 0 {
		log.Println("::" + msg)
	} else {
		if *verbose {
			log.Println("::" + msg)
		}
	}
}

func handleClient(conn net.Conn) {
	in := bufio.NewScanner(conn)
	who := getName(conn, in)
	cl := c.Client{Ch: make(chan string), Name: who, Address: conn.RemoteAddr().String()}
	go messageWriter(conn, cl.Ch) // goroutine that listens for messages directed at client

	msg := fmt.Sprintf("Welcome to %s! You are %s, connected from %s", *serverName, cl.Name, cl.Address)
	cl.Ch <- msg
	entering <- cl
	messages <- fmt.Sprintf("%s has connected to the server", cl.Name)

	for in.Scan() {
		broadcastMsg(cl, in.Text())
	}
	if err := in.Err(); err != nil {
		serverLog(fmt.Sprintf("handleClient: input: %v", err), 0)
	}

	leaving <- cl
	messages <- cl.Name + " has left the server"
	conn.Close()
}

func getName(conn net.Conn, in *bufio.Scanner) (who string) {
	who = ""
	for who == "" || clients.Exists(who) {
		fmt.Fprintf(conn, "Enter your name: ")
		if in.Scan() {
			who = in.Text()
		}
		if clients.Exists(who) {
			fmt.Fprintf(conn, "Name provided already exists in the server\n")
		}
	}
	return
}

func broadcastMsg(from c.Client, msg string) {
	serverLog(fmt.Sprintf("%s has sent a message", from.Address), 1)
	messages <- from.Name + ": " + msg
}

func messageWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		_, err := fmt.Fprintln(conn, "\r"+msg)
		if err != nil {
			log.Printf("sending message: %v", err)
			break
		}
	}
}

func hub() {
	for {
		select {
		case msg := <-messages:
			clients.Inform(msg)
		case cl := <-entering:
			clients.Add(cl)
			serverLog(fmt.Sprintf("%s has connected", cl.Address), 1)
		case cl := <-leaving:
			clients.Remove(cl)
			serverLog(fmt.Sprintf("%s has disconnected", cl.Address), 1)
		}

	}
}
