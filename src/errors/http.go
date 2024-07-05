package errors

import "errors"

var ErrExpiredToken error = errors.New("the access token expired")
var ErrUnauthorized error = errors.New("unauthorized")
var ErrNotFound error = errors.New("not found")
var ErrBadRequest error = errors.New("bad request")
var ErrSessionFull error = errors.New("this session is full")
