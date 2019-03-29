// Package models defines common data types used with the transport package.
package models

// Consumer - Holds information about an API consumer.
type Consumer struct {
	Tokens []*Token `json:"tokens,omitempty"` // The API consumers current access token.
}
