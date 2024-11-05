package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"websocket-demo/websocket"
)

func main() {
	http.HandleFunc("/generate", startAssetGeneration)
	http.HandleFunc("/ws", websocket.HandleWebSocket)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port for local testing
	}

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

type GenerateResponse struct {
	WsUrl string `json:"wsUrl"`
}

// Start asset generation and return WebSocket URL
func startAssetGeneration(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		return
	}

	fmt.Println("start generating and respond with new assetID")
	assetID := "12345"
	// wsURL := fmt.Sprintf("ws://localhost:8080/ws?assetID=%s", assetID)
	wsURL := fmt.Sprintf("ws://websocket-demo123-3a8d482bdc30.herokuapp.com/ws?assetID=%s", assetID)

	// Prepare the response with the WebSocket URL
	response := GenerateResponse{
		WsUrl: wsURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
