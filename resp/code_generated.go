package resp

import "net/http"

func init() {
	register(ErrSuccess, http.StatusOK, "OK")
	register(ErrUnknown, http.StatusInternalServerError, "Internal server error")
	register(ErrBind, http.StatusBadRequest, "Error occurred while binding the request body to the struct")
	register(ErrValidation, http.StatusBadRequest, "Validation failed")
	register(ErrTokenInvalid, http.StatusUnauthorized, "Token invalid")
	register(ErrPageNotFound, http.StatusNotFound, "Page not found")
	register(ErrDatabase, http.StatusInternalServerError, "Database error")
	register(ErrEncrypt, http.StatusUnauthorized, "Error occurred while encrypting the user password")
	register(ErrSignatureInvalid, http.StatusUnauthorized, "Signature is invalid")
	register(ErrExpired, http.StatusUnauthorized, "Token expired")
	register(ErrInvalidAuthHeader, http.StatusUnauthorized, "Invalid authorization header")
	register(ErrMissingHeader, http.StatusUnauthorized, "The `Authorization` header was empty")
	register(ErrPasswordIncorrect, http.StatusUnauthorized, "Password was incorrect")
	register(ErrPermissionDenied, http.StatusForbidden, "Permission denied")
	register(ErrEncodingFailed, http.StatusInternalServerError, "Encoding failed due to an error with the data")
	register(ErrDecodingFailed, http.StatusInternalServerError, "Decoding failed due to an error with the data")
	register(ErrInvalidJSON, http.StatusInternalServerError, "Data is not valid JSON")
	register(ErrEncodingJSON, http.StatusInternalServerError, "JSON data could not be encoded")
	register(ErrDecodingJSON, http.StatusInternalServerError, "JSON data could not be decoded")
	register(ErrInvalidYaml, http.StatusInternalServerError, "Data is not valid Yaml")
	register(ErrEncodingYaml, http.StatusInternalServerError, "Yaml data could not be encoded")
	register(ErrDecodingYaml, http.StatusInternalServerError, "Yaml data could not be decoded")
}
