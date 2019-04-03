package portfmt

import "fmt"

// Port represents a server port
type Port int

// String converts a port number to a string to be used with net
func (p Port) String() string {
	if p < 1 {
		return ":"
	}
	return fmt.Sprintf(":%d", p)
}
