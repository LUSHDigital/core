package portfmt_test

import (
	"testing"

	"github.com/LUSHDigital/core/workers/internal/portfmt"
)

func TestPort_String(t *testing.T) {
	tests := []struct {
		name string
		port int
		want string
	}{
		{
			name: "0 port",
			port: 0,
			want: ":",
		},
		{
			name: "1 port",
			port: 1,
			want: ":1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := portfmt.Port(tt.port).String(); got != tt.want {
				t.Errorf("Port.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
