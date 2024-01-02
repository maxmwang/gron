package grond

import (
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func Start() {
	errChan := make(chan error)
	addChan := make(chan *Cron)
	removeChan := make(chan uuid.UUID)
	snapshotChan := make(chan chan<- []Cron)

	m := manager{
		crons:  []Cron{},
		logger: logrus.New(),

		add:      addChan,
		remove:   removeChan,
		snapshot: snapshotChan,
	}
	s := RPCServer{
		add:      addChan,
		remove:   removeChan,
		snapshot: snapshotChan,
	}

	go s.start(errChan)
	go m.start(errChan)

	select {
	case err := <-errChan:
		panic(err)
	}
}
