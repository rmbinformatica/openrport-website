package notifications

import (
	"context"
	"sync"
)

type Processor interface {
	Close() error
}

type Consumer interface {
	Process(details NotificationDetails) error
	Target() Target
}

type Target string

const TargetMail Target = "smtp"
const TargetScript Target = "script"

var AllTargets = []Target{TargetMail, TargetScript}

func (t Target) Valid() bool {
	for _, target := range AllTargets {
		if t == target {
			return true
		}
	}
	return false
}

type Store interface {
	Create(ctx context.Context, details NotificationDetails) error
	LogRunning(ctx context.Context, nid string) error
	LogDone(ctx context.Context, nid string) error
	LogError(ctx context.Context, nid string, error string) error
	NotificationStream(target Target) chan NotificationDetails
}

type processor struct {
	store       Store
	done        context.CancelFunc
	waitForDead context.Context
	timeToDie   context.Context
	killMe      context.CancelFunc
	consumers   []Consumer
}

func (p *processor) start() {
	w := sync.WaitGroup{}
	w.Add(len(p.consumers))
	for _, c := range p.consumers {
		go func(consumer Consumer) {
			p.startConsumer(consumer)
			w.Done()
		}(c)
	}
	w.Wait()
	p.done()
}

func (p *processor) startConsumer(consumer Consumer) {
	updates := p.store.NotificationStream(consumer.Target())
root:
	for {
		select {
		case <-p.timeToDie.Done():
			break root
		case notification := <-updates:
			// TODO: should rport crush when save errors?
			p.store.LogRunning(context.Background(), notification.ID.ID())
			err := consumer.Process(notification)
			if err == nil {
				p.store.LogDone(context.Background(), notification.ID.ID())
			} else {
				p.store.LogError(context.Background(), notification.ID.ID(), err.Error())
			}
		}
	}
}

func (p *processor) Close() error {
	p.killMe()
	<-p.waitForDead.Done()
	return nil
}

func NewProcessor(store Store, consumers ...Consumer) Processor {
	ctx := context.Background()
	waitForDead, cancel := context.WithCancel(ctx)
	timeToDie, killMe := context.WithCancel(ctx)

	p := &processor{
		consumers:   consumers,
		store:       store,
		waitForDead: waitForDead,
		done:        cancel,
		timeToDie:   timeToDie,
		killMe:      killMe,
	}
	go p.start()
	return p
}