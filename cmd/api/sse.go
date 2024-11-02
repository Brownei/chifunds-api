package api

import (
	"fmt"
	"net/http"
)

// Client represents a connection to be notified
type Client struct {
	channel chan string
}

// global slice to manage connected clients
var clients = make([]*Client, 0)

// Handler to send balance updates
func SseRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Create new client
	client := &Client{channel: make(chan string)}
	clients = append(clients, client)

	// Close connection on client disconnect
	defer func() {
		close(client.channel)
		removeClient(client)
	}()

	// Listen to the client channel and send data as SSE events
	for {
		select {
		case msg := <-client.channel:
			fmt.Fprintf(w, "data: %s\n\n", msg)
			w.(http.Flusher).Flush()
		case <-r.Context().Done():
			return
		}
	}
}

// Function to remove client from the slice
func removeClient(client *Client) {
	for i, c := range clients {
		if c == client {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// Broadcast balance update to all connected clients
func broadcastBalanceUpdate(balance string) {
	for _, client := range clients {
		select {
		case client.channel <- balance:
		default:
			close(client.channel)
			removeClient(client)
		}
	}
}
