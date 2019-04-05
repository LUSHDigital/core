package keybroker_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"path"
	"testing"

	"github.com/LUSHDigital/core/workers/keybroker"
)

type SourceFunc func(ctx context.Context) ([]byte, error)

func (s SourceFunc) Get(ctx context.Context) ([]byte, error) {
	return s(ctx)
}

var (
	faultySource = SourceFunc(func(ctx context.Context) ([]byte, error) {
		return nil, fmt.Errorf("cannot get source")
	})
	zeroByteSource = SourceFunc(func(ctx context.Context) ([]byte, error) {
		return []byte{0}, nil
	})
	oneByteSource = SourceFunc(func(ctx context.Context) ([]byte, error) {
		return []byte{1}, nil
	})
	foobarStringSource = keybroker.StringSource("foobar")
	barbazStringSource = keybroker.StringSource("barbaz")

	onePubPath  = keybroker.FileSource(path.Join("testdata", "one.pub"))
	onePubBytes = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCckVt+i52W4M6XuPyd3u40SPql
YbhRB9XiOBZJztokBc5SJet0i9OsakKKnLbZevsM3MPI+Oj4hwsqp9oLDrJ1LXJy
IqI0OfMqq0f+YiPc0A6Uou1HiMDGSt7grwHkPVF7PDYeiNIAFR6e+rdTdWGLulx3
eCLysKk3KiS+JZF/twIDAQAB
-----END PUBLIC KEY-----`)
	twoPubPath  = keybroker.FileSource(path.Join("testdata", "two.pub"))
	twoPubBytes = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDkeaMV5IrxxcoK6xpFaR6wBCTp
1CZTOB3sCWuFG0YaGYo/4w4O2WVUUoYN4/dvZbHAyUAeeLT5+T4s6pLBebbzooU+
pAU+iLlsgMQHCqm5s+yUYjniST15suYIbjhvP1o6VNHp5JgOoL+b/EHfZZZUP6LB
Iy0Bo6vikx7871xzuwIDAQAB
-----END PUBLIC KEY-----`)
)

func TestSources(t *testing.T) {
	cases := []struct {
		name          string
		source        keybroker.Source
		expectedErr   error
		expectedBytes []byte
	}{
		{
			name: "when first source is healthy and correct",
			source: keybroker.Sources{
				zeroByteSource,
			},
			expectedBytes: []byte{0},
			expectedErr:   nil,
		},
		{
			name: "when all sources are healthy and correct",
			source: keybroker.Sources{
				oneByteSource,
				zeroByteSource,
			},
			expectedBytes: []byte{1},
			expectedErr:   nil,
		},
		{
			name: "when first source is faulty and the next is healthy and correct",
			source: keybroker.Sources{
				faultySource,
				zeroByteSource,
			},
			expectedBytes: []byte{0},
			expectedErr:   nil,
		},
		{
			name: "when all sources are faulty",
			source: keybroker.Sources{
				faultySource,
				faultySource,
				faultySource,
			},
			expectedBytes: nil,
			expectedErr: keybroker.ErrNoSourcesResolved{
				N: 3,
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bts, err := c.source.Get(context.Background())
			equals(t, c.expectedBytes, bts)
			equals(t, c.expectedErr, err)
		})
	}

}
func TestStringSource(t *testing.T) {
	cases := []struct {
		name          string
		source        keybroker.Source
		expectedErr   error
		expectedBytes []byte
	}{
		{
			name:          "when source is foobar",
			source:        foobarStringSource,
			expectedBytes: []byte("foobar"),
			expectedErr:   nil,
		},
		{
			name:          "when source is barbaz",
			source:        barbazStringSource,
			expectedBytes: []byte("barbaz"),
			expectedErr:   nil,
		},
		{
			name:          "when source is empty",
			source:        keybroker.StringSource(""),
			expectedBytes: nil,
			expectedErr:   keybroker.ErrEmptyString,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bts, err := c.source.Get(context.Background())
			equals(t, c.expectedBytes, bts)
			equals(t, c.expectedErr, err)
		})
	}
}
func TestFileSource(t *testing.T) {
	cases := []struct {
		name          string
		source        keybroker.Source
		expectedErr   error
		expectedBytes []byte
	}{
		{
			name:          "when source is one.pub",
			source:        onePubPath,
			expectedBytes: onePubBytes,
			expectedErr:   nil,
		},
		{
			name:          "when source is two.pub",
			source:        twoPubPath,
			expectedBytes: twoPubBytes,
			expectedErr:   nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bts, err := c.source.Get(context.Background())
			equals(t, c.expectedBytes, bts)
			equals(t, c.expectedErr, err)
		})
	}
}
func TestHTTPSource(t *testing.T) {
	l, err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	fs := http.FileServer(http.Dir(path.Join("testdata")))
	go http.Serve(l, fs)

	cases := []struct {
		name          string
		source        keybroker.Source
		expectedErr   error
		expectedBytes []byte
	}{
		{
			name:          "when source is 127.0.0.1/one.pub",
			source:        keybroker.HTTPSource(fmt.Sprintf("http://127.0.0.1:%d/one.pub", port)),
			expectedBytes: onePubBytes,
			expectedErr:   nil,
		},
		{
			name:          "when source is 127.0.0.1/two.pub",
			source:        keybroker.HTTPSource(fmt.Sprintf("http://127.0.0.1:%d/two.pub", port)),
			expectedBytes: twoPubBytes,
			expectedErr:   nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bts, err := c.source.Get(context.Background())
			equals(t, c.expectedBytes, bts)
			equals(t, c.expectedErr, err)
		})
	}
}
