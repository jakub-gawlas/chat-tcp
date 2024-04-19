package main

import (
	"fmt"
	"os"
	"os/signal"
)

func main() {
	server := NewServer("localhost", "8080")

	go func() {
		if err := server.Listen(); err != nil {
			panic(err)
		}
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh

	server.Close()

	fmt.Println("Server terminated.")
}
