package main

//Make a program that can send and receive messages on the network
//Using the bcast package

import ( 
	"github.com/SondreHauk/TTK4145-Sanntid/network/bcast"
)

func main() {
	//Create a channel for sending and receiving messages
	sendChan := make(chan string)
	receiveChan := make(chan string)
	
	//Start a transmitter and receiver
	go bcast.Transmitter(16569, sendChan)
	go bcast.Receiver(16569, receiveChan)
	
	//Send a message
	sendChan <- "Hello world!"
	sendChan <- "Hello again!"
	sendChan <- "Hello for the third time!"
	
	//Receive a message
	message := <- receiveChan
	println(message)
	message = <- receiveChan
	println(message)
	message = <- receiveChan
	println(message)
}