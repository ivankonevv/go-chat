package chatserver

import (
	"log"
	"math/rand"
	"sync"
	"time"
)

type messageUnit struct {
	ClientName        string
	MessageBody       string
	MessageUniqueCode int
	ClientUniqueCode  int
}

type messageHandle struct {
	MQue []messageUnit
	mu   sync.Mutex
}

var messageHandleObject = messageHandle{}

type ChatServer struct {
}

//define ChatService
func (is *ChatServer) ChatService(csi Service_ChatServiceServer) error {
	clientUniqueCode := rand.Intn(1e6)
	errch := make(chan error)

	//receive messages
	go receiveFromStream(csi, clientUniqueCode, errch)

	//send messages
	go sendToStream(csi, clientUniqueCode, errch)

	return <-errch
}

func receiveFromStream(csi_ Service_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	//implement a loop
	for {
		mssg, err := csi_.Recv()
		if err != nil {
			log.Printf("Error in receiving message from client :: %v", err)
			errch_ <- err
		} else {
			messageHandleObject.mu.Lock()

			messageHandleObject.MQue = append(messageHandleObject.MQue, messageUnit{
				ClientName:        mssg.Name,
				MessageBody:       mssg.Body,
				MessageUniqueCode: rand.Intn(1e8),
				ClientUniqueCode:  clientUniqueCode_,
			})

			messageHandleObject.mu.Unlock()

			log.Printf("%v", messageHandleObject.MQue[len(messageHandleObject.MQue) - 1])
		}
	}
}

//send message
func sendToStream(csi_ Service_ChatServiceServer, clientUniqueCode_ int, errch_ chan error) {
	//implement a loop
	for{
		//loop through messages in MQue
		for{
			time.Sleep(500 * time.Millisecond)

			messageHandleObject.mu.Lock()

			if len(messageHandleObject.MQue) == 0 {
				messageHandleObject.mu.Unlock()
				break
			}

			senderUniqueCode := messageHandleObject.MQue[0].ClientUniqueCode
			senderNameForClient := messageHandleObject.MQue[0].ClientName
			messageForClient := messageHandleObject.MQue[0].MessageBody

			messageHandleObject.mu.Unlock()

			//send message to designated client
			if senderUniqueCode != clientUniqueCode_ {
				err := csi_.Send(&FromServer{Name: senderNameForClient, Body: messageForClient})
				if err != nil {
					errch_ <- err
				}

				messageHandleObject.mu.Lock()

				if len(messageHandleObject.MQue) > 1 {
					messageHandleObject.MQue  = messageHandleObject.MQue[1:] //delete the message at index 0 after sending to receiver
				} else {
					messageHandleObject.MQue = []messageUnit{}
				}

				messageHandleObject.mu.Unlock()
			}
		}

		time.Sleep(100 * time.Millisecond)
	}
}