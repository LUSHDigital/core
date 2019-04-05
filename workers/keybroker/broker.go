package keybroker

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
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

// Config represents broker configuration
type Config struct {
	Interval time.Duration
	Source   Source
}

var (
	// DefaultRSA is an empty RSA public key
	DefaultRSA = &rsa.PublicKey{E: 0, N: big.NewInt(0)}
)

// NewRSA returns a rsa public key broker based on configuration.
func NewRSA(config *Config) *RSAPublicKeyBroker {
	if config == nil {
		config = &Config{}
	}
	if config.Source == nil {
		config.Source = JWTPublicKeySources
	}
	if config.Interval == 0 {
		config.Interval = 5 * time.Second
	}

	broker := &RSAPublicKeyBroker{
		interval:  config.Interval,
		source:    config.Source,
		ticker:    time.NewTicker(config.Interval),
		key:       DefaultRSA,
		renew:     make(chan struct{}, 1),
		cancelled: make(chan struct{}, 1),
	}

	// Make sure the broker is marked for renewal immediately
	broker.Renew()
	return broker
}

// RSAPublicKeyBroker defines the implementation for brokering an RSA public key
type RSAPublicKeyBroker struct {
	interval  time.Duration
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

// Run will periodically try and the public key.
func (b *RSAPublicKeyBroker) Run(ctx context.Context, out io.Writer) error {
	fmt.Fprintf(out, "running rsa public key broker checking for new key every %d second(s)\n", b.interval/time.Second)
	for {
		select {
		case <-b.cancelled:
			return fmt.Errorf("rsa public key broker cancelled")
		case <-b.ticker.C:
			select {
			case <-b.renew:
				if err := b.get(ctx); err != nil {
					fmt.Fprintf(out, "rsa public key broker interval error: %v\n", err)
					b.Renew()
				}
			default:
			}
		case <-ctx.Done():
			return fmt.Errorf("rsa public key broker quit due to context timeout")
		}
	}
}

func (b *RSAPublicKeyBroker) get(ctx context.Context) error {
	bts, err := b.source.Get(ctx)
	if err != nil {
		return fmt.Errorf("cannot get rsa public key: %v", err)
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(bts)
	if err != nil {
		return fmt.Errorf("cannot parse rsa public key: %v", err)
	}
	b.key = key
	return nil
}
