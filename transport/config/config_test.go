package config

import (
	"os"
	"testing"
)

func TestGetGatewayUri(t *testing.T) {
	tt := []struct {
		name string
		uri  string
	}{
		{
			name: "String URI",
			uri:  "test",
		},
		{
			name: "String with digits URI",
			uri:  "1234",
		},
		{
			name: "Blank URI",
			uri:  "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("SOA_GATEWAY_URI", tc.uri)

			uri := GetGatewayURI()
			if uri != tc.uri {
				t.Errorf("TestGetGatewayUri: %s: expected %v got %v", tc.name, tc.uri, uri)
			}
		})
	}
}

func TestGetServiceDomain(t *testing.T) {
	tt := []struct {
		name   string
		domain string
	}{
		{
			name:   "String domain",
			domain: "test.com",
		},
		{
			name:   "String with digits domain",
			domain: "1234.com",
		},
		{
			name:   "Blank domain",
			domain: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("SOA_DOMAIN", tc.domain)

			domain := GetServiceDomain()
			if domain != tc.domain {
				t.Errorf("TestGetServiceDomain: %s: expected %v got %v", tc.name, tc.domain, domain)
			}
		})
	}
}

func TestGetGatewayUrl(t *testing.T) {
	tt := []struct {
		name string
		url  string
	}{
		{
			name: "String URI",
			url:  "test.com",
		},
		{
			name: "String with digits URI",
			url:  "1234.com",
		},
		{
			name: "Blank URI",
			url:  "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			os.Setenv("SOA_GATEWAY_URL", tc.url)

			url := GetGatewayURL()
			if url != tc.url {
				t.Errorf("TestGetGatewayUri: %s: expected %v got %v", tc.name, tc.url, url)
			}
		})
	}
}
