package ecs

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/docker/docker/api/types/container"
)

// Task is the ECS task representation
type Task struct {
	Cluster       string
	TaskARN       string
	Family        string
	Revision      string
	DesiredStatus string
	KnownStatus   string
	Containers    []Container
	Limits        map[string]float64
	PullStartedAt time.Time
	PullStoppedAt time.Time
}

// Container is the ECS metadata container representation
type Container struct {
	ID            string `json:"DockerId"`
	Name          string
	DockerName    string
	Image         string
	ImageID       string
	Labels        map[string]string
	DesiredStatus string
	KnownStatus   string
	Limits        map[string]float64
	CreatedAt     time.Time
	StartedAt     time.Time
	Stats         container.StatsResponse
	Type          string
	Networks      []Network
}

// Network is a docker network configuration
type Network struct {
	NetworkMode   string
	IPv4Addresses []string
}

func unmarshalTask(r io.Reader) (*Task, error) {
	task := &Task{}
	err := json.NewDecoder(r).Decode(task)
	return task, err
}

// docker parsers
func unmarshalStats(r io.Reader) (map[string]*container.StatsResponse, error) {
	var statsMap map[string]*container.StatsResponse
	if err := json.NewDecoder(r).Decode(&statsMap); err != nil {
		return nil, err
	}
	return statsMap, nil
}

// interleaves Stats in to the Container objects in the Task
func mergeTaskStats(task *Task, stats map[string]*container.StatsResponse) {
	for i := range task.Containers {
		c := &task.Containers[i]
		if strings.Trim(c.ID, " ") == "" {
			continue
		}
		stat, ok := stats[c.ID]
		if !ok || stat == nil {
			continue
		}
		c.Stats = *stat
	}
}
