package keybroker

import (
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"log"
	"math/big"

	"github.com/dgrijalva/jwt-go"
)

// RSAPublicKeyCopier represents behaviour for distributing copies of public keys
type RSAPublicKeyCopier interface {
	Copy() rsa.PublicKey
}

// RSAPrivateKeyCopier represents behaviour for distributing copies of private keys
type RSAPrivateKeyCopier interface {
	Copy() rsa.PrivateKey
}

var (
	// DefaultPublicRSA is an empty RSA public key.
	DefaultPublicRSA = &rsa.PublicKey{E: 0, N: big.NewInt(0)}

	// DefaultPrivateRSA is an empty RSA private key.
	DefaultPrivateRSA = &rsa.PrivateKey{
		D:         big.NewInt(0),
		PublicKey: *DefaultPublicRSA,
		Primes:    []*big.Int{},
	}

	// DefaultRSA is an empty RSA public key.
	// DEPRECATED: DefaultRSA is deprecated in favour of DefaultPublicRSA
	DefaultRSA = DefaultPublicRSA
)

// NewRSA returns a rsa public key broker based on configuration.
// DEPRECATED: The function keybroker.NewRSA() has been deprecated in favour of keybroker.NewPublicRSA()
func NewRSA(config *Config) *RSAPublicKeyBroker {
	log.Println("DEPRECATED: The function keybroker.NewRSA() has been deprecated in favour of keybroker.NewPublicRSA()")
	return NewPublicRSA(config)
}

// NewPublicRSA returns a rsa public key broker based on configuration.
func NewPublicRSA(config *Config) *RSAPublicKeyBroker {
	if config == nil {
		config = &Config{}
	}
	if config.Source == nil {
		config.Source = JWTPublicKeySources
	}
	brk := &RSAPublicKeyBroker{
		broker: newBroker("rsa public key", config),
	}
	// Make sure the broker is marked for renewal immediately
	brk.Renew()
	return brk
}

// NewPrivateRSA returns a rsa private key broker based on configuration.
func NewPrivateRSA(config *Config) *RSAPrivateKeyBroker {
	if config == nil {
		config = &Config{}
	}
	if config.Source == nil {
		config.Source = JWTPrivateKeySources
	}
	brk := &RSAPrivateKeyBroker{
		broker: newBroker("rsa private key", config),
	}
	// Make sure the broker is marked for renewal immediately.
	brk.Renew()
	return brk
}

// RSAPublicKeyBroker defines the implementation for brokering an RSA public key.
type RSAPublicKeyBroker struct {
	broker *broker
	key    *rsa.PublicKey
}

// Copy returns a shallow copy o the RSA public key.
func (b *RSAPublicKeyBroker) Copy() rsa.PublicKey {
	if b.key == nil {
		return *DefaultPublicRSA
	}
	return *b.key
}

// Renew will inform the broker to force renewal of the key.
func (b *RSAPublicKeyBroker) Renew() {
	b.broker.Renew()
}

// Close stops the ticker and releases resources.
func (b *RSAPublicKeyBroker) Close() {
	b.broker.Close()
}

// Run will periodically try and the public key.
func (b *RSAPublicKeyBroker) Run(ctx context.Context, out io.Writer) error {
	go b.broker.Run(ctx, out)
	select {
	case res := <-b.broker.res:
		key, err := jwt.ParseRSAPublicKeyFromPEM(res)
		if err != nil {
			return fmt.Errorf("cannot parse rsa public key: %v", err)
		}
		fmt.Fprintf(out, "rsa public key broker found new key of size %d\n", key.Size())
		b.key = key
	case err := <-b.broker.err:
		return err
	case <-ctx.Done():
	}
	return nil
}

// Check will see if the broker is ready.
func (b *RSAPublicKeyBroker) Check() ([]string, bool) {
	if !b.broker.running {
		return []string{"rsa public key broker is not yet running"}, false
	}
	if b.key == nil {
		return []string{fmt.Sprintf("rsa public key broker has not yet retrieved a key")}, false
	}
	return []string{fmt.Sprintf("rsa public key broker has retrieved key of size %d", b.key.Size())}, true
}

// RSAPrivateKeyBroker defines the implementation for brokering an RSA public key
type RSAPrivateKeyBroker struct {
	broker *broker
	key    *rsa.PrivateKey
}

// Copy returns a shallow copy o the RSA private key.
func (b *RSAPrivateKeyBroker) Copy() rsa.PrivateKey {
	if b.key == nil {
		return *DefaultPrivateRSA
	}
	return *b.key
}

// Renew will inform the broker to force renewal of the key.
func (b *RSAPrivateKeyBroker) Renew() {
	b.broker.Renew()
}

// Close stops the ticker and releases resources.
func (b *RSAPrivateKeyBroker) Close() {
	b.broker.Close()
}

// Run will periodically try and the private key.
func (b *RSAPrivateKeyBroker) Run(ctx context.Context, out io.Writer) error {
	go b.broker.Run(ctx, out)
	select {
	case res := <-b.broker.res:
		key, err := jwt.ParseRSAPrivateKeyFromPEM(res)
		if err != nil {
			return fmt.Errorf("cannot parse rsa private key: %v", err)
		}
		fmt.Fprintf(out, "rsa private key broker found new key of size %d\n", key.Size())
		b.key = key
	case err := <-b.broker.err:
		return err
	case <-ctx.Done():
		return nil
	}
	return nil
}

// Check will see if the broker is ready.
func (b *RSAPrivateKeyBroker) Check() ([]string, bool) {
	if !b.broker.running {
		return []string{"rsa private key broker is not yet running"}, false
	}
	if b.key == nil {
		return []string{fmt.Sprintf("rsa private key broker has not yet retrieved a key")}, false
	}
	return []string{fmt.Sprintf("rsa private key broker has retrieved key of size %d", b.key.Size())}, true
}
