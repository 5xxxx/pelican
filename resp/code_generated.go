package resp

import "net/http"

func init() {
	Register(ErrSuccess, http.StatusOK, "OK")
	Register(ErrUnknown, http.StatusInternalServerError, "Internal server error")
	Register(ErrBind, http.StatusBadRequest, "Error occurred while binding the request body to the struct")
	Register(ErrValidation, http.StatusBadRequest, "Validation failed")
	Register(ErrTokenInvalid, http.StatusUnauthorized, "Token invalid")
	Register(ErrPageNotFound, http.StatusNotFound, "Page not found")
	Register(ErrDatabase, http.StatusInternalServerError, "Database error")
	Register(ErrEncrypt, http.StatusUnauthorized, "Error occurred while encrypting the user password")
	Register(ErrSignatureInvalid, http.StatusUnauthorized, "Signature is invalid")
	Register(ErrExpired, http.StatusUnauthorized, "Token expired")
	Register(ErrInvalidAuthHeader, http.StatusUnauthorized, "Invalid authorization header")
	Register(ErrMissingHeader, http.StatusUnauthorized, "The `Authorization` header was empty")
	Register(ErrPasswordIncorrect, http.StatusUnauthorized, "Password was incorrect")
	Register(ErrPermissionDenied, http.StatusForbidden, "Permission denied")
	Register(ErrEncodingFailed, http.StatusInternalServerError, "Encoding failed due to an error with the data")
	Register(ErrDecodingFailed, http.StatusInternalServerError, "Decoding failed due to an error with the data")
	Register(ErrInvalidJSON, http.StatusInternalServerError, "Data is not valid JSON")
	Register(ErrEncodingJSON, http.StatusInternalServerError, "JSON data could not be encoded")
	Register(ErrDecodingJSON, http.StatusInternalServerError, "JSON data could not be decoded")
	Register(ErrInvalidYaml, http.StatusInternalServerError, "Data is not valid Yaml")
	Register(ErrEncodingYaml, http.StatusInternalServerError, "Yaml data could not be encoded")
	Register(ErrDecodingYaml, http.StatusInternalServerError, "Yaml data could not be decoded")
}
