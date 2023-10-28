package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conn, err := net.Dial("tcp4", "localhost:5000")

	if err != nil {
		log.Fatal(err)
	}
	go func() {

		conn.Write([]byte(os.Args[1]))

		for {

			in := bufio.NewReader(os.Stdin)
			//out := bufio.NewWriter(os.Stdout)
			fmt.Printf("\033[%d;%dH", 30, 0)
			fmt.Print("You: ")

			message, err := in.ReadString('\n')

			if err != nil {
				log.Println(err)
				continue
			}

			conn.Write([]byte(message))
		}
	}()
	buffer := make([]byte, 1024)
	go func() {
		for {

			_, err := conn.Read(buffer)

			if err != nil {
				conn.Close()
				log.Fatal("Connection refused from server")
			}

			fmt.Println(string(buffer))
		}

	}()
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()

	<-ctx.Done()
	conn.Close()
}
