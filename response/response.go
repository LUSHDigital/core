package response

// MicroserviceReponse - A standardised reponse format for a microservice.
type MicroserviceReponse struct {
	Status  string                 `json:"string"`
	Code    int                    `json:"code"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data,omitempty"`
}
