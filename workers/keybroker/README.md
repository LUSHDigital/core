# Key Broker
The package `core/workers/keybroker` implements a background broker conmtinous retrieval of public keys from multiple different type of sources.

## Examples

```go
broker := keybroker.NewRSA(&keybroker.Config{
    Source:   keybroker.JWTPublicKeySources,
    Interval: 5 * time.Second,
})

// Run the broker
go broker.Run(ctx, ioutil.Discard)

// Queue retrieval of new key
broker.Renew()

// Copy the current public key held by the broker
broker.Copy()
```