package main

import (
    "fmt"
    "io"
    "net"
    "sync"
    "time"
)

const timeout = 10 * time.Minute

func main() {
    listener, err := net.Listen("tcp", "localhost:8888")
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    fmt.Println("Proxy listening on localhost:8888")

    for {
        client, err := listener.Accept()
        if err != nil {
            fmt.Println("Connection error:", err)
            continue
        }
        go handleConnection(client)
    }
}

func handleConnection(client net.Conn) {
    defer client.Close()
    fmt.Printf("Client connected: %s\n", client.RemoteAddr())

    server, err := net.Dial("tcp", "localhost:80")
    if err != nil {
        fmt.Println("Could not connect to server:", err)
        return
    }
    defer server.Close()

    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        n, err := io.Copy(server, client)
        if err != nil {
            fmt.Printf("Error forwarding client->server: %v\n", err)
        }
        fmt.Printf("Forwarded %d bytes from client to server\n", n)
        server.(*net.TCPConn).CloseWrite()
    }()

    go func() {
        defer wg.Done()
        n, err := io.Copy(client, server)
        if err != nil {
            fmt.Printf("Error forwarding server->client: %v\n", err)
        }
        fmt.Printf("Forwarded %d bytes from server to client\n", n)
        client.(*net.TCPConn).CloseWrite()
    }()

    wg.Wait()
    fmt.Printf("Disconnected: %s\n", client.RemoteAddr())
}