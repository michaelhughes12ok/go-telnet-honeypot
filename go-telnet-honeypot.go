package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
	"encoding/json"
)

type Attempt struct {
Time string `json:"time"`
SourceAddr string `json:"source_address"`
AttemptedInput string `json:"attempted_input"`
}

func main() {
	jsonData, err := json.Marshal(Attempt)
	if err != nil {
		
	}
	file, err := os.OpenFile()
	if err != nil {
		log.Fatal("File already exists or cannot be created: ", err)
	}

	logChannel := make(chan string, 100)

	go func() {
		for msg := range logChannel {
			file.WriteString(msg)
		}
	}()
	serverPort := 23
	serverIp := "0.0.0.0"

	addr := fmt.Sprintf("%s:%v", serverIp, serverPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn, logChannel)
	}
}

func handleConnection(conn net.Conn, logChannel chan string) {
	defer conn.Close()
	conn.Write([]byte("\r\nPassword: "))
	buffer := make([]byte, 1024)
	for {
		nInt, err := conn.Read(buffer)
		if err != nil {
			return
		}

		attemptedInput := string(buffer[:nInt])
		layout := "2006-01-02 15:04:05"
		result := time.Now().Format(layout)
		sourceAddr := conn.RemoteAddr().String()
		logChannel <- Attempt

		conn.Write([]byte("\r\nPassword incorrect\r\nPassword: "))
	}
}
