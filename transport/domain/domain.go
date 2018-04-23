package domain

// BuildServiceDNSName - Build the full DNS name for a service.
func BuildServiceDNSName(service, branch, environment, serviceNamespace string) string {
	return service + "-" + branch + "-" + environment + "." + serviceNamespace
}

// BuildCloudServiceURL - Build the full URL for a cloud service.
func BuildCloudServiceURL(apiGatewayURL, serviceNamespace, serviceName string) string {
	return apiGatewayURL + "/" + serviceNamespace + "/" + serviceName
}
