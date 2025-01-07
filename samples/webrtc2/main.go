package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Signal struct {
	Type      string `json:"type"`
	SDP       string `json:"sdp,omitempty"`
	Candidate string `json:"candidate,omitempty"`
}

func main() {
	// Connect to signaling server
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		panic("Failed to connect to signaling server: " + err.Error())
	}
	defer conn.Close()

	// Create a new WebRTC peer connection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		panic(err)
	}

	// Create a DataChannel
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Handle incoming messages on the DataChannel
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Received: %s\n", string(msg.Data))
	})

	// Send a message to the peer when the DataChannel is open
	dataChannel.OnOpen(func() {
		fmt.Println("Data channel open. Type a message:")
		go func() {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				dataChannel.SendText(scanner.Text())
			}
		}()
	})

	// Handle ICE candidates
	peerConnection.OnICECandidate(func(candidate *webrtc.ICECandidate) {
		if candidate != nil {
			// Send ICE candidate to signaling server
			signal := Signal{
				Type:      "candidate",
				Candidate: candidate.ToJSON().Candidate,
			}
			sendSignal(conn, signal)
		}
	})

	// Handle incoming messages from signaling server
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Read error:", err)
				return
			}
			handleSignal(conn, peerConnection, message)
		}
	}()

	// Create an SDP offer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Set the local description
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		panic(err)
	}

	// Send SDP offer to signaling server
	sendSignal(conn, Signal{Type: "offer", SDP: offer.SDP})

	// Gracefully close on interrupt
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	<-interrupt

	peerConnection.Close()
}

func sendSignal(conn *websocket.Conn, signal Signal) {
	message, _ := json.Marshal(signal)
	conn.WriteMessage(websocket.TextMessage, message)
}

func handleSignal(conn *websocket.Conn, peerConnection *webrtc.PeerConnection, message []byte) {
	var signal Signal
	if err := json.Unmarshal(message, &signal); err != nil {
		fmt.Println("Failed to parse signal:", err)
		return
	}

	switch signal.Type {
	case "offer":
		// Set the remote description
		err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeOffer,
			SDP:  signal.SDP,
		})
		if err != nil {
			fmt.Println("Failed to set remote description:", err)
			return
		}

		// Create an answer
		answer, err := peerConnection.CreateAnswer(nil)
		if err != nil {
			fmt.Println("Failed to create answer:", err)
			return
		}

		// Set the local description
		peerConnection.SetLocalDescription(answer)

		// Send SDP answer to signaling server
		sendSignal(conn, Signal{Type: "answer", SDP: answer.SDP})

	case "answer":
		// Set the remote description
		err := peerConnection.SetRemoteDescription(webrtc.SessionDescription{
			Type: webrtc.SDPTypeAnswer,
			SDP:  signal.SDP,
		})
		if err != nil {
			fmt.Println("Failed to set remote description:", err)
			return
		}

	case "candidate":
		// Add ICE candidate
		err := peerConnection.AddICECandidate(webrtc.ICECandidateInit{
			Candidate: signal.Candidate,
		})
		if err != nil {
			fmt.Println("Failed to add ICE candidate:", err)
			return
		}
	}
}
