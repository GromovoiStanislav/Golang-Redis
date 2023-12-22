package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	
	"github.com/gorilla/websocket"
)

func main() {
	fmt.Print("Enter your username: ")
	username := getUserInput()

	serverAddr := "ws://localhost:8080/chat/" + username
	conn, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			fmt.Println("Received message:", string(message))
		}
	}()

	for {
		fmt.Print("Enter message (or 'exit' to quit): ")
		message := getUserInput()

		if strings.ToLower(message) == "exit" {
			break
		}

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
			break
		}
	}
}

func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}
