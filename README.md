# goconverse
A simplistic console based chat application which works cross-platform based on the TCP protocol. 

# Installation
This requires a working Go environment to run. Follow the steps [here](http://golang.org/doc/install) to install the Go environment.

Once Go is running, you can download and build the application using the following command:

<code>go get github.com/kompiuter/goconverse</code>

Make sure you build separately on each different OS you are using.

Executables can then be found under
<code>%GOPATH%\bin</code>

# Usage
## Server
Start a new server instance with the following command: 

<code>server.exe -n "ServerName" -p "Port" &</code>

You may enable verbose logging using the -v flag:

<code>server.exe -v -n "ServerName" -p "Port" &</code>

The new server will sit and wait for incoming TCP connections by clients on the specified port. The server will spawn (2n + 2) goroutines,
where n represents the number of clients. 

One goroutine is listening for new clients, one goroutine is acting as a broadcaster to all clients, 
one goroutine (for each client) handles a new client connection, one goroutine (for each client) for sending messages to a client.

## Client
You may connect to the server with the following command:
<code>converse.exe -h "HostName" -p "Port"</code>

When successfully connected to a server you will receive a message asking for a name.
Once a valid name is provided (one that does not already exist in the server), you will join the chat and be able to 
converse with other connected users.
