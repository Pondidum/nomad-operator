package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/nomad/api"
)

type Consumer struct {
	client *api.Client
	onJob  func(eventType string, job *api.Job)

	stop func()
}

func NewConsumer(client *api.Client, onJob func(eventType string, job *api.Job)) *Consumer {
	return &Consumer{
		client: client,
		onJob:  onJob,
	}
}

func (c *Consumer) Stop() {
	if c.stop != nil {
		c.stop()
	}
}

func (c *Consumer) Start() {

	ctx := context.Background()
	ctx, c.stop = context.WithCancel(ctx)

	c.consume(ctx)
}

func (c *Consumer) consume(ctx context.Context) error {

	var index uint64 = 0
	if _, meta, err := c.client.Jobs().List(nil); err == nil {
		index = meta.LastIndex + 1
	}

	topics := map[api.Topic][]string{
		api.TopicJob: {"*"},
	}

	eventsClient := c.client.EventStream()
	eventCh, err := eventsClient.Stream(ctx, topics, index, &api.QueryOptions{})
	if err != nil {
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case event := <-eventCh:
			// Ignore heartbeats used to keep connection alive
			if event.IsHeartbeat() {
				continue
			}

			c.handleEvent(event)
		}
	}

}

func (c *Consumer) handleEvent(event *api.Events) {
	if event.Err != nil {
		fmt.Printf("received error %s\n", event.Err)
		return
	}

	for _, e := range event.Events {

		if e.Type != "JobRegistered" && e.Type != "JobDeregistered" {
			return
		}

		job, err := e.Job()
		if err != nil {
			fmt.Printf("    received error %s\n", err)
			return
		}
		if job == nil {
			return
		}

		fmt.Printf("==> %s: %s (%s)...\n", e.Type, *job.ID, *job.Status)

		c.onJob(e.Type, job)
	}

}
