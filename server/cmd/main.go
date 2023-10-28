package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os/signal"
	"syscall"
)

func main() {
	fmt.Println("console.log(\"PENIS\")")
	listener, err := net.Listen("tcp4", "localhost:5000")

	if err != nil {
		fmt.Println(err)
	}

	//connections := make([]net.Conn, 5)
	connections := make(map[net.Conn]string)

	//fmt.Println(connections)
	go func() {

		for i := 0; i < 5; i++ {

			conn, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			buffer := make([]byte, 1024)
			_, err = conn.Read(buffer)
			if errors.Is(err, io.EOF) {
				log.Println("User disconnected")
				conn.Close()
				i--
				continue
			}

			connections[conn] = string(cutNil(buffer))
			fmt.Println(connections)

			if err != nil {
				log.Fatal(err)
			}

			go func(conn *net.Conn) {
				for {
					buffer := make([]byte, 1024)

					_, err := (*conn).Read(buffer)
					if errors.Is(err, io.EOF) {
						log.Printf("%s disconnected\n", connections[*conn])
						return
					}

					cuttedBuffer := cutNil(buffer)

					log.Println(cuttedBuffer)

					for key, _ := range connections {

						if key != *conn {
							log.Println([]byte(connections[*conn]))
							//fmt.Println(fmt.Sprintf("%s:", connections[*conn]) + string(buffer))
							key.Write([]byte(fmt.Sprintf("%s:", connections[*conn]) + string(cuttedBuffer)))
						}

					}
				}
			}(&conn)

		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)

	defer stop()

	<-ctx.Done()

}

func cutNil(arr []byte) []byte {
	var newArr []byte

	for i := 0; arr[i] != 0; i++ {
		newArr = append(newArr, arr[i])
	}

	return newArr
}
