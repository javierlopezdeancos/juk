package main

import (
	"github.com/nats-io/nats"
)

// Dispatcher is responsable for registering new Job hooks and
// dispatching Jobs to the corresponding channels over NATS.
type Dispatcher struct {
	channels map[string]chan *Job
	conn     *nats.EncodedConn
	jobCh    chan *Job
	messages chan *JobMsg
}

// NewDispatcher returns a new Dispatcher instance given a URL
// or an error if the connection to NATS fails.
func NewDispatcher(url string) (*Dispatcher, error) {

	conn, err := CreateEncodedConn(url)
	if err != nil {
		return nil, err
	}

	dsp := &Dispatcher{
		channels: make(map[string]chan *Job),
		conn:     conn,
		jobCh:    make(chan *Job),
		messages: make(chan *JobMsg),
	}

	dsp.conn.Subscribe(RegisterJob, func(m *JobMsg) {

		if _, ok := dsp.channels[m.Name]; ok {
			return
		}

		dsp.channels[m.Name] = make(chan *Job)
		dsp.conn.BindSendChan(m.Name, dsp.channels[m.Name])

		// Notify
		dsp.messages <- m
	})

	return dsp, nil
}

func (d *Dispatcher) Close() {
	if d.conn != nil {
		d.conn.Close()
	}
}

func (d *Dispatcher) Messages() <-chan *JobMsg {
	return d.messages
}

// Schedule sends a Job to the Dispatcher's job channel.
func (d *Dispatcher) Schedule(j *Job) {
	d.jobCh <- j
}

// Run periodically checks the Dispatcher's job channel
// and redirects incoming Jobs to the correponding channel.
func (d *Dispatcher) Run() {
	for {
		select {
		case j := <-d.jobCh:
			go func(j *Job) {
				ch, ok := d.channels[j.Name]
				if ok {
					ch <- j
				}
			}(j)
		}
	}
}
