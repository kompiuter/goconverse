// Package client allows safe concurrent access to clients
package client

import (
	"sync"
)

// Client represents a client connected to the server
type Client struct {
	Ch      chan string // (outgoing) channel where client receives messages
	Name    string
	Address string
}

// Clients are all the clients that are connected to a server
type Clients map[Client]bool

var rw = sync.RWMutex{}

// New returns a new empty struct of type Clients
func New() Clients {
	return make(Clients)
}

// Add adds a client to the Clients connection
func (c Clients) Add(cl Client) {
	rw.Lock()
	defer rw.Unlock()
	c[cl] = true
}

// Inform sends the message msg to all clients
func (c Clients) Inform(msg string) {
	rw.RLock()
	defer rw.RUnlock()
	for cl := range c {
		cl.Ch <- msg
	}
}

// Remove deletes a client from the client connection and closes its outgoing channel
func (c Clients) Remove(cl Client) {
	rw.Lock()
	defer rw.Unlock()
	delete(c, cl)
	close(cl.Ch)
}

// Exists returns true if a client with the give name already exists
func (c Clients) Exists(name string) bool {
	rw.RLock()
	defer rw.RUnlock()
	for cl := range c {
		if cl.Name == name {
			return true
		}
	}
	return false
}
