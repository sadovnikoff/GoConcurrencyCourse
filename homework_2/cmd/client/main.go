package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_2/internal/common"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP network address")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxBufSize := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	bufSize, err := common.ParseBufSize(*maxBufSize)
	if err != nil {
		fmt.Println("Error parsing maxBufSize:", err)
		return
	}

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer func() { _ = conn.Close() }()

	if *idleTimeout != 0 {
		if err := conn.SetDeadline(time.Now().Add(*idleTimeout)); err != nil {
			fmt.Println("Failed to set deadline for connection:", err)
			return
		}
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type a command with arguments and press Enter (available commands: SET, GET, DEL).")
	for {
		fmt.Print("> ")
		request, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading request: %s", err)
		}

		_, err = conn.Write([]byte(request))
		if err != nil {
			fmt.Println("Error sending request:", err)
			return
		}

		response := make([]byte, bufSize)
		count, err := conn.Read(response)
		if err != nil && err != io.EOF {
			fmt.Println("Error reading server response:", err)
			break
		} else if count == bufSize {
			fmt.Println("Error reading server response: too small buffer size")
			break
		}

		fmt.Println(string(response[:count]))
	}
}
