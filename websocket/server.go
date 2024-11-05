package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type ProgressUpdate struct {
	AssetID string `json:"assetID"`
	Status  string `json:"status"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Store active connections with a mutex for thread-safe access
var connections = make(map[*websocket.Conn]bool)
var connMutex = sync.Mutex{}

func getClientIP(r *http.Request) string {
	// Check if "X-Forwarded-For" header is present
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// The client IP is typically the first IP in the X-Forwarded-For list
		ips := strings.Split(forwarded, ",")
		return strings.TrimSpace(ips[0])
	}
	// Fallback to RemoteAddr if X-Forwarded-For is not set
	return r.RemoteAddr
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("HandleWebSocket")
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade WebSocket:", err)
		return
	}
	// defer conn.Close()

	// Add connection to the active connections map
	connMutex.Lock()
	connections[conn] = true
	connectionID := fmt.Sprintf("%p", &conn) // Unique identifier for logging

	clientIP := getClientIP(r)
	userAgent := r.UserAgent()
	startTime := time.Now()

	log.Printf("New connection [%s] from IP: %s, User-Agent: %s at %s", connectionID, clientIP, userAgent, startTime.Format(time.RFC3339))

	log.Printf("New WebSocket connection: %p", conn)
	log.Printf("Total active connections: %d", len(connections))

	connMutex.Unlock()

	// Get assetID from query parameters (for tracking purposes)
	assetID := r.URL.Query().Get("assetID")
	fmt.Println("assetID", assetID)
	if assetID == "" {
		log.Println("Missing assetID")
		return
	}

	// Start the asset generation process in a separate goroutine
	go generateAsset(assetID, conn)
}

func generateAsset(assetID string, conn *websocket.Conn) {
	fmt.Println("generating Asset...")
	// Simulate a process that takes a few seconds to complete
	time.Sleep(5 * time.Second) // Replace this with actual generation logic

	// Create a completion message
	update := ProgressUpdate{
		AssetID: assetID,
		Status:  "Asset generation complete",
	}

	fmt.Println("updating...")

	// Send the completion message to the client
	sendUpdate(conn, update)
}

func sendUpdate(conn *websocket.Conn, update ProgressUpdate) {
	fmt.Println("sending update")
	// Marshal the update message to JSON format
	message, err := json.Marshal(update)
	if err != nil {
		log.Println("Error encoding message:", err)
		return
	}

	// Send the message to the WebSocket client
	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		fmt.Println("error sending message??")

		log.Println("Error sending message:", err)
	}
	fmt.Println("message sent?")

}
