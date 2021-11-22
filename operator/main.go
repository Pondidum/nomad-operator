package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/nomad/api"
)

func main() {

	if err := run(os.Args[:]); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(args []string) error {

	client, err := api.NewClient(&api.Config{})
	if err != nil {
		return err
	}

	backup := NewBackup(client)
	consumer := NewConsumer(client, backup.OnJob)

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		s := <-signals
		fmt.Printf("Received %s, stopping\n", s)

		consumer.Stop()
		os.Exit(0)
	}()

	// blocks
	consumer.Start()
	return nil
}
