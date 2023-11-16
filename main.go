package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

func createNode(node1 string) host.Host {
	node, err := libp2p.New()
	if err != nil {
		panic(err)
	}

	//you can generate a key here
	//createPairPublicPrivateKey
	d := node1
	log.Printf("Private key: %s", d)

	return node
}

func readHelloProtocol(s network.Stream) error {
	//   Read the stream and print its content
	buf := bufio.NewReader(s)
	message, err := buf.ReadString('\n')
	if err != nil {
		return err
	}

	connection := s.Conn()

	log.Printf("Message from '%s': %s", connection.RemotePeer().String(), message)
	return nil
}

func runTargetNode() peer.AddrInfo {
	log.Printf("Creating target node...")
	targetNode := createNode("0xadd613596b13c8695004309e960f31a3e596d4bfb3d457e5f716076c8c8c5df8d0e")
	log.Printf("Target node created with ID '%s'", targetNode.ID().String())

	// Set stream handler for the "/hello/1.0.0" protocol
	targetNode.SetStreamHandler("/hello/1.0.0", func(s network.Stream) {
		log.Printf("/hello/1.0.0 stream created")
		err := readHelloProtocol(s)
		if err != nil {
			s.Reset()
		} else {
			s.Close()
		}
	})

	return *host.InfoFromHost(targetNode)
}

func runSourceNode(targetNodeInfo peer.AddrInfo) {
	log.Printf("Creating source node...")
	sourceNode := createNode("0xadd613596b13c8695004309e960f31a3e596d4bfb3d457e5f716076c8c8c5df8d01")
	log.Printf("Source node created with ID '%s'", sourceNode.ID().String())
	//print targetNodeInfo
	log.Printf("Target node info: %s %s", targetNodeInfo, targetNodeInfo.ID)
	//convert string address to peer.AddrInfo
	sourceNode.Connect(context.Background(), targetNodeInfo)

	// Open stream and send message
	// we listen on this stream
	stream, err := sourceNode.NewStream(context.Background(), targetNodeInfo.ID, "/hello/1.0.0")
	if err != nil {
		panic(err)
	}

	//we send messages on this channel
	message := "Hello World!\n"
	log.Printf("Sending message...")
	_, err = stream.Write([]byte(message))
	if err != nil {
		log.Println(err)
		panic(err)
	}
}

//dispatcher control server spawn
//parameter arg function
func dispatcherControl() {
	argsWithoutProg := os.Args[1:]

	arg := os.Args[1]

	fmt.Println(argsWithoutProg)
	fmt.Println(arg)

}

func main() {
	ctx, _ := context.WithCancel(context.Background())

	// make this 2 nodes to communicate using messages
	// a chat between them
	// also use your own input
	info := runTargetNode()
	runSourceNode(info)
	//dispatcherControl()

	<-ctx.Done()
}
