package keys

import (
	"context"
	"crypto/rsa"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// Renewer represents behaviour for marking a broker for renewal
type Renewer interface {
	Renew()
}

// Closer represents behaviour for closing a broker
type Closer interface {
	Close()
}

// RSAPublicKeyCopier represents behaviour for distributing copies of public keys
type RSAPublicKeyCopier interface {
	Copy() rsa.PublicKey
}

var (
	// DefaultRSA is an empty RSA public key
	DefaultRSA = &rsa.PublicKey{E: 0, N: big.NewInt(0)}
)

// BrokerRSAPublicKey will broker a public key from a source on an interval
func BrokerRSAPublicKey(ctx context.Context, source Source, interval time.Duration) *RSAPublicKeyBroker {
	broker := &RSAPublicKeyBroker{
		source:    source,
		ticker:    time.NewTicker(interval),
		key:       DefaultRSA,
		renew:     make(chan struct{}, 1),
		cancelled: make(chan struct{}, 1),
	}

	// Make sure the broker is marked for renewal immediately
	broker.Renew()

	// Begin the key renewal
	go broker.run(ctx)

	// Return the broker together with a separate cancel function
	// We do this to ensure cancellation is handeled correctly
	return broker
}

// RSAPublicKeyBroker defines the implementation for brokering an RSA public key
type RSAPublicKeyBroker struct {
	source    Source
	ticker    *time.Ticker
	key       *rsa.PublicKey
	renew     chan struct{}
	cancelled chan struct{}
}

// Copy returns a shallow copy o the RSA public key
func (b *RSAPublicKeyBroker) Copy() rsa.PublicKey {
	return *b.key
}

// Renew will inform the broker to force renewal of the key
func (b *RSAPublicKeyBroker) Renew() {
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
func (b *RSAPublicKeyBroker) Close() {
	// Close the cancelled channel first to stop all select switches
	close(b.cancelled)
	b.ticker.Stop()
	close(b.renew)
}

// Run will periodically try and the public key
func (b *RSAPublicKeyBroker) run(ctx context.Context) {
	for {
		select {
		case <-b.cancelled:
			log.Printf("rsa public key broker cancelled\n")
			return
		case <-b.ticker.C:
			select {
			case <-b.renew:
				if err := b.get(ctx); err != nil {
					log.Printf("rsa public key broker interval error: %v\n", err)
				}
			default:
			}
		case <-ctx.Done():
			log.Printf("rsa public key broker quit due to context timeout\n")
			return
		}
	}
}

func (b *RSAPublicKeyBroker) get(ctx context.Context) error {
	bts, err := b.source.Get(ctx)
	if err != nil {
		return fmt.Errorf("cannot get key: %v", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(bts)
	if err != nil {
		return fmt.Errorf("cannot parse key: %v", err)
	}
	b.key = key
	return nil
}
