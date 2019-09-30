# Workers
The `core/workers` package is used to setup gracefully terminating workers for services.

## Examples

### The worker interfaces
All workers implement a simple interface for running and halting, using the context package for timeouts.

```go
// Runner represents the behaviour for running a service worker.
type Runner interface {
	// Run should run start processing the worker and be a blocking operation.
	Run(context.Context) error
}

// Halter represents the behaviour for stopping a service worker.
type Halter interface {
	// Halt should tell the worker to stop doing work.
	Halt(context.Context) error
}

// Worker represents the behaviour for a service worker.
type Worker interface {
	Runner
	Halter
}
```
