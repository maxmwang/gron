package grond

import (
	"net"
	"net/http"
	"net/rpc"
	"sync"

	"github.com/google/uuid"
)

type Server interface {
	Add(c *Cron, result *string) error
	Remove(uuid uuid.UUID, result *string) error
	Snapshot(any, result *[]Cron) error
	start(errChan chan<- error)
}

type RPCServer struct {
	add      chan<- *Cron
	remove   chan<- uuid.UUID
	snapshot chan<- chan<- []Cron
}

func (s *RPCServer) Add(c *Cron, result *string) error {
	s.add <- c
	*result = "success"

	return nil
}

func (s *RPCServer) Remove(uuid uuid.UUID, result *string) error {
	s.remove <- uuid
	*result = "success"

	return nil
}

func (s *RPCServer) Snapshot(_ any, result *[]Cron) error {
	c := make(chan []Cron)
	s.snapshot <- c
	*result = <-c

	return nil
}

func (s *RPCServer) start(errChan chan<- error) {
	wg := sync.WaitGroup{}

	if err := rpc.Register(s); err != nil {
		errChan <- err
		return
	}

	rpc.HandleHTTP()
	l, err := net.Listen("tcp", ":1234")
	if err != nil {
		errChan <- err
		return
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := http.Serve(l, nil); err != nil {
			errChan <- err
			return
		}
	}()

	wg.Wait()
}
