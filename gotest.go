package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	file, err := os.OpenFile("log.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
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
	fmt.Print("DANGER!!! PLEASE MAKE SURE THIS APPLICATION IS CONTAINERIZED AND RUN AS ROOT/SUDO IN SAID CONTAINER, thank you :).\n")
	fmt.Print("Expose to (1) localhost only (2) all interfaces: ")
	var choice string
	fmt.Scanln(&choice)

	serverIp := "127.0.0.1"
	if choice == "2" {
		serverIp = "0.0.0.0"
	}

	addr := fmt.Sprintf("%s:%v", serverIp, serverPort)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf("Server running on %s...\r", addr)

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
	conn.Write([]byte("BusyBox v1.19.4\r\nlogin: "))
	buffer := make([]byte, 1024)
	for {
		nInt, err := conn.Read(buffer)
		if err != nil {
			return
		}
		attempt := string(buffer[:nInt])

		layout := "2006-01-02 15:04:05"
		result := time.Now().Format(layout)
		sourceAddr := conn.RemoteAddr().String()

		var attackLog string = fmt.Sprintf("%s - Source:%s, Attempt:%s", result, sourceAddr, attempt)
		logChannel <- attackLog

		conn.Write([]byte("\r\nLogin incorrect\r\nlogin: "))
	}
}
