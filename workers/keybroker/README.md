# Key Broker
The package `core/workers/keybroker` implements a background broker conmtinous retrieval of public keys from multiple different type of sources.

## Configuration
The key broker will by default try to retrieve keys from sources specified in the environment. These are the available environment variables:

- `JWT_PUBLIC_KEY` the key as a string
- `JWT_PUBLIC_KEY_URL` the http url where the key can be retrieved
- `JWT_PUBLIC_KEY_PATH` the file path on disk where the key can be read

You can also put your key on the location `/usr/local/var/jwt.pub` and it will by default attempt to read it.

## Examples

```go
broker := keybroker.NewPublicRSA(&keybroker.Config{
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