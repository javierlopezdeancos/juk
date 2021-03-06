package main

import "github.com/nats-io/nats"

const (
	RegisterJob = "job.register"
	UpdateJob   = "job.update"
)

type Job struct {
	Data    []byte              `json:"data,omitempty"`
	Name    string              `json:"name"`
	Headers map[string][]string `json:"headers"`
	Vars    map[string]string   `json:"vars"`
}

type JobMsg struct {
	Name            string            `json:"name"`
	Path            string            `json:"path"`
	Methods         []string          `json:"methods,omitempty"`
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`
	StatusCode      int               `json:"status_code,omitempty"`
}

func CreateEncodedConn(url string) (*nats.EncodedConn, error) {
	nc, err := nats.Connect("nats://" + url)
	if err != nil {
		return nil, err
	}
	conn, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
