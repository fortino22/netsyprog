package main

import (
    "bufio"
    "fmt"
    "io"
    "net"
    "time"
)

const timeout = 10 * time.Minute

func main() {
    listener, err := net.Listen("tcp", "localhost:80")
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    fmt.Println("Server listening on localhost:80")

    for {
        conn, err := listener.Accept()
        if err != nil {
            fmt.Println("Accept error:", err)
            continue
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    fmt.Printf("Client connected: %s\n", conn.RemoteAddr())
    conn.SetDeadline(time.Now().Add(timeout))

    reader := bufio.NewReader(conn)
    
    for {
        message, err := reader.ReadBytes('\n')
        
        if err != nil {
            if err == io.EOF {
                if len(message) > 0 {
                    processMessage(message)
                }
                fmt.Println("Connection closed by client")
            } else {
                fmt.Println("Read error:", err)
            }
            break
        }
        
        processMessage(message[:len(message)-1])
    }
}

func processMessage(message []byte) {
    bytesReceived := len(message)
    fmt.Printf("Received message (%d bytes): %s\n", bytesReceived, string(message))
}