package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Attempt struct {
	Time           string `json:"time"`
	SourceAddr     string `json:"source_address"`
	AttemptedInput string `json:"attempted_input"`
}

func main() {
	logChannel := make(chan []byte, 100)
	file, err := os.OpenFile("attempts.jsonl", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	go func() {
		for msg := range logChannel {
			file.Write(msg)
		}
	}()

	serverPort := 2323
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

func handleConnection(conn net.Conn, logChannel chan []byte) {
	defer conn.Close()
	conn.Write([]byte("\r\nPassword: "))
	buffer := make([]byte, 1024)
	for {
		nInt, err := conn.Read(buffer)
		if err != nil {
			return
		}

		attemptedInput := string(buffer[:nInt])
		result := time.Now().Format("2006-01-02 15:04:05")
		sourceAddr := conn.RemoteAddr().String()

		attackLog := Attempt{Time: result, SourceAddr: sourceAddr, AttemptedInput: attemptedInput}

		jsonData, err := json.Marshal(attackLog)
		if err != nil {
			log.Println(err)
			continue
		}

		jsonData = append(jsonData, '\n')
		logChannel <- jsonData

		conn.Write([]byte("\r\nPassword incorrect\r\nPassword: "))
	}
}
