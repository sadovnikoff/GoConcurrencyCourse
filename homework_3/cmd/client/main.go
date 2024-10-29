package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/sadovnikoff/GoConcurrencyCourse/homework_3/internal/network/tcp"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "HTTP network address")
	idleTimeout := flag.Duration("idle_timeout", time.Minute, "Idle timeout for connection")
	maxBufSize := flag.String("max_message_size", "4KB", "Max message size for connection")
	flag.Parse()

	client, err := tcp.NewClient(*addr, *maxBufSize, *idleTimeout)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Type a command with arguments and press Enter (available commands: SET, GET, DEL).")
	for {
		fmt.Print("> ")
		request, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading request: %s", err)
		}

		response, err := client.Communicate(request)
		if errors.Is(err, syscall.EPIPE) {
			log.Fatal("broken pipe (EPIPE): ", err)
		} else if err != nil {
			log.Println(err)
		}

		fmt.Println(string(response))
	}
}
