package apierrors

import "errors"

var ErrBadJWT error = errors.New("Malformed or missing JWT token.")
var ErrBadRequestBody error = errors.New("Bad request body.")
