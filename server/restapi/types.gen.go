// Package restapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.16.2 DO NOT EDIT.
package restapi

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Message string `json:"message"`
}

// UserInput defines model for UserInput.
type UserInput struct {
	Name string `json:"name" validate:"required"`
}

// UserResponse defines model for UserResponse.
type UserResponse struct {
	CreatedAt string `json:"createdAt"`
	Id        int    `json:"id"`
	Name      string `json:"name"`
}

// CreateUserJSONRequestBody defines body for CreateUser for application/json ContentType.
type CreateUserJSONRequestBody = UserInput
