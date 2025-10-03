package apierrors

import "errors"

var ErrBadJWT error = errors.New("malformed or missing JWT token")
var ErrBadRequest error = errors.New("bad request: Check documentation for further informations")
var ErrEndpointNotFound error = errors.New("endpoint not found")
