package ioc

import (
	spb "google.golang.org/genproto/googleapis/rpc/status"
	"google.golang.org/grpc/status"
)

func ErrorFromGoError(err error) *Error {
	s, _ := status.FromError(err)
	return ErrorFromStatus(s)
}

func ErrorFromStatus(s *status.Status) *Error {
	pbs := s.Proto()
	return &Error{
		Code:    pbs.Code,
		Message: pbs.Message,
		Details: pbs.Details,
	}
}

func (e *Error) ToStatus() *status.Status {
	return status.FromProto(&spb.Status{
		Code:    e.Code,
		Message: e.Message,
		Details: e.Details,
	})
}
