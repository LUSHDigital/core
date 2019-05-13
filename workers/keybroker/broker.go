package keybroker

import (
	"context"
	"fmt"
	"io"
	"time"
)

// Renewer represents behaviour for marking a broker for renewal
type Renewer interface {
	Renew()
}

// Closer represents behaviour for closing a broker
type Closer interface {
	Close()
}

// Config represents broker configuration
type Config struct {
	Interval time.Duration
	Source   Source
}

func newBroker(keyType string, config *Config) *broker {
	if config.Interval == 0 {
		config.Interval = 5 * time.Second
	}
	return &broker{
		interval:  config.Interval,
		source:    config.Source,
		ticker:    time.NewTicker(config.Interval),
		renew:     make(chan struct{}, 1),
		cancelled: make(chan struct{}, 1),
		res:       make(chan []byte, 1),
		err:       make(chan error, 1),
		keyType:   keyType,
	}
}

type broker struct {
	interval  time.Duration
	source    Source
	ticker    *time.Ticker
	renew     chan struct{}
	cancelled chan struct{}
	res       chan []byte
	err       chan error
	running   bool
	keyType   string
}

// Renew will inform the broker to force renewal of the key
func (b *broker) Renew() {
	select {
	// Return early if the cancelled channel is already closed
	case <-b.cancelled:
		return
	case b.renew <- struct{}{}:
	// Exit if we can't send or receive on any channels
	default:
	}
}

// Close stops the ticker and releases resources
func (b *broker) Close() {
	// Close the cancelled channel first to stop all select switches.
	close(b.cancelled)
	b.ticker.Stop()
}

func (b *broker) Run(ctx context.Context, out io.Writer) {
	fmt.Fprintf(out, "running %s broker checking for new key every %d second(s)\n", b.keyType, b.interval/time.Second)
	b.running = true
	defer func() { b.running = false }()
	defer close(b.renew)
	for {
		select {
		case <-b.cancelled:
			b.err <- fmt.Errorf("%s broker cancelled", b.keyType)
			return
		case <-b.ticker.C:
			select {
			case <-b.renew:
				bts, err := b.source.Get(ctx)
				if err != nil {
					fmt.Fprintf(out, "%s broker interval error: %v\n", b.keyType, err)
					b.Renew()
				}
				b.res <- bts
			default:
			}
		case <-ctx.Done():
			b.err <- fmt.Errorf("%s broker quit due to context timeout", b.keyType)
			return
		}
	}
}
