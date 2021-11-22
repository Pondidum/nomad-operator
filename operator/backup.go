package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"text/template"

	_ "embed"

	"github.com/hashicorp/nomad/api"
	"github.com/hashicorp/nomad/jobspec"
)

//go:embed backup.nomad
var backupHcl string

type Backup struct {
	client *api.Client
}

func NewBackup(client *api.Client) *Backup {
	return &Backup{
		client: client,
	}
}

const (
	BackupFlag     = "auto-backup"
	BackupSchedule = "backup-schedule"
	BackupTargetDB = "backup-target-db"
)

func (b *Backup) OnJob(eventType string, job *api.Job) {

	if strings.HasPrefix(*job.ID, "backup-") {
		fmt.Println("    Job is a backup, skipping")
		return
	}

	settings, enabled := b.parseMeta(*job.ID, job.Meta)

	if eventType == "JobDeregistered" {
		fmt.Println("    Trying to remove a backup, if any")
		b.tryRemoveBackupJob(settings.JobID)
		return
	}

	if !enabled {
		fmt.Println("    Job has no auto-backup")
		fmt.Println("    Trying to remove a backup, if any")
		b.tryRemoveBackupJob(settings.JobID)
		return
	}

	fmt.Println("    Registering backup job")
	if err := b.createBackupJob(settings.JobID, settings); err != nil {
		fmt.Printf("--> Error creating job: %v\n", err)
	}

	fmt.Println("--> Done")
}

type settings struct {
	JobID       string
	SourceJobID string
	Schedule    string
	TargetDB    string
}

func (b *Backup) parseMeta(jobID string, meta map[string]string) (settings, bool) {

	s := settings{
		JobID:       "backup-" + jobID,
		SourceJobID: jobID,
	}

	enabled, found := meta[BackupFlag]
	if !found {
		return s, false
	}
	if active, _ := strconv.ParseBool(enabled); !active {
		return s, false
	}

	if schedule, found := meta[BackupSchedule]; found {
		s.Schedule = schedule
	}

	if target, found := meta[BackupTargetDB]; found {
		s.TargetDB = target
	}

	return s, true
}

func (b *Backup) tryRemoveBackupJob(jobID string) {
	b.client.Jobs().Deregister(jobID, false, &api.WriteOptions{})
}

func (b *Backup) createBackupJob(id string, bs settings) error {

	t, err := template.New("").Delims("[[", "]]").Parse(backupHcl)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	if err := t.Execute(&buffer, bs); err != nil {
		return err
	}

	backup, err := jobspec.Parse(&buffer)
	if err != nil {
		return err
	}

	_, _, err = b.client.Jobs().Register(backup, nil)
	if err != nil {
		return err
	}

	fmt.Printf("    Backup created: %s\n", id)
	return nil
}
