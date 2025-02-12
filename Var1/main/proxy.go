package main

import (
    "bufio"
    "fmt"
    "net"
    "time"
)

const proxyTimeout = 30 * time.Second

func main() {
    listener, err := net.Listen("tcp", "localhost:8888")
    if err != nil {
        panic(err)
    }
    defer listener.Close()

    fmt.Println("Proxy is running on localhost:8888")

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

    fmt.Printf("New client connected: %s\n", client.RemoteAddr())

    server, err := net.Dial("tcp", "localhost:80")
    if err != nil {
        fmt.Printf("Unable to connect to server: %v\n", err)
        return
    }
    defer server.Close()

    fmt.Printf("Connected to target server: %s\n", server.RemoteAddr())

    done := make(chan bool)

    go relay(client, server, done)
    go relay(server, client, done)

    <-done
    fmt.Printf("Client disconnected: %s\n", client.RemoteAddr())
}

func relay(src net.Conn, dst net.Conn, done chan bool) {
    reader := bufio.NewScanner(src)
    for reader.Scan() {
        data := reader.Text()
        fmt.Fprintf(dst, "%s\n", data)
        src.SetDeadline(time.Now().Add(proxyTimeout))
        dst.SetDeadline(time.Now().Add(proxyTimeout))
    }
    if err := reader.Err(); err != nil {
        fmt.Printf("Relay error: %v\n", err)
    }
    done <- true
}
