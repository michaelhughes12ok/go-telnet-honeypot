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
	
	defer func() {
		if err := file.Close(); err != nil {
        		log.Fatal(err)
    		}
	}()

	go func() {
		for msg := range logChannel {
			if _, err := file.Write(msg); err != nil {
				log.Println(err)
			}
		}
	}()

	serverPort := 2323
	serverIp := "0.0.0.0"

	addr := fmt.Sprintf("%s:%v", serverIp, serverPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer func(){
		if err := listener.Close(); err != nil {
			log.Fatal(err)
		}
	}()

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
	defer func(){
		if err := conn.Close(); err != nil {
			log.Println(err)
		}
	}()
	if _, err := conn.Write([]byte("\r\nPassword: ")); err != nil {
		log.Println(err)
		return
	}
	buffer := make([]byte, 1024)
	for {
		nInt, err := conn.Read(buffer)
		if err != nil {
			return
		}

		attemptedInput := string(buffer[:nInt])
		result := time.Now().Format(time.RFC3339)
		sourceAddr := conn.RemoteAddr().String()

		attackLog := Attempt{Time: result, SourceAddr: sourceAddr, AttemptedInput: attemptedInput}

		jsonData, err := json.Marshal(attackLog)
		if err != nil {
			log.Println(err)
			continue
		}

		jsonData = append(jsonData, '\n')
		logChannel <- jsonData

		if _, err := conn.Write([]byte("\r\nPassword incorrect\r\nPassword: ")); err != nil{
			log.Println(err)
		return
		}
	}
}
