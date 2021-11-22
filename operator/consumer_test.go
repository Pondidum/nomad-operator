package main

import (
	_ "embed"
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
	"github.com/stretchr/testify/assert"
)

//go:embed data/events.json
var eventsJson string

//go:embed data/withbackup.nomad
var withBackupHcl string

func TestConsumingEvents(t *testing.T) {

	seenEvents := []string{}
	c := NewConsumer(nil, func(eventType string, job *api.Job) {

		seenEvents = append(seenEvents, eventType)
		assert.Equal(t, "example", *job.ID)
	})

	for _, line := range strings.Split(eventsJson, "\n") {
		var events api.Events
		json.Unmarshal([]byte(line), &events)

		c.handleEvent(&events)
	}

	assert.Len(t, seenEvents, 2)
	assert.Equal(t, []string{"JobRegistered", "JobDeregistered"}, seenEvents)
}

func TestConsumingRealApi(t *testing.T) {
	if val := os.Getenv("CI"); val != "true" {
		t.Skip("Only runs in CI (set CI envvar to 'true')")
	}

	wait := make(chan bool, 1)

	client, err := api.NewClient(&api.Config{})
	assert.NoError(t, err)

	seenJobID := ""
	c := NewConsumer(client, func(eventType string, job *api.Job) {
		seenJobID = *job.ID
		wait <- true
	})

	go c.Start()

	//register a job
	job, err := jobspec.Parse(strings.NewReader(withBackupHcl))
	assert.NoError(t, err)

	client.Jobs().Register(job, &api.WriteOptions{})

	<-wait

	assert.Equal(t, *job.ID, seenJobID)
}
