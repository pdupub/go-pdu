package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pion/webrtc/v3"
)

func main() {
	fmt.Println("Choose role (offer/answer):")
	reader := bufio.NewReader(os.Stdin)
	role, _ := reader.ReadString('\n')
	role = role[:len(role)-1]

	if role == "offer" {
		runOffer()
	} else if role == "answer" {
		runAnswer()
	} else {
		fmt.Println("Invalid role")
	}
}

func runOffer() {
	// Create a new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			// {URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed to %s\n", state.String())
	})

	// Create a data channel for message exchange
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Handle data channel open event
	dataChannel.OnOpen(func() {
		fmt.Println("Data channel opened!")
	})

	// Handle data channel message event
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Received message: %s\n", string(msg.Data))
		dataChannel.SendText("Hi") // Respond with "Hi"
	})

	// Generate the SDP Offer
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}
	if err := peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	// Output the SDP Offer
	fmt.Println("SDP Offer (copy this to the Answer side):")
	fmt.Println(offer.SDP)
	fmt.Println("END")

	// Read the SDP Answer
	answerSDP := readSDP("Paste SDP Answer:")
	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answerSDP,
	}
	if err := peerConnection.SetRemoteDescription(answer); err != nil {
		panic(err)
	}

	select {} // Keep the program running
}

func runAnswer() {
	// Create a new PeerConnection
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{URLs: []string{"stun:stun.l.google.com:19302"}},
		},
	})
	if err != nil {
		panic(err)
	}

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("ICE Connection State has changed to %s\n", state.String())
	})
	// Handle data channel creation
	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		fmt.Printf("New DataChannel: %s\n", dc.Label())

		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Received message: %s\n", string(msg.Data))
		})

		dc.OnOpen(func() {
			fmt.Println("Data channel opened!")
			dc.SendText("Hello") // Send "Hello" to the offer side
		})
	})

	// Read the SDP Offer
	offerSDP := readSDP("Paste SDP Offer:")
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerSDP,
	}

	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	// Generate the SDP Answer
	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}
	if err := peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	// Output the SDP Answer
	fmt.Println("SDP Answer (copy this to the Offer side):")
	fmt.Println(answer.SDP)
	fmt.Println("END")

	select {} // Keep the program running
}

func readSDP(prompt string) string {
	fmt.Println(prompt)
	fmt.Println("(Paste SDP and type 'END' on a new line to finish):")

	var sdpLines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" {
			sdpLines = append(sdpLines, "\n")
			break
		}
		sdpLines = append(sdpLines, line)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return strings.Join(sdpLines, "\n")
}
