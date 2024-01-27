package xyz

import "context"

// XyzService is an interface that represents a service
//
// @service()
type XyzService interface {
	// MyMethod is a method that does something
	//
	// @http(method="GET", path="/my-method")
	MyMethod(ctx context.Context)
}

type xyzService struct {
}

func (s *xyzService) MyMethod(ctx context.Context) {

}

func NewXyzService() XyzService {
	return &xyzService{}
}
