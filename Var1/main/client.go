package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "time"
    "var1/config"
)

func main() {
    connection, err := net.Dial("tcp", "localhost:8888")
    if err != nil {
        panic(err)
    }
    defer connection.Close()

    fmt.Println("Successfully connected to the proxy server.")
    connection.SetDeadline(time.Now().Add(30 * time.Second))

    keyboardScanner := bufio.NewScanner(os.Stdin)
    messageBuffer := make([]byte, config.BaseBufferSize)
    messageWriter := bufio.NewWriterSize(connection, config.BaseBufferSize)

    for {
        fmt.Print("Type your message: ")
        keyboardScanner.Scan()
        userMessage := keyboardScanner.Text()

        fmt.Printf("Current message buffer size: %d bytes\n", len(messageBuffer))
        connection.SetDeadline(time.Now().Add(30 * time.Second))

        if len(userMessage) > len(messageBuffer) && len(messageBuffer) < config.MaxBufferSize {
            newBufferSize := len(messageBuffer) * 2
            if newBufferSize > config.MaxBufferSize {
                newBufferSize = config.MaxBufferSize
            }
            messageBuffer = make([]byte, newBufferSize)
            messageWriter = bufio.NewWriterSize(connection, newBufferSize)
            fmt.Printf("Buffer size expanded to: %d bytes\n", len(messageBuffer))
        }

        _, err := fmt.Fprintf(messageWriter, userMessage+"\n")
        if err != nil {
            if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
                fmt.Println("The connection has timed out.")
            } else {
                fmt.Println("Failed to send the message:", err)
            }
            return
        }
        messageWriter.Flush()
    }
}
