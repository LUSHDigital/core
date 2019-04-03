# Keys
The package `core/workers/keybroker` implements a background broker conmtinous retrieval of public keys from multiple different type of sources.

## Examples

```go
broker := keys.BrokerRSAPublicKey(context.Background(), keys.JWTPublicKeySources, 5*time.Second)
defer broker.Close()

// Queue retrieval of new key
broker.Renew()

// Copy the current public key held by the broker
publicKey := broker.Copy()
```