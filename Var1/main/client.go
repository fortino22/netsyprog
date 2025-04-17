package main

import (
    "bufio"
    "fmt"
    "io"
    "net"
    "os"
    "strings"
    "time"
)

type Payload interface {
    GetContent() string
    GetSize() int
    Send(writer io.Writer) (int, error)
}

type TextPayload struct {
    Content string
}

func (p TextPayload) GetContent() string {
    return p.Content
}

func (p TextPayload) GetSize() int {
    return len(p.Content)
}

func (p TextPayload) Send(writer io.Writer) (int, error) {
    messageWithNewline := p.Content + "\n"
    return fmt.Fprint(writer, messageWithNewline)
}

func main() {
    conn, err := net.Dial("tcp", "localhost:8888")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    conn.SetDeadline(time.Now().Add(10 * time.Minute))
    fmt.Println("Connected to proxy server.")

    keyboardScanner := bufio.NewScanner(os.Stdin)

    for {
        fmt.Print("Enter your message (type 'exit' to quit): ")
        keyboardScanner.Scan()
        message := keyboardScanner.Text()

        if strings.ToLower(message) == "exit" {
            fmt.Println("Exiting the program.")
            break
        }

        payload := TextPayload{Content: message}
        
        bytesSent, err := payload.Send(conn)
        if err != nil {
            fmt.Println("Error sending message:", err)
            break
        }

        fmt.Printf("Message sent (%d bytes)\n", bytesSent)
    }
}