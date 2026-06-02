package gapi

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func fieldViolation(field string, err error) error {
	return status.Errorf(codes.InvalidArgument, "invalid field %s: %s", field, err)
}

func unauthenticatedError(err error) error {
	return status.Errorf(codes.Unauthenticated, "unauthorized: %s", err)
}
