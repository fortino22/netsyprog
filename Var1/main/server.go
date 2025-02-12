package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
	"var1/config"
)

const connTimeout = 30 * time.Second

func main() {
    server, err := net.Listen("tcp", "localhost:80")
    if err != nil {
        panic(err)
    }
    defer server.Close()
    fmt.Println("TCP Server is running on localhost:80")

    for {
        client, err := server.Accept()
        if err != nil {
            fmt.Println("Failed to establish connection:", err)
            continue
        }
        client.SetDeadline(time.Now().Add(connTimeout))
        go processClient(client)
    }
}

func processClient(client net.Conn) {
    defer client.Close()
    fmt.Printf("New client connected: %s\n", client.RemoteAddr())

    recvBuffer := make([]byte, config.BaseBufferSize)
    bufReader := bufio.NewReaderSize(client, config.BaseBufferSize)

    for {
        client.SetDeadline(time.Now().Add(connTimeout))
        bytesRead, err := bufReader.Read(recvBuffer)
        if err != nil {
            if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                fmt.Printf("Connection timed out for client: %s\n", client.RemoteAddr())
            }
            return
        }

        fmt.Printf("Current buffer: %d bytes, Received data: %d bytes\n", len(recvBuffer), bytesRead)
        resizeBuffer(&recvBuffer, &bufReader, client, bytesRead)
        fmt.Printf("Message from %s: %s\n", client.RemoteAddr(), string(recvBuffer[:bytesRead]))
    }
}

func resizeBuffer(recvBuffer *[]byte, bufReader **bufio.Reader, client net.Conn, bytesRead int) {
    if bytesRead == len(*recvBuffer) && len(*recvBuffer) < config.MaxBufferSize {
        newBufferSize := len(*recvBuffer) * 2
        if newBufferSize > config.MaxBufferSize {
            newBufferSize = config.MaxBufferSize
        }
        *recvBuffer = make([]byte, newBufferSize)
        *bufReader = bufio.NewReaderSize(client, newBufferSize)
        fmt.Printf("Buffer expanded to: %d bytes\n", len(*recvBuffer))
    }
}
