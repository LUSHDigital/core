package response

import (
	"fmt"
	"reflect"
	"testing"
)

// An example data type.
var expectedResponseDataType = "tests"

// An example data set for testing with.
var expectedResponseData = map[string]interface{}{
	"tests":    "ok",
	"language": "golang",
}

// An example response object for testing with.
var expectedResponse = &MicroserviceReponse{
	Status:  StatusOk,
	Code:    200,
	Message: "",
	Data: map[string]interface{}{
		"tests": expectedResponseData,
	},
}

// An example response object (with no data) for testing with.
var expectedResponseNoData = &MicroserviceReponse{
	Status:  StatusOk,
	Code:    200,
	Message: "",
}

// TestResponseObject - Check the response creator is providing a valid
// response.
func TestResponseObject(t *testing.T) {
	// Create a response.
	response := CreateResponse(expectedResponseDataType, expectedResponseData, 200, StatusOk, "")

	// Check the response.
	if !reflect.DeepEqual(response, expectedResponse) {
		t.Error(fmt.Sprintf("Expected %v, got %v", response, expectedResponse))
	}
}

// TestResponseObjectNoData - Check the response creator is providing a valid
// response if there is no data.
func TestResponseObjectNoData(t *testing.T) {
	// Create a response.
	response := CreateResponse(expectedResponseDataType, nil, 200, StatusOk, "")

	// Check the response.
	if !reflect.DeepEqual(response, expectedResponseNoData) {
		t.Error(fmt.Sprintf("Expected %v, got %v", response, expectedResponseNoData))
	}
}

// ExampleCreateResponse - Example usage for the CreateResponse function.
func ExampleCreateResponse() {
	someThings := map[string]string{
		"thing_one": "a thing",
		"thing_two": "another thing",
	}

	response := CreateResponse("things", someThings, 200, StatusOk, "")
	fmt.Printf("%+v", response)

	// Output:
	// &{Status:ok Code:200 Message: Data:map[things:map[thing_one:a thing thing_two:another thing]]}
}
