package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/pion/webrtc/v3"
)

func readSDP(prompt string) string {
	fmt.Println(prompt)
	fmt.Println("(Paste SDP and type 'END' on a new line to finish):")

	var sdpLines []string
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "END" {
			break
		}
		sdpLines = append(sdpLines, line)
	}

	if len(sdpLines) == 0 {
		panic("No SDP provided")
	}

	return strings.Join(sdpLines, "\n")
}

func getConfig() webrtc.Configuration {
	return webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:turn.cloudflare.com:3478"},
			},
			{
				URLs: []string{"stun:stunserver2024.stunprotocol.org:3478"},
			},
			{
				URLs: []string{"stun:stun.isp.net.au:3478"},
			},
			{
				URLs: []string{"stun:stun.freeswitch.org:3478"},
			},
			{
				URLs: []string{"stun:stun.voip.blackberry.com:3478"},
			},
		},
	}
}

func main() {
	fmt.Println("Choose role (offer/answer):")
	reader := bufio.NewReader(os.Stdin)
	role, _ := reader.ReadString('\n')
	role = strings.TrimSpace(role)

	if role == "offer" {
		runOffer()
	} else if role == "answer" {
		runAnswer()
	} else {
		fmt.Println("Invalid role")
	}
}

func runOffer() {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	dataChannel.OnOpen(func() {
		fmt.Println("Data channel opened!")
		dataChannel.SendText("Hello from offer!")
	})

	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
		fmt.Printf("Received message: %s\n", string(msg.Data))
	})

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("ICE State: %s\n", state.String())
	})

	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	if err := peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	fmt.Println("SDP Offer (paste to answer):")
	fmt.Println(offer.SDP)

	answerSDP := readSDP("Paste SDP Answer:")
	answer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeAnswer,
		SDP:  answerSDP,
	}
	if err := peerConnection.SetRemoteDescription(answer); err != nil {
		panic(err)
	}

	select {}
}

func runAnswer() {
	peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})
	if err != nil {
		panic(err)
	}

	peerConnection.OnDataChannel(func(dc *webrtc.DataChannel) {
		fmt.Printf("DataChannel opened: %s\n", dc.Label())
		dc.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("Received message: %s\n", string(msg.Data))
		})
		dc.OnOpen(func() {
			dc.SendText("Hi from answer!")
		})
	})

	peerConnection.OnICEConnectionStateChange(func(state webrtc.ICEConnectionState) {
		fmt.Printf("ICE State: %s\n", state.String())
	})

	offerSDP := readSDP("Paste SDP Offer:")
	offer := webrtc.SessionDescription{
		Type: webrtc.SDPTypeOffer,
		SDP:  offerSDP,
	}
	if err := peerConnection.SetRemoteDescription(offer); err != nil {
		panic(err)
	}

	answer, err := peerConnection.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	if err := peerConnection.SetLocalDescription(answer); err != nil {
		panic(err)
	}

	fmt.Println("SDP Answer (paste to offer):")
	fmt.Println(answer.SDP)

	select {}
}
