package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"sync"

	"github.com/slayer321/chatApp/proto"

	"google.golang.org/grpc"
)

var client proto.BroadcastClient
var wait *sync.WaitGroup

// syn the channel using init function

func init() {
	wait = &sync.WaitGroup{}
}

// Connection the client with the open port

func connect(user *proto.User) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &proto.Connect{
		User:   user,
		Active: true,
	})

	if err != nil {
		return fmt.Errorf("connection failed : %v", err)
	}

	wait.Add(1)
	go func(str proto.Broadcast_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("Error reading message: %v", err)
				break
			}
			fmt.Printf("%v : %s\n", msg.Id, msg.Content)
		}
	}(stream)

	return streamerror
}

func main() {

	timestamp := time.Now()
	done := make(chan int)

	name := flag.String("n", "sachin", "The name of the user")
	flag.Parse()
	//id := sha256.Sum256([]byte(timestamp.String() + *name))
	id := name

	//Listing on port 8080 over a Insecure medium

	conn, err := grpc.Dial("localhost:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Couldnt connect to service: %v", err)
	}

	client = proto.NewBroadcastClient(conn)

	// adding new Client

	user := &proto.User{
		Id: *id,
		//Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	// Calling the Connect function
	connect(user)

	wait.Add(1)

	// Chatting
	// taking the user input from the user
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &proto.Message{
				Id:        user.Id,
				Content:   scanner.Text(),
				Timestamp: timestamp.String(),
			}
			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error Sending Message : %v", err)
				break
			}
		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
}
