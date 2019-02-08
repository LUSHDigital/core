package keys_test

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/LUSHDigital/microservice-core-golang/keys"
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
	foobarStringSource = keys.StringSource("foobar")
	barbazStringSource = keys.StringSource("barbaz")

	onePubPath  = keys.FileSource(path.Join(wd(), "fixtures", "one.pub"))
	onePubBytes = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCckVt+i52W4M6XuPyd3u40SPql
YbhRB9XiOBZJztokBc5SJet0i9OsakKKnLbZevsM3MPI+Oj4hwsqp9oLDrJ1LXJy
IqI0OfMqq0f+YiPc0A6Uou1HiMDGSt7grwHkPVF7PDYeiNIAFR6e+rdTdWGLulx3
eCLysKk3KiS+JZF/twIDAQAB
-----END PUBLIC KEY-----`)
	twoPubPath  = keys.FileSource(path.Join(wd(), "fixtures", "two.pub"))
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
		source        keys.Source
		expectedErr   error
		expectedBytes []byte
	}{
		{
			name: "when first source is healthy and correct",
			source: keys.Sources{
				zeroByteSource,
			},
			expectedBytes: []byte{0},
			expectedErr:   nil,
		},
		{
			name: "when all sources are healthy and correct",
			source: keys.Sources{
				oneByteSource,
				zeroByteSource,
			},
			expectedBytes: []byte{1},
			expectedErr:   nil,
		},
		{
			name: "when first source is faulty and the next is healthy and correct",
			source: keys.Sources{
				faultySource,
				zeroByteSource,
			},
			expectedBytes: []byte{0},
			expectedErr:   nil,
		},
		{
			name: "when all sources are faulty",
			source: keys.Sources{
				faultySource,
				faultySource,
				faultySource,
			},
			expectedBytes: nil,
			expectedErr:   keys.ErrNoSourcesResolved,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bts, err := c.source.Get(context.Background())
			deepEqual(t, c.expectedBytes, bts)
			deepEqual(t, c.expectedErr, err)
		})
	}

}
func TestStringSource(t *testing.T) {
	cases := []struct {
		name          string
		source        keys.Source
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
			source:        keys.StringSource(""),
			expectedBytes: nil,
			expectedErr:   keys.ErrEmptyString,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bts, err := c.source.Get(context.Background())
			deepEqual(t, c.expectedBytes, bts)
			deepEqual(t, c.expectedErr, err)
		})
	}
}
func TestFileSource(t *testing.T) {
	cases := []struct {
		name          string
		source        keys.Source
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
			deepEqual(t, c.expectedBytes, bts)
			deepEqual(t, c.expectedErr, err)
		})
	}
}
func TestHTTPSource(t *testing.T) {
	l, err := net.Listen("tcp", ":")
	if err != nil {
		panic(err)
	}
	port := l.Addr().(*net.TCPAddr).Port
	fs := http.FileServer(http.Dir(path.Join(wd(), "fixtures")))
	go http.Serve(l, fs)

	cases := []struct {
		name          string
		source        keys.Source
		expectedErr   error
		expectedBytes []byte
	}{
		{
			name:          "when source is 127.0.0.1/one.pub",
			source:        keys.HTTPSource(fmt.Sprintf("http://127.0.0.1:%d/one.pub", port)),
			expectedBytes: onePubBytes,
			expectedErr:   nil,
		},
		{
			name:          "when source is 127.0.0.1/two.pub",
			source:        keys.HTTPSource(fmt.Sprintf("http://127.0.0.1:%d/two.pub", port)),
			expectedBytes: twoPubBytes,
			expectedErr:   nil,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			bts, err := c.source.Get(context.Background())
			deepEqual(t, c.expectedBytes, bts)
			deepEqual(t, c.expectedErr, err)
		})
	}
}

func wd() string {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return pwd
}
