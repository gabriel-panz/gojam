package errors

import "errors"

var ErrExpiredToken error = errors.New("the access token expired")
var ErrUnauthorized error = errors.New("unauthorized")
var ErrNotFound error = errors.New("not foud")
var ErrBadRequest error = errors.New("not foud")
