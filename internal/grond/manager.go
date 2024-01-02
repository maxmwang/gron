package grond

import (
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type manager struct {
	crons  []Cron
	logger *logrus.Logger

	add      chan *Cron
	remove   chan uuid.UUID
	snapshot chan chan<- []Cron
}

func (m *manager) start(errChan chan<- error) {
	for {
		select {
		case c := <-m.add:
			m.crons = append(m.crons, *c)

			m.logger.WithField("uuid", c.uuid).Info("added cron")
		case uid := <-m.remove:
			r := false
			for i, c := range m.crons {
				if c.uuid == uid {
					m.crons = append(m.crons[:i], m.crons[i+1:]...)
					r = true
					break
				}
			}

			if r {
				m.logger.WithField("uuid", uid).Info("removed cron")
			} else {
				m.logger.WithField("uuid", uid).Warn("cron not found")
			}
		case c := <-m.snapshot:
			c <- m.crons
		}
	}
}

type Cron struct {
	uuid     uuid.UUID
	job      Job
	schedule scheduler
	next     time.Time
}

type Job struct{}

type scheduler interface {
	next(time time.Time) time.Time
}
