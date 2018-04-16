package domain

import (
	"fmt"
	"testing"
)

func TestBuildServiceDNSName(t *testing.T) {
	tt := []struct {
		name             string
		service          string
		branch           string
		environment      string
		serviceNamespace string
		expectedDNSName  string
	}{
		{
			name:             "Normal data",
			service:          "test",
			branch:           "master",
			environment:      "staging",
			serviceNamespace: "test",
			expectedDNSName:  "test-master-staging.test",
		},
		{
			name:             "Extreme data",
			service:          "21323kl1j3913issvxc9vx0",
			branch:           "(!()*)(*!KJ",
			environment:      "sljsjlfjdkgj",
			serviceNamespace: ")ID`hdfy7d7f",
			expectedDNSName:  "21323kl1j3913issvxc9vx0-(!()*)(*!KJ-sljsjlfjdkgj.)ID`hdfy7d7f",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actualDNSName := BuildServiceDNSName(tc.service, tc.branch, tc.environment, tc.serviceNamespace)
			if actualDNSName != tc.expectedDNSName {
				t.Errorf("TestBuildServiceDNSName: %s: expected %v got %v", tc.name, tc.expectedDNSName, actualDNSName)
			}
		})
	}
}

func ExampleBuildServiceDNSName() {
	dnsName := BuildServiceDNSName("myservice", "master", "staging", "services")
	fmt.Println(dnsName)

	// Output: myservice-master-staging.services
}

func TestBuildCloudServiceUrl(t *testing.T) {
	tt := []struct {
		name                    string
		apiGatewayURL           string
		serviceNamespace        string
		serviceName             string
		expectedCloudServiceURL string
	}{
		{
			name:                    "Normal data",
			apiGatewayURL:           "test.com",
			serviceNamespace:        "test",
			serviceName:             "test",
			expectedCloudServiceURL: "test.com/test/test",
		},
		{
			name:                    "Extreme data",
			apiGatewayURL:           "te(SDS(sdsdsdst.com",
			serviceNamespace:        "sdfisfpsif9((DF",
			serviceName:             "D&D&*FDHFHSDFHDF",
			expectedCloudServiceURL: "te(SDS(sdsdsdst.com/sdfisfpsif9((DF/D&D&*FDHFHSDFHDF",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			actualCloudServiceURL := BuildCloudServiceURL(tc.apiGatewayURL, tc.serviceNamespace, tc.serviceName)
			if actualCloudServiceURL != tc.expectedCloudServiceURL {
				t.Errorf("TestBuildServiceDNSName: %s: expected %v got %v", tc.name, tc.expectedCloudServiceURL, actualCloudServiceURL)
			}
		})
	}
}

func ExampleBuildCloudServiceUrl() {
	cloudServiceURL := BuildCloudServiceURL("my-api-gateway.com", "services", "myservice")
	fmt.Println(cloudServiceURL)

	// Output: my-api-gateway.com/services/myservice
}
